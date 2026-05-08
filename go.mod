module github.com/jackc/pgx/v5

go 1.21

require (
	github.com/jackc/pgpassfile v1.0.0
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d787cb
	github.com/jackc/puddle/v2 v2.2.1
	golang.org/x/crypto v0.20.0
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
//   - NOTE: bumped golang.org/x/crypto to v0.20.0 locally to test if the
//     idle connection issue is related to TLS handshake timing; see
//     https://github.com/golang/go/issues/XXXXX for context.
//   - UPDATE: after testing with crypto v0.20.0, the lingering idle connection
//     issue persists — likely not TLS-related. Next step: trace puddle's
//     reaper goroutine to see if MaxConnIdleTime is being checked correctly.
//   - UPDATE 2: added some debug logging around puddle's reaper in a local
//     branch (debug/reaper-trace); confirmed MaxConnIdleTime IS being checked
//     but the ticker interval defaults to 1s — connections can linger up to
//     ~1s past MaxConnIdleTime. Probably acceptable; closing this thread.
