package opa

// nolint:lll
//go:generate go-bindata -pkg $GOPACKAGE -o policy.bindata.go -ignore .*_test.rego -ignore Makefile -ignore README\.md policy/...

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"github.com/open-policy-agent/opa/topdown"

	"github.com/tx7do/kratos-authz/engine"
)

var _ engine.Engine = (*State)(nil)

type State struct {
	store                storage.Store
	queries              map[string]ast.Body
	compiler             *ast.Compiler
	modules              map[string]*ast.Module
	preparedEvalProjects rego.PreparedEvalQuery

	regoVersion       ast.RegoVersion
	enableQueryTracer bool

	authzProjectsQuery    string
	filteredPairsQuery    string
	filteredProjectsQuery string

	log *log.Helper
}

func NewEngine(_ context.Context, opts ...OptFunc) (*State, error) {
	s := State{
		store:                 inmem.New(),
		queries:               make(map[string]ast.Body),
		log:                   log.NewHelper(log.With(log.DefaultLogger, "module", "opa.authz.engine")),
		regoVersion:           ast.DefaultRegoVersion,
		enableQueryTracer:     false,
		authzProjectsQuery:    defaultAuthzProjectsQuery,
		filteredPairsQuery:    defaultFilteredPairsQuery,
		filteredProjectsQuery: defaultFilteredProjectsQuery,
	}

	if err := s.init(opts...); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *State) init(opts ...OptFunc) error {
	var err error

	for _, opt := range opts {
		opt(s)
	}

	if err = s.initQueries(); err != nil {
		return errors.Wrap(err, "init queries")
	}

	if err = s.initModules(); err != nil {
		return errors.Wrap(err, "init OPA modules")
	}

	return nil
}

func (s *State) Name() string {
	return string(engine.Opa)
}

func (s *State) ParseProjectsQuery(query string) error {
	if query == "" {
		query = defaultAuthzProjectsQuery
	}

	s.authzProjectsQuery = query

	authzProjectsQueryParsed, err := ast.ParseBody(query)
	if err != nil {
		s.log.Errorf("failed to parse authz projects query %q: %v", query, err)
		return errors.Wrapf(err, "parse query %q", query)
	}

	if s.queries == nil {
		s.queries = make(map[string]ast.Body)
	}

	s.queries[AuthzProjectsQueryKey] = authzProjectsQueryParsed

	return nil
}

func (s *State) ParseFilterPairsQuery(query string) error {
	if query == "" {
		query = defaultFilteredPairsQuery
	}

	s.filteredPairsQuery = query

	filteredPairsQueryParsed, err := ast.ParseBody(query)
	if err != nil {
		s.log.Errorf("failed to parse filtered pairs query %q: %v", query, err)
		return errors.Wrapf(err, "parse query %q", query)
	}

	if s.queries == nil {
		s.queries = make(map[string]ast.Body)
	}

	s.queries[FilteredPairsQueryKey] = filteredPairsQueryParsed

	return nil
}

func (s *State) ParseFilterProjectsQuery(query string) error {
	if query == "" {
		query = defaultFilteredProjectsQuery
	}

	s.filteredProjectsQuery = query

	filteredProjectsQueryParsed, err := ast.ParseBody(query)
	if err != nil {
		s.log.Errorf("failed to parse filtered projects query %q: %v", query, err)
		return errors.Wrapf(err, "parse query %q", query)
	}

	if s.queries == nil {
		s.queries = make(map[string]ast.Body)
	}

	s.queries[FilteredProjectsQueryKey] = filteredProjectsQueryParsed

	return nil
}

func (s *State) ProjectsAuthorized(
	ctx context.Context,
	subjects engine.Subjects,
	action engine.Action,
	resource engine.Resource,
	projects engine.Projects,
) (engine.Projects, error) {
	var subs []*ast.Term
	for _, sub := range subjects {
		subs = append(subs, ast.NewTerm(ast.String(sub)))
	}

	var projs []*ast.Term
	for _, proj := range projects {
		projs = append(projs, ast.NewTerm(ast.String(proj)))
	}

	input := ast.NewObject(
		[2]*ast.Term{ast.NewTerm(ast.String("subjects")), ast.ArrayTerm(subs...)},
		[2]*ast.Term{ast.NewTerm(ast.String("resource")), ast.NewTerm(ast.String(resource))},
		[2]*ast.Term{ast.NewTerm(ast.String("action")), ast.NewTerm(ast.String(action))},
		[2]*ast.Term{ast.NewTerm(ast.String("projects")), ast.ArrayTerm(projs...)},
	)
	resultSet, err := s.preparedEvalProjects.Eval(ctx, rego.EvalParsedInput(input))
	if err != nil {
		s.log.Errorf("failed to evaluate projects query: %v", err)
		return engine.Projects{}, &EvaluationError{e: err}
	}

	return s.projectsFromPreparedEvalQuery(resultSet)
}

