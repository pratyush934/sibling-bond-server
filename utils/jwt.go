package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pratyush934/sibling-bond-server/cjson"
	"github.com/pratyush934/sibling-bond-server/models"
	"net/http"
	"strings"
	"time"
)

var privateKey = []byte("iampratyushiampratyushiampratyushiampratyush")

/*
	1. CreateToken
	2. GetToken
	3. GetTokenFromHeader
	4. ValidateToken
*/

func CreateToken(u *models.User) (string, error) {
	ttl := 1800

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    u.Id,
		"email": u.Email,
		"name":  u.FirstName,
		"role":  u.RoleId,
		"st":    time.Now(),
		"et":    time.Now().Add(time.Second * time.Duration(ttl)).Unix(),
	})
	return claims.SignedString(privateKey)
}

func ValidateAdminRole(r *http.Request) {
	token := GetToken(r)
	claims, ok := token.Claims.(jwt.MapClaims)
	roleId := (uint)(claims["role"].(float64))
	if !ok || !token.Valid || roleId != 2 {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Ye Banda is not admin",
			InternalError: fmt.Errorf("issue issue issue in admin wala banda"),
		})
	}
}

func ValidateToken(r *http.Request) {
	token := GetToken(r)
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Token claim is not correct",
			InternalError: nil,
		})
	}
}

func GetToken(r *http.Request) *jwt.Token {
	header, err := GetTokenFromHeader(r)
	if err != nil {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Issue from GetToken",
			InternalError: err,
		})
	}
	parse, err := jwt.Parse(header, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return privateKey, err
	})
	return parse
}

func GetTokenFromHeader(r *http.Request) (string, error) {
	str := r.Header.Get("Authorization")
	split := strings.Split(str, " ")
	if len(split) != 2 || split[0] != "Bearer" {
		panic(&cjson.HTTPError{
			Status:        http.StatusUnauthorized,
			Message:       "Not getting any token or right token from GetTokenFromHeader",
			InternalError: fmt.Errorf("issue issue"),
		})
	}
	return split[1], nil

}
