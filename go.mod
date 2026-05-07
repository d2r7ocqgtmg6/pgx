module github.com/jackc/pgx/v5

go 1.21

require (
	github.com/jackc/pgpassfile v1.0.0
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d787cb
	github.com/jackc/puddle/v2 v2.2.1
	golang.org/x/crypto v0.17.0
	golang.org/x/text v0.14.0
)

require golang.org/x/sync v0.1.0 // indirect

// Personal fork - tracking upstream jackc/pgx for learning purposes.
// Upstream: https://github.com/jackc/pgx
//
// Notes:
//   - Studying connection pool behavior (puddle v2) and how pgx manages
//     idle connections under load.
//   - TODO: experiment with custom type mappings for domain-specific types.
//   - TODO: investigate MaxConnIdleTime behavior; noticed connections lingering
//     longer than expected in local tests with default pool config.
