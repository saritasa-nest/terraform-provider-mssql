help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "⚡ \033[34m%-30s\033[0m %s\n", $$1, $$2}'

init: build clear ## Build terraform provider and initialize example with just built binary
	terraform init

clear: ## Remove terraform lock file (with hash of binary) and modules/providers cache
	rm -rf .terraform && rm -f .terraform.lock.hcl

build: ## Build MSSQL terraform provider from source
	cd .. && make fmt && make install

rundb: ## Start 'mssql' container in Docker
	docker run -e "ACCEPT_EULA=Y" -e 'SA_PASSWORD=yourStrong(!)Password' --name mssql -p 1433:1433 -d mcr.microsoft.com/mssql/server:2019-latest