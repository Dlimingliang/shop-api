package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UserId   uint64 `json:"userId"`
	UserName string `json:"userName"`
	RoleId   int    `json:"roleId"`
	jwt.StandardClaims
}
