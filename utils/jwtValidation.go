package utils

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"net/http"
)

func ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		token := GetToken(request)

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(request.Context(), "userId", claims["id"])
			ctx = context.WithValue(ctx, "email", claims["email"])
			ctx = context.WithValue(ctx, "role", claims["role"])
			ctx = context.WithValue(ctx, "name", claims["name"])

			request = request.WithContext(ctx)
		} else {
			panic(&cjson.HTTPError{
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

		if ok && token.Valid {

			if role != 2 {
				panic(&cjson.HTTPError{
					Status:        http.StatusForbidden,
					Message:       "Admin access required",
					InternalError: fmt.Errorf("user role %d is not admin (required: 2)", role),
				})
			}

			ctx := context.WithValue(request.Context(), "userId", claims["id"])
			ctx = context.WithValue(ctx, "email", claims["email"])
			ctx = context.WithValue(ctx, "role", claims["role"])
			ctx = context.WithValue(ctx, "name", claims["name"])

			request = request.WithContext(ctx)
		} else {
			panic(&cjson.HTTPError{
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

		role := (uint)(claims["role"].(float64))

		if ok && token.Valid {

			if role != 3 {
				panic(&cjson.HTTPError{
					Status:        http.StatusForbidden,
					Message:       "Tenant access required",
					InternalError: fmt.Errorf("user role %d is not tenant (required: 3)", role),
				})
			}

			ctx := context.WithValue(request.Context(), "userId", claims["id"])
			ctx = context.WithValue(ctx, "email", claims["email"])
			ctx = context.WithValue(ctx, "role", claims["role"])
			ctx = context.WithValue(ctx, "name", claims["name"])

			request = request.WithContext(ctx)
		} else {
			panic(&cjson.HTTPError{
				Status:        http.StatusUnauthorized,
				Message:       "there is token while validating the Token for tenant role, the call is from ValidateTenant",
				InternalError: fmt.Errorf("look at the Validate Middleware, I think user is not tenant"),
			})
		}

		next.ServeHTTP(writer, request)

	})
}
