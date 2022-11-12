package engine

type Type int

const (
	CasbinEngine Type = 1
	OpaEngine    Type = 2
)

type Subject string
type Subjects []Subject

func MakeSubjects(subs ...Subject) Subjects {
	return subs
}

type Project string
type Projects []Project

func MakeProjects(projects ...Project) Projects {
	return projects
}

type Action string
type Actions []Action

func MakeActions(actions ...Action) Actions {
	return actions
}

type Resource string
type Resources []Resource

func MakeResources(resources ...Resource) Resources {
	return resources
}

type Pair struct {
	Resource Resource `json:"resource"`
	Action   Action   `json:"action"`
}
type Pairs []Pair

func MakePair(res, act string) Pair {
	return Pair{Resource(res), Action(act)}
}
func MakePairs(pairs ...Pair) Pairs {
	return pairs
}
