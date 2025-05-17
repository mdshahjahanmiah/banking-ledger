package transaction

import (
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/http"
	"github.com/mdshahjahanmiah/explore-go/logging"
)

func MakeHandler(ms Service, logger *logging.Logger) http.Endpoint {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(error.EncodeError),
	}

	depositHandler := kithttp.NewServer(
		makeDepositEndpoint(ms, logger),
		decodeDepositRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	withdrawHandler := kithttp.NewServer(
		makeWithdrawEndpoint(ms, logger),
		decodeWithdrawRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	auditHandler := kithttp.NewServer(
		makeAuditEndpoint(ms),
		decodeAuditRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	r := chi.NewRouter()

	r.Method("POST", "/accounts/deposit", depositHandler)
	r.Method("POST", "/accounts/withdraw", withdrawHandler)
	r.Method("GET", "/accounts/{id}/transactions", auditHandler)

	return http.Endpoint{Pattern: "/accounts/*", Handler: r}
}
