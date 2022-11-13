APP_VERSION=v0.0.1

PACKAGE_LIST = engine/opa/ engine/noop/ engine/casbin/ authz/

.PHONY: tag
tag:
	git tag -f $(APP_VERSION) && $(foreach item, $(PACKAGE_LIST), git tag -f $(item)$(APP_VERSION) && ) git push --tags --force
