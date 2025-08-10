package ikprovider

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/imagekit-developer/imagekit-go"
	"os"
	"time"
)

type ServiceHandler struct {
	ik *imagekit.ImageKit `json:"ik"`
}

/*
	1. SetUp NewImageKitService
    2. GenerateAuthentication Parameters
*/

func NewImageKitService() (*ServiceHandler, error) {

	publicKey := os.Getenv("IMAGEKIT_PUBLIC_KEY")
	privateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	urlEndpoint := os.Getenv("IMAGEKIT_URL_ENDPOINT")

	if publicKey == "" || privateKey == "" || urlEndpoint == "" {
		return nil, fmt.Errorf("required ImageKit environment variables not set")
	}

	kit := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		UrlEndpoint: urlEndpoint,
	})

	return &ServiceHandler{kit}, nil
}

func (s *ServiceHandler) GenerateAuthenticationParams(file string) (map[string]interface{}, error) {
	token := uuid.New().String()
	expire := time.Now().Add(5 * time.Minute).Unix()

	signToken := s.ik.SignToken(imagekit.SignTokenParam{
		Token:   token,
		Expires: expire,
	})

	return map[string]interface{}{
		"signature": signToken.Signature,
		"expire":    signToken.Expires,
		"token":     signToken.Token,
	}, nil
}
