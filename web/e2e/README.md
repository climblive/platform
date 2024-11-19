# E2E

## Install dependencies

```bash
pnpm exec playwright install --with-deps
```

## Run tests

```bash
DEBUG=testcontainers* PW_TEST_HTML_REPORT_OPEN='never' pnpm exec playwright test --project=chromium --headed
```
