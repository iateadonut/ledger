package swagger

import (
	"errors"

	"github.com/communitybridge/ledger/gen/restapi/operations/health"
	"github.com/communitybridge/ledger/gen/restapi/operations/transactions"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/go-openapi/runtime/middleware"
)

type codedResponse interface {
	Code() string
}

// ErrorResponse wraps the error in the api standard models.ErrorResponse object
func ErrorResponse(err error) *models.ErrorResponse {
	cd := ""
	if e, ok := err.(codedResponse); ok {
		cd = e.Code()
	}

	e := models.ErrorResponse{
		Code:    cd,
		Message: err.Error(),
	}
	return &e
}

var (
	ErrNotFound         = errors.New("not found")
	ErrNotValidCurrency = errors.New("asset not valid")
	ErrInvalid          = errors.New("invalid request")
	ErrDuplicate        = errors.New("duplicate resource")
)

// HealthErrorHandler handles error resp from calls to the health endpoint
func HealthErrorHandler(label string, err error) middleware.Responder {
	logrus.WithError(err).Error(label)

	return health.NewGetHealthBadRequest()

}

// TransactionErrorHandler handles
func TransactionErrorHandler(label string, err error) middleware.Responder {
	switch err.Error() {
	case ErrDuplicate.Error():
		return transactions.NewCreateTransactionConflict().WithPayload(ErrorResponse(err))
	default:
		return transactions.NewListTransactionsBadRequest()
	}
}
