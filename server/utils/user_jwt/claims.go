package user_jwt

import (
	"git.echol.cn/loser/menu-server/model/app"
	"git.echol.cn/loser/menu-server/model/app/request"
	"net"
	"time"

	"git.echol.cn/loser/menu-server/global"
	"github.com/gin-gonic/gin"
)

func ClearToken(c *gin.Context) {
	// 增加cookie Authorization 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("Authorization", "", -1, "/", "", false, false)
	} else {
		c.SetCookie("Authorization", "", -1, "/", host, false, false)
	}
}

func SetToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie Authorization 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("Authorization", token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie("Authorization", token, maxAge, "/", host, false, false)
	}
}

func GetToken(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		j := NewUserJWT()
		token, _ = c.Cookie("Authorization")
		claims, err := j.ParseToken(token)
		if err != nil {
			global.GVA_LOG.Error("重新写入cookie token失败,未能成功解析token,请检查请求头是否存在Authorization且claims是否为规定结构")
			return token
		}
		SetToken(c, token, int((claims.ExpiresAt.Unix()-time.Now().Unix())/60))
	}
	return token
}

func GetClaims(c *gin.Context) (*request.CustomClaims, error) {
	token := GetToken(c)
	j := NewUserJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GVA_LOG.Error(err.Error())
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {

	if cl, err := GetClaims(c); err != nil {
		return 0
	} else {
		return cl.BaseClaims.ID
	}

}

// GetNickName 从Gin的Context中获取从jwt解析出来的用户名
func GetNickName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.NickName
		}
	} else {
		waitUse := claims.(*request.CustomClaims)
		return waitUse.NickName
	}
}

func LoginToken(user app.User) (token string, claims request.CustomClaims, err error) {
	j := NewUserJWT()
	claims = j.CreateClaims(request.BaseClaims{
		ID:       user.ID,
		NickName: user.NickName,
		Phone:    user.Phone,
		OpenId:   user.OpenId,
	})
	token, err = j.CreateToken(claims)
	return
}
