CERT_DIR := .local/certs
WEB_DIR := ../../backend/cmd/api/web

.PHONY: images certs web clean test test-ci

certs: $(CERT_DIR)/cert.pem

$(CERT_DIR)/cert.pem:
	mkdir -p $(CERT_DIR)
	openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 \
		-keyout $(CERT_DIR)/key.pem -out $(CERT_DIR)/cert.pem \
		-days 1 -nodes -subj "/CN=localhost"
	chmod 644 $(CERT_DIR)/key.pem

web:
	cd .. && pnpm i
	cd ../packages/lib && cp src/config.json src/config.json.bak \
		&& sed 's/"API_URL":.*/"API_URL": "",/' src/config.json.bak > src/config.json
	cd ../.. && pnpm --filter=* build
	cd ../packages/lib && mv src/config.json.bak src/config.json
	rm -rf $(WEB_DIR)/admin $(WEB_DIR)/scoreboard $(WEB_DIR)/scorecard $(WEB_DIR)/www
	cp -r ../admin/dist $(WEB_DIR)/admin
	cp -r ../scoreboard/dist $(WEB_DIR)/scoreboard
	cp -r ../scorecard/dist $(WEB_DIR)/scorecard
	cp -r ../www/dist $(WEB_DIR)/www

images: certs web
	cd ../../backend && CGO_ENABLED=0 go build -o ./climblive ./cmd/api
	docker build ../../backend -t climblive-api:latest

test: images
	PW_TEST_HTML_REPORT_OPEN='never' pnpm exec playwright test --project="Mobile Chrome" --ui --ui-host 0.0.0.0 --ui-port 9323

test-ci: images
	DEBUG=testcontainers* pnpm exec playwright test

clean:
	rm -rf .local
	rm -rf $(WEB_DIR)/admin $(WEB_DIR)/scoreboard $(WEB_DIR)/scorecard $(WEB_DIR)/www
