.PHONY: images
images:
	docker build .. -t climblive-web:latest
	docker build ../../backend -t climblive-api:latest

.PHONY: test
test: images
	PW_TEST_HTML_REPORT_OPEN='never' pnpm exec playwright test --headed --project="Mobile Chrome"

.PHONY: test-ci
test-ci: images
	DEBUG=testcontainers* pnpm exec playwright test