package auth

import (
	"andreasho/scalable-ecomm/pgk/errors"
	"andreasho/scalable-ecomm/pgk/rest"
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authorizationHeader := request.Header.Get("Authorization")
		if authorizationHeader == "" {
			rest.ErrorResponse(writer, 401, errors.Unauthorized)
			return
		}

		token := strings.TrimPrefix(authorizationHeader, "Bearer ")
		accessToken, err := parseAccessToken(token)
		if err != nil {
			rest.ErrorResponse(writer, 401, errors.Unauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), "claims", accessToken)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
