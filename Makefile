APP_VERSION=v0.0.1

.PHONY: tag
tag:
	git tag -f $(APP_VERSION) && git tag -f engine/opa/$(APP_VERSION) && git tag -f engine/noop/$(APP_VERSION) && git tag -f engine/casbin/$(APP_VERSION) && git tag -f authz/$(APP_VERSION) && git push --tags --force
