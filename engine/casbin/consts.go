package casbin

const DefaultWildcardItem = "*"

const DefaultAuthorizedProjectsMatcher = "g(r.sub, p.sub, p.dom) && (keyMatch(r.dom, p.dom) || p.dom == '*')"
