package ikprovider

import (
	"encoding/json"
	"errors"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"sync"
)

var imageKitService *ServiceHandler
var once sync.Once
var initErr error

func validateFileName(fileName string) error {
	if len(fileName) == 0 || len(fileName) > 255 {
		return errors.New("file name should be between 1 and 255")
	}
	return nil
}

func GetImageKitService() (*ServiceHandler, error) {
	once.Do(func() {
		imageKitService, initErr = NewImageKitService()
		if initErr != nil {
			log.Error().Err(initErr).Msg("Failed to initialize ImageKit service")

		}
	})
	return imageKitService, initErr
}

type MultipleAuthRequest struct {
	FileNames []string `json:"fileNames"`
}

func GetImageKitAuthHandler(w http.ResponseWriter, r *http.Request) {

	imageKitService, err := GetImageKitService()
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "ImageKit service not available",
			InternalError: err,
		})
	}

	var multipleAuthRequest MultipleAuthRequest

	const maxRequestSize = 5 * (1 << 20)
	const maxFiles = 5
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	if err := json.NewDecoder(r.Body).Decode(&multipleAuthRequest); err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to decode the MultipleAuthRequest",
			InternalError: err,
		})
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to close request Body")
		}
	}(r.Body)

	if len(multipleAuthRequest.FileNames) == 0 || len(multipleAuthRequest.FileNames) > maxFiles {
		panic(&cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "File should be between 1 to 5 inclusively",
			InternalError: nil,
		})
	}

	var newQueue []map[string]interface{}

	for _, value := range multipleAuthRequest.FileNames {

		if err := validateFileName(value); err != nil {
			panic(&cjson.HTTPError{
				Status:        http.StatusBadRequest,
				Message:       "Invalid file name: " + value,
				InternalError: err,
			})
		}

		params, err := imageKitService.GenerateAuthenticationParams(value)
		if err != nil {
			panic(&cjson.HTTPError{
				Status:        http.StatusInternalServerError,
				Message:       "Not able to generateAuthenticationParam",
				InternalError: err,
			})
		}
		newQueue = append(newQueue, params)
	}

	_ = cjson.WriteJSON(w, http.StatusOK, newQueue)
}
