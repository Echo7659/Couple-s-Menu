package request

import "github.com/golang-jwt/jwt/v5"

type BaseClaims struct {
	NickName string `json:"nickName"`
	ID       uint   `json:"id"`
	Phone    string `json:"phone"`
	OpenId   string `json:"openId"`
}

type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}
