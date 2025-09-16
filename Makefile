# =====================
# Config
# =====================
MIGRATE = migrate
MIGRATIONS_DIR = ./cmd/migrate/migrations
DB_URL = "mysql://niga:123456789@tcp(127.0.0.1:3306)/myapp?parseTime=true"

# =====================
# Commands
# =====================

## Tạo migration mới: make migration name=add_users_table
migration:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## Apply all up migrations
up:
	$(MIGRATE) -source "file://$(MIGRATIONS_DIR)" -database $(DB_URL) up

## Rollback 1 step
down:
	$(MIGRATE) -source "file://$(MIGRATIONS_DIR)" -database $(DB_URL) down 1

## Rollback all
reset:
	$(MIGRATE) -source "file://$(MIGRATIONS_DIR)" -database $(DB_URL) down -all

## Drop DB (cẩn thận!)
drop:
	$(MIGRATE) -source "file://$(MIGRATIONS_DIR)" -database $(DB_URL) drop -f

## Force version (fix dirty)
force:
	$(MIGRATE) -source "file://$(MIGRATIONS_DIR)" -database $(DB_URL) force $(v)

## Check current version
version:
	$(MIGRATE) -source "file://$(MIGRATIONS_DIR)" -database $(DB_URL) version

# =====================
# Seed
# =====================

.PHONY: seed seed-window seed-docker

seed: ## default seed sẽ báo cách dùng
	@echo "👉 Hãy dùng: make seed-window hoặc make seed-docker"

seed-window:
	set DB_ADDR=niga:123456789@tcp(localhost:3306)/myapp?parseTime=true&& \
	go run ./cmd/migrate/seed/main.go

seed-docker:
	set DB_ADDR=niga:123456789@tcp(db:3306)/myapp?parseTime=true&& \
	go run ./cmd/migrate/seed/main.go


.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt