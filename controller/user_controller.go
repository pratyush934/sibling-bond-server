package controller

import (
	"encoding/json"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/dto"
	"github.com/pratyush934/sibling-bond-server/models"
	"github.com/pratyush934/sibling-bond-server/utils"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var register dto.RegisterModel

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Err(err).Msg("Issue in Closing r.Body in Register")
			return
		}
	}(r.Body)

	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusForbidden,
			Message:       "Not able to read the register",
			InternalError: err,
		})
	}
	user := models.User{
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Email:     register.Email,
		PassWord:  register.Password,
	}

	createUser, err := user.CreateUser()
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusForbidden,
			Message:       "Not able to CreateUser in Register",
			InternalError: err,
		})
	}
	//createUser.PassWord = "" just ignore this for now
	_ = cjson.WriteJSON(w, http.StatusCreated, createUser)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var login dto.LoginModel

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Err(err).Msg("Issue while closing the r.Body in Login")
			return
		}
	}(r.Body)

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusForbidden,
			Message:       "Not able to read the Login",
			InternalError: err,
		})
	}

	if login.Email == "" || login.PassWord == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Email and password are required",
			InternalError: nil,
		})
	}

	userByEmail, err := models.GetUserByEmail(login.Email)
	if err != nil || !userByEmail.ValidatePassWord(login.PassWord) {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Invalid Email or Password",
			InternalError: err,
		})
	}

	token, err := utils.CreateToken(userByEmail)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to Generate Token",
			InternalError: err,
		})
	}
	randomValue := cjson.CreateRandomToken()
	http.SetCookie(w, &http.Cookie{
		Name:     "user-logged-in",
		Value:    randomValue,
		Expires:  time.Now().Add(time.Hour * 1),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	_ = cjson.WriteJSON(w, http.StatusOK, dto.LoginResponse{User: *userByEmail, Token: token})
}

func LogOut(w http.ResponseWriter, r *http.Request) {

}
