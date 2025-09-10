package apperror

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

var (
	// ------------------------
	// Generic errors
	// ------------------------
	ErrInternalServer = errors.New("internal server error") // 500
	ErrUnknown        = errors.New("unknown error")         // 500
	ErrTimeout        = errors.New("timeout")               // 504
	ErrUnauthorized   = errors.New("unauthorized")          // 401
	ErrForbidden      = errors.New("forbidden")             // 403
	ErrNotImplemented = errors.New("not implemented")       // 501

	// ------------------------
	// MongoDB / persistence errors
	// ------------------------
	// Map "record not found" to Mongo's ErrNoDocuments.
	ErrRecordNotFound = mongo.ErrNoDocuments // 404
	// You can raise this explicitly in your domain; detection also checks the real Mongo error codes.
	ErrDuplicatedKey = errors.New("duplicated key not allowed") // 409

	// ------------------------
	// Validation errors
	// ------------------------
	ErrInvalidData          = errors.New("invalid data")                                        // 400
	ErrInvalidID            = errors.New("invalid id")                                          // 400
	ErrRequiredField        = errors.New("required field missing")                              // 400
	ErrInvalidFormat        = errors.New("invalid format")                                      // 400
	ErrOutOfRange           = errors.New("value out of range")                                  // 400
	ErrUnprocessable        = errors.New("unprocessable entity")                                // 422
	ErrInvalidValue         = errors.New("invalid value, should be pointer to struct or slice") // 400
	ErrInvalidValueOfLength = errors.New("invalid association values, length doesn't match")    // 400
	ErrInvalidField         = errors.New("invalid field")                                       // 400

	// ------------------------
	// Business logic / domain-specific errors
	// ------------------------
	ErrAlreadyExists   = errors.New("already exists")   // 409
	ErrNotAvailable    = errors.New("not available")    // 409
	ErrLimitExceeded   = errors.New("limit exceeded")   // 429
	ErrOperationDenied = errors.New("operation denied") // 403

	// ------------------------
	// Other errors
	// ------------------------
	ErrConflict         = errors.New("conflict")            // 409
	ErrDependencyFail   = errors.New("dependency failure")  // 502
	ErrTransactionAbort = errors.New("transaction aborted") // 500
)

// StatusCode maps errors to Fiber HTTP status codes
func StatusCode(err error) int {
	switch {
	// Generic
	case errors.Is(err, ErrInternalServer), errors.Is(err, ErrUnknown), errors.Is(err, ErrTransactionAbort):
		return fiber.StatusInternalServerError
	case errors.Is(err, ErrTimeout), errors.Is(err, context.DeadlineExceeded):
		return fiber.StatusGatewayTimeout
	case errors.Is(err, context.Canceled):
		return fiber.StatusRequestTimeout
	case errors.Is(err, ErrUnauthorized):
		return fiber.StatusUnauthorized
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrOperationDenied):
		return fiber.StatusForbidden
	case errors.Is(err, ErrNotImplemented):
		return fiber.StatusNotImplemented

	// MongoDB / persistence errors
	case errors.Is(err, ErrRecordNotFound), errors.Is(err, mongo.ErrNoDocuments):
		return fiber.StatusNotFound
	case isMongoDuplicateKey(err), errors.Is(err, ErrDuplicatedKey), errors.Is(err, ErrConflict), errors.Is(err, ErrAlreadyExists), errors.Is(err, ErrNotAvailable):
		return fiber.StatusConflict
	case errors.Is(err, ErrDependencyFail):
		return fiber.StatusBadGateway

	// Validation / business logic
	case errors.Is(err, ErrInvalidData), errors.Is(err, ErrInvalidID), errors.Is(err, ErrRequiredField),
		errors.Is(err, ErrInvalidFormat), errors.Is(err, ErrOutOfRange), errors.Is(err, ErrInvalidValue),
		errors.Is(err, ErrInvalidValueOfLength), errors.Is(err, ErrInvalidField):
		return fiber.StatusBadRequest
	case errors.Is(err, ErrUnprocessable):
		return fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrLimitExceeded):
		return fiber.StatusTooManyRequests

	// Default
	default:
		return fiber.StatusInternalServerError
	}
}

// GRPCCode maps errors to gRPC status codes
func GRPCCode(err error) codes.Code {
	switch {
	// Generic
	case errors.Is(err, ErrInternalServer), errors.Is(err, ErrUnknown), errors.Is(err, ErrTransactionAbort):
		return codes.Internal
	case errors.Is(err, ErrTimeout), errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, ErrUnauthorized):
		return codes.Unauthenticated
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrOperationDenied):
		return codes.PermissionDenied
	case errors.Is(err, ErrNotImplemented):
		return codes.Unimplemented

	// MongoDB / persistence errors
	case errors.Is(err, ErrRecordNotFound), errors.Is(err, mongo.ErrNoDocuments):
		return codes.NotFound
	case isMongoDuplicateKey(err), errors.Is(err, ErrDuplicatedKey), errors.Is(err, ErrConflict), errors.Is(err, ErrAlreadyExists), errors.Is(err, ErrNotAvailable):
		return codes.AlreadyExists
	case errors.Is(err, ErrDependencyFail):
		return codes.FailedPrecondition

	// Validation / business logic
	case errors.Is(err, ErrInvalidData), errors.Is(err, ErrInvalidID), errors.Is(err, ErrRequiredField),
		errors.Is(err, ErrInvalidFormat), errors.Is(err, ErrOutOfRange), errors.Is(err, ErrInvalidValue),
		errors.Is(err, ErrInvalidValueOfLength), errors.Is(err, ErrInvalidField):
		return codes.InvalidArgument
	case errors.Is(err, ErrUnprocessable):
		return codes.FailedPrecondition
	case errors.Is(err, ErrLimitExceeded):
		return codes.ResourceExhausted

	// Default
	default:
		return codes.Unknown
	}
}

// isMongoDuplicateKey returns true if the error corresponds to a MongoDB duplicate key violation.
func isMongoDuplicateKey(err error) bool {
	if err == nil {
		return false
	}

	// Prefer the helper if available in the driver version.
	// This will return true for code 11000/11001 and related duplicate key scenarios.
	if mongo.IsDuplicateKeyError(err) {
		return true
	}

	// Fallback checks for older driver versions or wrapped errors.

	// WriteException
	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 || e.Code == 11001 || e.Code == 12582 {
				return true
			}
		}
		if we.WriteConcernError != nil {
			if we.WriteConcernError.Code == 11000 || we.WriteConcernError.Code == 11001 || we.WriteConcernError.Code == 12582 {
				return true
			}
		}
	}

	// BulkWriteException
	var bwe mongo.BulkWriteException
	if errors.As(err, &bwe) {
		for _, e := range bwe.WriteErrors {
			if e.Code == 11000 || e.Code == 11001 || e.Code == 12582 {
				return true
			}
		}
	}

	// CommandError
	var ce mongo.CommandError
	if errors.As(err, &ce) {
		if ce.Code == 11000 || ce.Code == 11001 || ce.Code == 12582 {
			return true
		}
	}

	return false
}
