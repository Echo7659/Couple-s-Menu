package user_jwt

import (
	"errors"
	"git.echol.cn/loser/menu-server/model/app/request"
	"git.echol.cn/loser/menu-server/utils"
	"time"

	"git.echol.cn/loser/menu-server/global"
	jwt "github.com/golang-jwt/jwt/v5"
)

type UserJWT struct {
	SigningKey []byte
}

func NewUserJWT() *UserJWT {
	return &UserJWT{
		[]byte(global.GVA_CONFIG.JWT.SigningKey),
	}
}

func (j *UserJWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	bf, _ := utils.ParseDuration(global.GVA_CONFIG.JWT.BufferTime)
	ep, _ := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	claims := request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"LCKT"},                  // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    global.GVA_CONFIG.JWT.Issuer,              // 签名的发行者
		},
	}
	return claims
}

// CreateToken 创建一个token
func (j *UserJWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *UserJWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// ParseToken 解析 token
func (j *UserJWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			claims, ok := token.Claims.(*request.CustomClaims)
			if ok {
				return claims, nil
			}
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, utils.TokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, utils.TokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, utils.TokenNotValidYet
		default:
			return nil, utils.TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, utils.TokenValid
}
