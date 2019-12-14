//SetMiddlewareJSON: This will format all responses to JSON
//SetMiddlewareAuthentication: This will check for the validity of the authentication token provided.

package middlewares

import (
	"errors"
	"net/http"

	"github.com/diegoclair/ApiGolang/api/auth"
	"github.com/diegoclair/ApiGolang/api/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
