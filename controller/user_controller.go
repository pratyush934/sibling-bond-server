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

	/* setting expired cookie to logout */

	http.SetCookie(w, &http.Cookie{
		Name:     "user-logged-in",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	_ = cjson.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

/*
GetProfile - Get user details
UpdateProfile - Update user details
ChangePassword - Allow users to change password
ForgotPassword - Send password reset email
ResetPassword - Process password reset token
*/

func GetProfile(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("userId").(string)

	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "User Id not found in context",
			InternalError: nil,
		})
	}

	userById, err := models.GetUserById(userId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "user not found",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, userById)

}

/*
	update profile doesn't make sense
*/

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)

	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "User Id not found in context",
			InternalError: nil,
		})
	}

	userById, err := models.GetUserById(userId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "User with Id not found in context",
			InternalError: nil,
		})
	}

	type newPassWord struct {
		Pass string `json:"pass"`
	}

	var newPass newPassWord

	if err := json.NewDecoder(r.Body).Decode(&newPass); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to read the newPassword",
			InternalError: err,
		})
	}

	if err := userById.ResetPassword(newPass.Pass); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able Reset the password",
			InternalError: nil,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, "Reset the PassWord!!")
}

func ForgotPassWord(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Unable to parse request",
			InternalError: err,
		})
	}

	// Validate input
	if request.Email == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Email is required",
			InternalError: nil,
		})
	}

	byEmail, err := models.GetUserByEmail(request.Email)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "user is not there via email",
			InternalError: err,
		})
	}
	token := byEmail.GeneratePasswordResetToken()

	user, err := models.UpdateUser(byEmail)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to save the user",
			InternalError: err,
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "reset-token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})
	user.PassWord = "kya-karoge-jaan-kar"

	_ = cjson.WriteJSON(w, http.StatusCreated, user)
}

func ResetPasswordFromToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("reset-token")

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadGateway,
			Message:       "Not able to get the cookie",
			InternalError: err,
		})
	}

	token := cookie.Value

	var EmailStruct struct {
		Password string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&EmailStruct); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadGateway,
			Message:       "Not able to get the EmailStruct",
			InternalError: err,
		})
	}

	if EmailStruct.Password == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusBadGateway,
			Message:       "Password is empty",
			InternalError: nil,
		})
	}

	userByToken, err := models.GetUserByResetToken(token)

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadGateway,
			Message:       "Not able to get the EmailStruct",
			InternalError: err,
		})
	}

	if err := userByToken.ResetPassword(EmailStruct.Password); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to reset the Password",
			InternalError: err,
		})
	}

	user, err := models.UpdateUser(userByToken)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to update the user",
			InternalError: err,
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "reset-token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})
	user.PassWord = ""
	_ = cjson.WriteJSON(w, http.StatusOK, user)
}
