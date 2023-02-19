package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//自定义一个字符串
var jwtkey = []byte("eryajf.net")
var str string

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func main() {
	r := gin.Default()
	r.GET("/set", setting)
	r.GET("/get", getting)
	//监听端口默认为8080
	r.Run(":8080")
}

//颁发token
func setting(ctx *gin.Context) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: 2,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",  // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println(token)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		fmt.Println(err)
	}
	// str = tokenString
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "token": tokenString})
}

//解析token
func getting(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization") //从header中取Authorization这个key，再与刚刚的token进行比对
	// 先判断取到的是否为空，为空则跳出
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		ctx.Abort()
		return
	}
	// 再来解析token，解析失败则跳出
	token, claims, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		ctx.Abort()
		return
	}
	// 最后成功了
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "id": claims.UserId, "msg": "认证通过"})

}

// ParseToken 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtkey, nil
	})
	return token, Claims, err
}
