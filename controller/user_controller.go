package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/dto"
	"github.com/pratyush934/sibling-bond-server/models"
	"github.com/pratyush934/sibling-bond-server/utils"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
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

/*
GetAddresses - List user addresses
AddAddress - Create new user address
UpdateAddress - Modify existing address
DeleteAddress - Remove address
*/

func GetAddresses(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)

	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to get the Id",
			InternalError: nil,
		})
	}

	addressById, err := models.GetAddressByUserId(userId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to get the Id",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, addressById)
}

func AddAddress(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)

	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to get the Id",
			InternalError: nil,
		})
	}

	var addressModel dto.AddressModel
	if err := json.NewDecoder(r.Body).Decode(&addressModel); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusForbidden,
			Message:       "Not able to get the AddressModel",
			InternalError: err,
		})
	}
	realAddress := models.Address{
		UserId:     userId,
		StreetName: addressModel.StreetName,
		LandMark:   addressModel.LandMark,
		ZipCode:    addressModel.ZipCode,
		City:       addressModel.City,
		State:      addressModel.State,
	}

	create, err := realAddress.Create()
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to create the address",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, create)

}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (from JWT token)
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to get the user ID",
			InternalError: nil,
		})
	}

	// Parse the request body
	var updateRequest struct {
		AddressID  string `json:"addressId"`
		StreetName string `json:"streetName"`
		LandMark   string `json:"landMark"`
		ZipCode    string `json:"zipCode"`
		City       string `json:"city"`
		State      string `json:"state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to parse address update request",
			InternalError: err,
		})
	}

	// Ensure address ID is provided
	if updateRequest.AddressID == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Address ID is required",
			InternalError: nil,
		})
	}

	// Retrieve the existing address
	existingAddress, err := models.GetAddressById(updateRequest.AddressID)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Address not found",
			InternalError: err,
		})
	}

	// Security check: Verify address belongs to current user
	if existingAddress.UserId != userId {
		panic(cjson.HTTPError{
			Status:        http.StatusForbidden,
			Message:       "You don't have permission to update this address",
			InternalError: nil,
		})
	}

	// Update the address fields
	existingAddress.StreetName = updateRequest.StreetName
	existingAddress.LandMark = updateRequest.LandMark
	existingAddress.ZipCode = updateRequest.ZipCode
	existingAddress.City = updateRequest.City
	existingAddress.State = updateRequest.State

	// Save the updated address
	updatedAddress, err := models.UpdateAddress(existingAddress)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusInternalServerError,
			Message:       "Not able to update the address",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, updatedAddress)
}

func DeleteAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressId := vars["id"]

	err := models.DeleteAddress(addressId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Fail to delete",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, "Address deleted successfully")
}

/*
Orders & History
GetOrderHistory - List user's past orders
GetOrderDetails - Get specific order details
*/

func GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to get the user ID",
			InternalError: nil,
		})
	}
	page := 1
	pageSize := 10

	// Parse pagination parameters if provided
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if parsedSize, err := strconv.Atoi(sizeStr); err == nil && parsedSize > 0 {
			pageSize = parsedSize
		}
	}

	ordersById, err := models.GetOrdersByUserId(userId, page, pageSize)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the Orders",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, ordersById)
}

func GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to get the user ID",
			InternalError: nil,
		})
	}
	vars := mux.Vars(r)
	orderId := vars["id"]

	if orderId == "" {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Order Id is required",
			InternalError: nil,
		})
	}

	orderById, err := models.GetOrderByUserIdAndOrderId(userId, orderId)

	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to get the Orders",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, orderById)
}

/*
Admin Operations (if applicable)
GetAllUsers - List all users (admin only)
GetUserById - Get specific user details
DeleteUser - Remove user account
UpdateUserRole - Change user role/permissions
*/

func GetAllUsersByAdmin(w http.ResponseWriter, r *http.Request) {
	roleId, okk := r.Context().Value("role").(int)

	if !okk {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to extract the id and role from the context",
			InternalError: nil,
		})
	}

	if roleId != 2 {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "you are not admin",
			InternalError: nil,
		})
	}

	users, err := models.GetAllUsers(5, 10)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Not able to getAll the Users",
			InternalError: err,
		})
	}

	_ = cjson.WriteJSON(w, http.StatusOK, users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	roleId, okk := r.Context().Value("role").(int)
	vars := mux.Vars(r)
	userId := vars["id"]

	if !okk {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to extract the id and role from the context",
			InternalError: nil,
		})
	}

	if roleId != 2 {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "you are not admin",
			InternalError: nil,
		})
	}

	userById, err := models.GetUserById(userId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusNotFound,
			Message:       "Message not found",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, userById)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {

	roleId, ok := r.Context().Value("role").(int)

	vars := mux.Vars(r)
	userId := vars["id"]

	if !ok {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not able to extract the id and role from the context",
			InternalError: nil,
		})
	}

	if roleId != 2 {
		panic(cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "you are not admin",
			InternalError: nil,
		})
	}

	err := models.DeleteUser(userId)
	if err != nil {
		panic(cjson.HTTPError{
			Status:        http.StatusBadRequest,
			Message:       "Not able to delete",
			InternalError: err,
		})
	}
	_ = cjson.WriteJSON(w, http.StatusOK, "User Deleted Successfully")
}
