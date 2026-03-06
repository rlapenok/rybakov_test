package pg

import pkgErr "github.com/rlapenok/rybakov_test/pkg/errors"

// errPoolCreate - create new pool create error
func errPoolCreate(err error) *pkgErr.Error {
	meta := map[string]any{
		"original_error": err.Error(),
	}

	return pkgErr.NewError(
		pkgErr.CodeInternalError,
		"Failed to create PostgreSQL pool",
		"FAILED_TO_CREATE_POSTGRESQL_POOL",
		meta,
	)
}

// errPoolMigrate - create new pool migrate error
func errPoolMigrate(err error) *pkgErr.Error {
	meta := map[string]any{
		"original_error": err.Error(),
	}

	return pkgErr.NewError(
		pkgErr.CodeInternalError,
		"Failed to migrate PostgreSQL pool",
		"FAILED_TO_MIGRATE_POSTGRESQL_POOL",
		meta,
	)
}
