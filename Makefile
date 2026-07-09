# =========================
# CONFIG
# =========================
APP_NAME := firstbeegoapi
VERSION ?= 1.0

# =========================
# START
# =========================
.PHONY: server
server:
	bee run

# =========================
# DOCKER
# =========================
.PHONY: docker-build
docker-build:
	docker build -t $(APP_NAME):$(VERSION) .

.PHONY: docker-up
docker-up:
	cd ./conf && docker compose run --rm migrate
	cd ./conf && docker compose up --build

.PHONY: docker-down
docker-down:
	cd ./conf && docker compose down