func (s *State) FilterAuthorizedPairs(
	ctx context.Context,
	subjects engine.Subjects,
	pairs engine.Pairs,
) (engine.Pairs, error) {
	opaInput := map[string]interface{}{
		"subjects": subjects,
		"pairs":    pairs,
	}

	rs, err := s.evalQuery(ctx, s.queries[FilteredPairsQueryKey], opaInput, s.store)
	if err != nil {
		s.log.Errorf("failed to evaluate filtered pairs query: %v", err)
		return nil, &EvaluationError{e: err}
	}

	return s.pairsFromResults(rs)
}

func (s *State) FilterAuthorizedProjects(ctx context.Context, subjects engine.Subjects) (engine.Projects, error) {
	opaInput := map[string]interface{}{
		"subjects": subjects,
	}

	rs, err := s.evalQuery(ctx, s.queries[FilteredProjectsQueryKey], opaInput, s.store)
	if err != nil {
		s.log.Errorf("failed to evaluate filtered projects query: %v", err)
		return nil, &EvaluationError{e: err}
	}

	return s.projectsFromPartialResults(rs)
}

func (s *State) IsAuthorized(
	ctx context.Context,
	subject engine.Subject,
	action engine.Action,
	resource engine.Resource,
	project engine.Project,
) (bool, error) {
	if len(project) > 0 {
		input := ast.NewObject(
			[2]*ast.Term{ast.NewTerm(ast.String("subjects")), ast.ArrayTerm(ast.NewTerm(ast.String(subject)))},
			[2]*ast.Term{ast.NewTerm(ast.String("resource")), ast.NewTerm(ast.String(resource))},
			[2]*ast.Term{ast.NewTerm(ast.String("action")), ast.NewTerm(ast.String(action))},
			[2]*ast.Term{ast.NewTerm(ast.String("projects")), ast.ArrayTerm(ast.NewTerm(ast.String(project)))},
		)
		resultSet, err := s.preparedEvalProjects.Eval(ctx, rego.EvalParsedInput(input))
		if err != nil {
			s.log.Errorf("failed to evaluate projects query: %v", err)
			return false, &EvaluationError{e: err}
		}
		return s.allowedFromPreparedEvalQuery(resultSet)
	} else {
		opaInput := map[string]interface{}{
			"subjects": engine.MakeSubjects(subject),
			"pairs":    engine.MakePairs(engine.Pair{Resource: resource, Action: action}),
		}

		rs, err := s.evalQuery(ctx, s.queries[FilteredPairsQueryKey], opaInput, s.store)
		if err != nil {
			s.log.Errorf("failed to evaluate filtered pairs query: %v", err)
			return false, &EvaluationError{e: err}
		}

		return s.pairsFromAllowed(rs)
	}
}

func (s *State) SetPolicies(ctx context.Context, policyMap engine.PolicyMap, roleMap engine.RoleMap) error {
	s.store = inmem.NewFromObject(map[string]interface{}{
		"policies": policyMap,
		"roles":    roleMap,
	})

	return s.makeAuthorizedProjectPreparedQuery(ctx)
}

func (s *State) InitModulesFromFiles(modules map[string]string) error {
	parsedModules := map[string]*ast.Module{}
	for name, path := range modules {
		moduleData, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "read module file %q", path)
		}

		parsed, err := ast.ParseModule(name, string(moduleData))
		if err != nil {
			s.log.Errorf("failed to parse module file %q: %v", name, err)
			return errors.Wrapf(err, "parse module %q", name)
		}

		parsedModules[name] = parsed
	}

	s.modules = parsedModules

	return nil
}

func (s *State) InitModulesFromString(modules map[string]string) error {
	parsedModules := map[string]*ast.Module{}
	for name, moduleData := range modules {
		parsed, err := ast.ParseModule(name, moduleData)
		if err != nil {
			s.log.Errorf("failed to parse module file %q: %v", name, err)
			return errors.Wrapf(err, "parse module %q", name)
		}

		parsedModules[name] = parsed
	}

	s.modules = parsedModules

	return nil
}

