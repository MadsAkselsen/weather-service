ENV ?= local
CONFIG_DIR ?= config

run:
ifeq ($(ENV),none)
	@echo "Running without specific environment configuration..."
	go run main.go
else
	@echo "Running with environment $(ENV)..."
	@set -a; . $(CONFIG_DIR)/$(ENV).env; set +a; go run main.go
endif
