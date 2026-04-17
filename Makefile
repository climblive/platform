CERT_DIR := .local/certs
WEB_DIR := backend/cmd/api/web
APPS := admin scoreboard scorecard www

.PHONY: localhost-certs web build run stop clean

localhost-certs: $(CERT_DIR)/app-cert.pem $(CERT_DIR)/www-cert.pem

$(CERT_DIR)/app-cert.pem:
	mkdir -p $(CERT_DIR)
	openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 \
		-keyout $(CERT_DIR)/app-key.pem -out $(CERT_DIR)/app-cert.pem \
		-days 365 -nodes -subj "/CN=app.localhost" \
		-addext "subjectAltName=DNS:app.localhost"
	chmod 644 $(CERT_DIR)/app-key.pem

$(CERT_DIR)/www-cert.pem:
	mkdir -p $(CERT_DIR)
	openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 \
		-keyout $(CERT_DIR)/www-key.pem -out $(CERT_DIR)/www-cert.pem \
		-days 365 -nodes -subj "/CN=www.localhost" \
		-addext "subjectAltName=DNS:www.localhost"
	chmod 644 $(CERT_DIR)/www-key.pem

web:
	cd web && pnpm install
	cd web && sed 's/"API_URL":.*/"API_URL": "",/' packages/lib/src/config.json > packages/lib/src/config.json.tmp \
		&& mv packages/lib/src/config.json.tmp packages/lib/src/config.json
	cd web && pnpm --filter=* build
	$(foreach app,$(APPS),rm -rf $(WEB_DIR)/$(app) && cp -r web/$(app)/dist $(WEB_DIR)/$(app);)

build: web
	cd backend && CGO_ENABLED=0 go build -o climblive cmd/api/main.go

run: localhost-certs build
	docker compose up --build

stop:
	docker compose down

clean: stop
	docker compose down -v
	rm -rf .local backend/climblive
	$(foreach app,$(APPS),rm -rf $(WEB_DIR)/$(app);)
