package utils

import (
	"encoding/json"
	"github.com/pratyush934/sibling-bond-server/models"
	"log"
	"net/http"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					writer.Header().Set("Content-Type", "application/json")

					if httpError := r.(*models.HTTPError); httpError != nil {

						writer.WriteHeader(httpError.Status)

						if httpError.InternalError != nil {
							log.Printf("Expected Error exist : Status : %v, Message : %v, InternalError : %v", httpError.Status, httpError.Message, httpError.InternalError)
						} else {
							log.Printf("Expected Error exist : Status : %v, Message : %v", httpError.Status, httpError.Message)
						}

						_ = json.NewEncoder(writer).Encode(models.ErrorResponse{
							Status:        httpError.Status,
							Message:       httpError.Message,
							InternalError: httpError.InternalError,
						})

					} else {

						writer.WriteHeader(http.StatusInternalServerError)

						log.Printf("Unexpected panic which was not expected %v", r)

						_ = json.NewEncoder(writer).Encode(models.ErrorResponse{
							Status:        http.StatusInternalServerError,
							Message:       "Unexpected Error has happened and we need to fix this",
							InternalError: nil,
						})
					}
				}
			}()
			next.ServeHTTP(writer, request)
		},
	)
}
