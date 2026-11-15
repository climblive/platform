# Competition UUID migration

Migration 43 renames `contest` to `competition` and `contest_id` to
`competition_id`. Only competition identifiers become UUIDs, stored as
`BINARY(16)`. Existing competitions and new inserts use UUIDv7.

Before deploying, back up the database and stop all other API instances and
database writers. MySQL DDL commits implicitly, so this migration is not
transactional. Its persistent `competition_uuid_mapping` table allows the same
migration to resume after an interruption; do not remove it during recovery.
Foreign-key checks are connection-local and restored even on failure. Constraints
are recreated with checks enabled to validate existing relationships.

Downgrading requires restoring the backup. Numeric competition URLs/API IDs are
no longer supported; registration codes and all other entity IDs are unchanged.

Run integration tests against a disposable MySQL/MariaDB server:

```sh
cd backend
MIGRATION_TEST_DSN='user:password@tcp(localhost:3306)/' go test ./cmd/api/migrations
```

The test user must have permission to create and drop databases. Tests create
uniquely named databases and remove them afterward.
