package opa

// nolint:lll
//go:generate go-bindata -pkg $GOPACKAGE -o policy.bindata.go -ignore .*_test.rego -ignore Makefile -ignore README\.md policy/...

import (
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

type State struct {
	store                storage.Store
	queries              map[string]ast.Body
	compiler             *ast.Compiler
	modules              map[string]*ast.Module
	preparedEvalProjects rego.PreparedEvalQuery
}

const (
	authzProjectsQuery    = "data.authz.authorized_project[project]"
	filteredPairsQuery    = "data.authz.introspection.authorized_pair[_]"
	filteredProjectsQuery = "data.authz.introspection.authorized_project"
)

func New(_ context.Context, opts ...OptFunc) (*State, error) {
	authzProjectsQueryParsed, err := ast.ParseBody(authzProjectsQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "parse query %q", authzProjectsQuery)
	}

	filteredPairsQueryParsed, err := ast.ParseBody(filteredPairsQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "parse query %q", filteredPairsQuery)
	}

	filteredProjectsQueryParsed, err := ast.ParseBody(filteredProjectsQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "parse query %q", filteredProjectsQuery)
	}

	s := State{
		store: inmem.New(),
		queries: map[string]ast.Body{
			authzProjectsQuery:    authzProjectsQueryParsed,
			filteredPairsQuery:    filteredPairsQueryParsed,
			filteredProjectsQuery: filteredProjectsQueryParsed,
		},
	}

	for _, opt := range opts {
		opt(&s)
	}

	if err := s.initModules(); err != nil {
		return nil, errors.Wrap(err, "init OPA modules")
	}

	return &s, nil
}

func (s *State) initModules() error {
	if len(s.modules) == 0 {
		mods := map[string]*ast.Module{}
		for _, name := range AssetNames() {
			if !strings.HasSuffix(name, ".rego") {
				continue
			}
			parsed, err := ast.ParseModule(name, string(MustAsset(name)))
			if err != nil {
				return errors.Wrapf(err, "parse policy file %q", name)
			}
			mods[name] = parsed
		}
		s.modules = mods
	}

	compiler, err := s.newCompiler()
	if err != nil {
		return errors.Wrap(err, "init compiler")
	}
	s.compiler = compiler
	return nil
}

func (s *State) makeAuthorizedProjectPreparedQuery(ctx context.Context) error {
	compiler, err := s.newCompiler()
	if err != nil {
		return err
	}

	r := rego.New(
		rego.Store(s.store),
		rego.Compiler(compiler),
		rego.ParsedQuery(s.queries[authzProjectsQuery]),
		rego.DisableInlining([]string{
			"data.authz.denied_project",
		}),
	)

	pq, err := r.Partial(ctx)
	if err != nil {
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
		return compiler.Errors
	}

	r2 := rego.New(
		rego.Store(s.store),
		rego.Compiler(compiler),
		rego.Query("data.__partialauthz.authorized_project[project]"),
	)

	query, err := r2.PrepareForEval(ctx)
	if err != nil {
		return errors.Wrap(err, "prepare query for eval (authorized_project)")
	}

	s.preparedEvalProjects = query

	return nil
}

func (s *State) newCompiler() (*ast.Compiler, error) {
	compiler := ast.NewCompiler()
	compiler.Compile(s.modules)
	if compiler.Failed() {
		return nil, errors.Wrap(compiler.Errors, "compile modules")
	}

	return compiler, nil
}

func (s *State) DumpData(ctx context.Context) error {
	return dumpData(ctx, s.store)
}

func dumpData(ctx context.Context, store storage.Store) error {
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

	log.Info("data: ", string(jsonData))
	return store.Commit(ctx, txn)
}

func (s *State) ProjectsAuthorized(ctx context.Context, subjects engine.Subjects, action engine.Action, resource engine.Resource, projects engine.Projects) (engine.Projects, error) {
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
		return engine.Projects{}, &EvaluationError{e: err}
	}

	return s.projectsFromPreparedEvalQuery(resultSet)
}

func (s *State) FilterAuthorizedPairs(ctx context.Context, subjects engine.Subjects, pairs engine.Pairs) (engine.Pairs, error) {

	opaInput := map[string]interface{}{
		"subjects": subjects,
		"pairs":    pairs,
	}

	rs, err := s.evalQuery(ctx, s.queries[filteredPairsQuery], opaInput, s.store)
	if err != nil {
		return nil, &EvaluationError{e: err}
	}

	return s.pairsFromResults(rs)
}

func (s *State) FilterAuthorizedProjects(ctx context.Context, subjects engine.Subjects) (engine.Projects, error) {

	opaInput := map[string]interface{}{
		"subjects": subjects,
	}

	rs, err := s.evalQuery(ctx, s.queries[filteredProjectsQuery], opaInput, s.store)
	if err != nil {
		return nil, &EvaluationError{e: err}
	}

	return s.projectsFromPartialResults(rs)
}

func (s *State) evalQuery(ctx context.Context, query ast.Body, input interface{}, store storage.Store) (rego.ResultSet, error) {

	var tracer *topdown.BufferTracer

	rs, err := rego.New(
		rego.ParsedQuery(query),
		rego.Input(input),
		rego.Compiler(s.compiler),
		rego.Store(store),
		rego.QueryTracer(tracer),
	).Eval(ctx)
	if err != nil {
		return nil, err
	}

	if tracer.Enabled() {
		topdown.PrettyTrace(os.Stderr, *tracer) //nolint: govet // tracer can be nil only if tracer.Enabled() == false
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
	for i := range rawArray {
		v, ok := rawArray[i].(string)
		if !ok {
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

func (s *State) SetPolicies(ctx context.Context, policyMap map[string]interface{}, roleMap map[string]interface{}) error {
	s.store = inmem.NewFromObject(map[string]interface{}{
		"policies": policyMap,
		"roles":    roleMap,
	})

	return s.makeAuthorizedProjectPreparedQuery(ctx)
}

type UnexpectedResultExpressionError struct {
	exps []*rego.ExpressionValue
}

func (e *UnexpectedResultExpressionError) Error() string {
	return fmt.Sprintf("unexpected result expressions: %v", e.exps)
}

type UnexpectedResultSetError struct {
	set rego.ResultSet
}

func (e *UnexpectedResultSetError) Error() string {
	return fmt.Sprintf("unexpected result set: %v", e.set)
}

type EvaluationError struct {
	e error
}

func (e *EvaluationError) Error() string {
	return fmt.Sprintf("error in query evaluation: %s", e.e.Error())
}
