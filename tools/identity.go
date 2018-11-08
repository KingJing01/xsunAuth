package tools

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"xsunAuth/models"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	SecretKey = "sfljdsfjsljdslfdsfsdfjdsf"
)

func CheckLogin(token string) (result bool, claims jwt.MapClaims, err error) {
	//result, _ := models.GetUserByName(token)
	result, claims, err = checkToken(token)
	return result, claims, err
}

func checkToken(token string) (checkResult bool, claims jwt.MapClaims, err error) {
	var ttoken *jwt.Token
	ttoken, err = jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		fmt.Println("parase with claims failed.", err)
		return false, nil, err
	}
	claims, _ = ttoken.Claims.(jwt.MapClaims)

	iiat, _ := strconv.ParseInt(strconv.FormatFloat(claims["iat"].(float64), 'f', -1, 64), 10, 64)
	iat := time.Unix(iiat, 0).Format("2006-01-02 15:04:05")
	fmt.Println(iat)

	iexp, _ := strconv.ParseInt(strconv.FormatFloat(claims["exp"].(float64), 'f', -1, 64), 10, 64)
	exp := time.Unix(iexp, 0).Format("2006-01-02 15:04:05")
	fmt.Println(exp)

	return true, claims, nil
}

func CheckAuthority(stoken string, permissionName string) (result bool, claims jwt.MapClaims, err error) {
	var userid int64
	var token *jwt.Token
	token, err = jwt.Parse(stoken, func(*jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		fmt.Println("parase with claims failed.", err)
		return false, nil, err
	}
	claims, _ = token.Claims.(jwt.MapClaims)
	//userid, _ = strconv.ParseInt(claims["jti"].(string), 10, 64)
	tmp := strconv.FormatFloat(claims["jti"].(float64), 'f', -1, 64)
	userid, _ = strconv.ParseInt(tmp, 10, 64)
	_, err = models.GetPermissionByUserAndPermission(userid, permissionName)
	fmt.Println("userid.", userid)
	fmt.Println("permissionName.", permissionName)
	if err != nil {
		result = false
		err = errors.New("没有操作权限")
	} else {
		result = true
		err = nil
	}
	return result, claims, err

}
