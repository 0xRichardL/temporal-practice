
# Makefile for managing docker-compose services

# Use .PHONY to ensure these targets run even if files with the same name exist. It also improves performance.
.PHONY: help up down clean build rebuild logs ps restart tctl \
        _build _rebuild _logs _restart _tctl

# Default target when no command is specified.
default: help

help:
	@echo "Usage: make [command] [arguments]"
	@echo ""
	@echo "Available commands:"
	@echo "  up        - Start all services in detached mode."
	@echo "  down      - Stop and remove containers."
	@echo "  clean     - Stop and remove containers, networks, and volumes."
	@echo "  build [s] - Build images. 's' is an optional service name. e.g., 'make build account'"
	@echo "  rebuild s - Rebuild and restart a specific service 's'. e.g., 'make rebuild account'"
	@echo "  logs [a]  - View output from containers. 'a' are optional docker-compose logs args. e.g., 'make logs -f'"
	@echo "  ps        - List containers."
	@echo "  restart [s]- Restart services. 's' is an optional service name. e.g., 'make restart account'"
	@echo "  tctl [a]  - Run tctl commands 'a'. e.g., 'make tctl workflow list'"

up:
	docker-compose up -d

down:
	docker-compose down

clean:
	docker-compose down --volumes

# Public targets that capture arguments and call private implementation targets
build rebuild restart logs tctl:
	@$(MAKE) _$(@) CMD_ARGS="$(filter-out $@,$(MAKECMDGOALS))"

_build:
	@echo "Building service(s): $(or $(CMD_ARGS), all)"
	docker-compose build $(CMD_ARGS)

_rebuild:
	@if [ -z "$(CMD_ARGS)" ]; then \
		echo "Error: Service name not provided. Usage: make rebuild <service-name>"; \
		exit 1; \
	fi
	@echo "Rebuilding and recreating service: $(CMD_ARGS)"
	docker-compose build $(CMD_ARGS)
	docker-compose up -d --no-deps $(CMD_ARGS)

_logs:
	docker-compose logs $(CMD_ARGS)

ps:
	docker-compose ps

_restart:
	@echo "Restarting service(s): $(or $(CMD_ARGS), all)"
	docker-compose restart $(CMD_ARGS)

_tctl:
	@if [ -z "$(CMD_ARGS)" ]; then \
		echo "Error: tctl arguments not provided. Usage: make tctl <tctl-args>"; \
		exit 1; \
	fi
	docker-compose exec temporal-admin-tools tctl $(CMD_ARGS)