func (s *State) InitModulesFromAssets() error {
	mods := map[string]*ast.Module{}
	for _, name := range AssetNames() {
		if !strings.HasSuffix(name, ".rego") {
			continue
		}
		parsed, err := ast.ParseModule(name, string(MustAsset(name)))
		if err != nil {
			s.log.Errorf("failed to parse policy file %q: %v", name, err)
			return errors.Wrapf(err, "parse policy file %q", name)
		}
		mods[name] = parsed
	}

	s.modules = mods

	return nil
}

func (s *State) doCompile() error {
	compiler, err := s.newCompiler()
	if err != nil {
		s.log.Errorf("failed to create compiler: %v", err)
		return errors.Wrap(err, "init compiler")
	}

	s.compiler = compiler

	return nil
}

func (s *State) initQueries() error {
	var err error

	if err = s.ParseProjectsQuery(s.authzProjectsQuery); err != nil {
		return errors.Wrap(err, "parse projects query")
	}
	if err = s.ParseFilterPairsQuery(s.filteredPairsQuery); err != nil {
		return errors.Wrap(err, "parse filter pairs query")
	}
	if err = s.ParseFilterProjectsQuery(s.filteredProjectsQuery); err != nil {
		return errors.Wrap(err, "parse filter projects query")
	}

	return nil
}

func (s *State) initModules() error {
	if len(s.modules) == 0 {
		if err := s.InitModulesFromAssets(); err != nil {
			return errors.Wrap(err, "init modules from assets")
		}
	}

	if err := s.doCompile(); err != nil {
		return errors.Wrap(err, "init compiler")
	}

	return nil
}

func (s *State) makeAuthorizedProjectPreparedQuery(ctx context.Context) error {
	compiler, err := s.newCompiler()
	if err != nil {
		s.log.Errorf("failed to create compiler: %v", err)
		return err
	}

	r := rego.New(
		rego.Store(s.store),
		rego.Compiler(compiler),
		rego.ParsedQuery(s.queries[AuthzProjectsQueryKey]),
		rego.DisableInlining([]string{
			"data.authz.denied_project",
		}),
		rego.SetRegoVersion(s.regoVersion),
	)

	pq, err := r.Partial(ctx)
	if err != nil {
		s.log.Errorf("failed to create partial query for authorized projects: %v", err)
		return err
	}

	for i, module := range pq.Support {
		compiler.Modules[fmt.Sprintf("support%d", i)] = module
	}

	main := &ast.Module{
		Package: ast.MustParsePackage("package __partialauthz"),
	}

	for i := range pq.Queries {
		rule := &ast.Rule{
			Module: main,
			Head:   ast.NewHead("authorized_project", ast.VarTerm("project")),
			Body:   pq.Queries[i],
		}
		main.Rules = append(main.Rules, rule)
	}

	compiler.Modules["__partialauthz"] = main

	compiler.Compile(compiler.Modules)

	if compiler.Failed() {
		s.log.Errorf("failed to compile authorized projects: %v", compiler.Errors)
		return compiler.Errors
	}

	r2 := rego.New(
		rego.Store(s.store),
		rego.Compiler(compiler),
		rego.Query("data.__partialauthz.authorized_project[project]"),
		rego.SetRegoVersion(s.regoVersion),
	)

	query, err := r2.PrepareForEval(ctx)
	if err != nil {
		s.log.Errorf("failed to prepare for eval: %v", err)
		return errors.Wrap(err, "prepare query for eval (authorized_project)")
	}

	s.preparedEvalProjects = query

	return nil
}

func (s *State) newCompiler() (*ast.Compiler, error) {
	compiler := ast.NewCompiler()
	compiler.Compile(s.modules)
	if compiler.Failed() {
		s.log.Errorf("failed to compile modules: %v", compiler.Errors)
		return nil, errors.Wrap(compiler.Errors, "compile modules")
	}

	return compiler, nil
}

func (s *State) DumpData(ctx context.Context) error {
	return s.dumpData(ctx, s.store)
}

func (s *State) dumpData(ctx context.Context, store storage.Store) error {
	txn, err := store.NewTransaction(ctx)
	if err != nil {
		return err
	}
	data, err := store.Read(ctx, txn, []string{})
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	s.log.Info("data: ", string(jsonData))

	return store.Commit(ctx, txn)
}

