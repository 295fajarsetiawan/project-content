package handler

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// requestValidator is the single shared validator instance for all gRPC request DTOs.
var requestValidator = validator.New(validator.WithRequiredStructEnabled())

// validateRequest checks a request DTO against its `validate` struct tags and
// returns a gRPC InvalidArgument error listing the failing fields.
func validateRequest(v any) error {
	err := requestValidator.Struct(v)
	if err == nil {
		return nil
	}

	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return status.Error(codes.InvalidArgument, "invalid request")
	}

	fields := make([]string, 0, len(verrs))
	for _, fe := range verrs {
		fields = append(fields, fe.Field()+" failed on "+fe.Tag())
	}
	return status.Error(codes.InvalidArgument, "validation failed: "+strings.Join(fields, ", "))
}
