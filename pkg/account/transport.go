package account

import (
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/http"
)

func MakeHandler(ms Service) http.Endpoint {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(error.EncodeError),
	}

	postAccountHandler := kithttp.NewServer(
		makePostAccountEndpoint(ms),
		decodeCreateAccountRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	r := chi.NewRouter()

	r.Method("POST", "/accounts", postAccountHandler)

	return http.Endpoint{Pattern: "/accounts", Handler: r}
}
