package repository

import (
	"errors"
	"testing"

	domainErrors "content-service/internal/domain/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapUniqueViolation(t *testing.T) {
	internalDBError := &pgconn.PgError{
		Code:           "23505",
		ConstraintName: "uq_content_galleries_app_slug",
		Message:        "sensitive database detail",
	}

	if got := mapUniqueViolation(internalDBError); !errors.Is(got, domainErrors.ErrDuplicateSlug) {
		t.Fatalf("mapped error = %v, want ErrDuplicateSlug", got)
	}

	otherError := errors.New("database unavailable")
	got := mapUniqueViolation(otherError)
	if errors.Is(got, domainErrors.ErrDuplicateSlug) {
		t.Fatalf("non-unique error mapped to duplicate slug: %v", got)
	}
	if got.Error() == domainErrors.ErrDuplicateSlug.Error() {
		t.Fatal("non-unique database error lost its wrapped context")
	}
}