func (s *State) evalQuery(ctx context.Context, query ast.Body, input interface{}, store storage.Store) (rego.ResultSet, error) {
	var tracer *topdown.BufferTracer
	if s.enableQueryTracer {
		tracer = topdown.NewBufferTracer()
	}

	rs, err := rego.New(
		rego.ParsedQuery(query),
		rego.Input(input),
		rego.Compiler(s.compiler),
		rego.Store(store),
		rego.QueryTracer(tracer),
		rego.SetRegoVersion(s.regoVersion),
	).Eval(ctx)
	if err != nil {
		s.log.Errorf("failed to evaluate query: %v", err)
		return nil, err
	}

	if tracer != nil && tracer.Enabled() {
		var buffer bytes.Buffer
		topdown.PrettyTrace(&buffer, *tracer)
		s.log.Debug(buffer.String())
	}

	return rs, nil
}

func (s *State) pairsFromResults(rs rego.ResultSet) (engine.Pairs, error) {
	pairs := make(engine.Pairs, len(rs))
	for i, r := range rs {
		if len(r.Expressions) != 1 {
			return nil, &UnexpectedResultExpressionError{exps: r.Expressions}
		}
		m, ok := r.Expressions[0].Value.(map[string]interface{})
		if !ok {
			return nil, &UnexpectedResultExpressionError{exps: r.Expressions}
		}
		res, ok := m["resource"].(string)
		if !ok {
			return nil, &UnexpectedResultExpressionError{exps: r.Expressions}
		}
		act, ok := m["action"].(string)
		if !ok {
			return nil, &UnexpectedResultExpressionError{exps: r.Expressions}
		}
		pairs[i] = engine.Pair{Resource: engine.Resource(res), Action: engine.Action(act)}
	}

	return pairs, nil
}

func (s *State) pairsFromAllowed(rs rego.ResultSet) (bool, error) {
	for _, r := range rs {
		if len(r.Expressions) != 1 {
			return false, &UnexpectedResultExpressionError{exps: r.Expressions}
		}
		m, ok := r.Expressions[0].Value.(map[string]interface{})
		if !ok {
			return false, &UnexpectedResultExpressionError{exps: r.Expressions}
		}
		_, ok = m["resource"].(string)
		if !ok {
			return false, &UnexpectedResultExpressionError{exps: r.Expressions}
		}
		_, ok = m["action"].(string)
		if !ok {
			return false, &UnexpectedResultExpressionError{exps: r.Expressions}
		}
		return true, nil
	}

	return false, nil
}

func (s *State) projectsFromPartialResults(rs rego.ResultSet) (engine.Projects, error) {
	if len(rs) != 1 {
		return nil, &UnexpectedResultSetError{set: rs}
	}

	r := rs[0]
	if len(r.Expressions) != 1 {
		return nil, &UnexpectedResultExpressionError{exps: r.Expressions}
	}

	projects, err := s.stringArrayFromResults(r.Expressions)
	if err != nil {
		s.log.Errorf("failed to parse projects: %v", err)
		return nil, &UnexpectedResultExpressionError{exps: r.Expressions}
	}

	return projects, nil
}

func (s *State) stringArrayFromResults(exps []*rego.ExpressionValue) (engine.Projects, error) {
	rawArray, ok := exps[0].Value.([]interface{})
	if !ok {
		return nil, &UnexpectedResultExpressionError{exps: exps}
	}

	vals := make(engine.Projects, len(rawArray))
	var v string
	for i := range rawArray {
		if v, ok = rawArray[i].(string); !ok {
			return nil, errors.New("error casting to string")
		}
		vals[i] = engine.Project(v)
	}

	return vals, nil
}

func (s *State) projectsFromPreparedEvalQuery(rs rego.ResultSet) (engine.Projects, error) {
	projectsFound := make(map[string]bool, len(rs))
	result := make(engine.Projects, 0, len(rs))
	var ok bool
	var proj string
	for i := range rs {
		proj, ok = rs[i].Bindings["project"].(string)
		if !ok {
			return nil, &UnexpectedResultExpressionError{exps: rs[i].Expressions}
		}
		if !projectsFound[proj] {
			result = append(result, engine.Project(proj))
			projectsFound[proj] = true
		}
	}
	return result, nil
}

func (s *State) allowedFromPreparedEvalQuery(rs rego.ResultSet) (bool, error) {
	var ok bool
	for i := range rs {
		_, ok = rs[i].Bindings["project"].(string)
		if !ok {
			return false, &UnexpectedResultExpressionError{exps: rs[i].Expressions}
		}
		return true, nil
	}
	return false, nil
}
