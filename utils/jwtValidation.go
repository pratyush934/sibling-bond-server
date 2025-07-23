package utils

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pratyush934/sibling-bond-server/models"
	"net/http"
)

func ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		token := GetToken(request)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(request.Context(), "userid", claims["id"])
			ctx = context.WithValue(ctx, "email", claims["email"])
			ctx = context.WithValue(ctx, "role", claims["role"])
			ctx = context.WithValue(ctx, "name", claims["name"])

			request = request.WithContext(ctx)
		} else {
			panic(models.HTTPError{
				Status:        http.StatusUnauthorized,
				Message:       "there is an issue while validating the Token, the call is from middleware",
				InternalError: fmt.Errorf("look at the ValidateUser Middleware, I think the User is not Validated"),
			})
		}

		next.ServeHTTP(writer, request)
	})
}

func ValidateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		token := GetToken(request)

		claims, ok := token.Claims.(jwt.MapClaims)
		role := (uint)(claims["role"].(float64))

		if ok && role == 2 {
			ctx := context.WithValue(request.Context(), "userid", claims["id"])
			ctx = context.WithValue(ctx, "email", claims["email"])
			ctx = context.WithValue(ctx, "role", claims["role"])
			ctx = context.WithValue(ctx, "name", claims["name"])

			request = request.WithContext(ctx)
		} else {
			panic(models.HTTPError{
				Status:        http.StatusUnauthorized,
				Message:       "there is token while validating the Token for admin role, the call is from ValidatedAdmin",
				InternalError: fmt.Errorf("look at the Validate Middleware, I think user is not admin"),
			})
		}

		next.ServeHTTP(writer, request)
	})
}

func ValidateTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		token := GetToken(request)

		claims, ok := token.Claims.(jwt.MapClaims)

	})
}
