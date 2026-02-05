package auth

import (
	"andreasho/scalable-ecomm/pgk/errors"
	"andreasho/scalable-ecomm/pgk/rest"
	"context"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Extract cookie and retrieve userID
		cookie, err := request.Cookie("access_token")
		if err != nil {
			rest.ErrorResponse(writer, 401, errors.Unauthorized)
			return
		}

		accessToken, err := parseAccessToken(cookie.Value)
		if err != nil {
			rest.ErrorResponse(writer, 401, errors.Unauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), "claims", accessToken)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
