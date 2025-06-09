package opa

const (
	AuthzProjectsQueryKey    = "AuthzProjectsQuery"
	FilteredPairsQueryKey    = "FilteredPairsQuery"
	FilteredProjectsQueryKey = "FilteredProjectsQuery"
)

const (
	defaultAuthzProjectsQuery    = "data.authz.authorized_project[project]"
	defaultFilteredPairsQuery    = "data.authz.introspection.authorized_pair[_]"
	defaultFilteredProjectsQuery = "data.authz.introspection.authorized_project"
)
