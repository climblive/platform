# E2E

## Install dependencies

```bash
pnpm exec playwright install --with-deps
```

## Run tests

```bash
make test
```

The admin test reads `ADMIN_USERNAME` and `ADMIN_PASSWORD` from the environment.
Use a confirmed Cognito test user with a permanent password and no MFA in the
pool configured in `packages/lib/src/config.json`. The test clicks Sign in in
the admin app, logs in through the Cognito hosted UI, and follows the redirect
back to the app. The Cognito app client must allow the test callback URL
(`https://localhost:8443/admin` by default, or `${BASE_URL}/admin`). The
competition data is created in the disposable test database.

In GitHub Actions, set repository secrets `E2E_ADMIN_USERNAME` and
`E2E_ADMIN_PASSWORD`. These are mapped to the environment variables above.
The admin test fails with a setup error if credentials are missing.

To run only the admin test after building the test image:

```bash
make images
pnpm exec playwright test --project=admin
```
