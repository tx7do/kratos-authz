package assets

import (
	_ "embed"
)

//go:embed rbac.conf
var DefaultRbacModel string

//go:embed rbac_with_domains.conf
var DefaultRbacWithDomainModel string

//go:embed abac.conf
var DefaultAbacModel string

//go:embed acl.conf
var DefaultAclModel string

//go:embed restfull.conf
var DefaultRestfullModel string

//go:embed restfull_with_role.conf
var DefaultRestfullWithRoleModel string
