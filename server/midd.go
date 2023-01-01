package server

import (
	"MDServer/baseLoad"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Claims struct {
	Data string `json:"username"`
	jwt.StandardClaims
}

type TokenInterface struct {
	plugin []func(ctx *fiber.Ctx, token, origin string) error
	path   []string
}

func (token *TokenInterface) Use(fun func(ctx *fiber.Ctx, token, origin string) error) {
	token.plugin = append(token.plugin, fun)
}

func (token *TokenInterface) AddPath(path ...string) {
	token.path = append(token.path, path...)
}

func (token *TokenInterface) Handler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		relToken, err := Token(ctx)
		if err != nil {
			res, _ := Message.Error(fmt.Sprint("3:", err.Error()), ctx)
			return res
		} else {
			if len(relToken) > 0 {
				for i := 0; i < len(token.plugin); i++ {
					pluginFunc := token.plugin[i]
					err := pluginFunc(ctx, relToken, ctx.Get("token"))
					if err != nil {
						res, _ := Message.Error(err.Error(), ctx)
						return res
					}
				}
			}
		}
		ctx.Next()
		return nil
	}
}

func Token(ctx *fiber.Ctx) (string, error) {
	wantTo := ctx.Path()
	if UrlGuardsInclude(wantTo) {
		token := ctx.Get("token")
		if len(token) > 0 {
			res, err := ParsingToken(token)
			if err != nil {
				return "", errors.New("error token")
			} else {
				ctx.Set("relToken", res) //存储key为relToken共后续处理器使用
				ctx.Set("originToken", token)
				return res, nil
			}
		} else {
			return "", errors.New(fmt.Sprint("The api ", wantTo, " want to the correct verification token"))
		}
	}
	ctx.Next()
	return "", nil
}

func NewToken() *TokenInterface {
	return &TokenInterface{}
}

// MakeToken 制作一个过期时间为一周的token
func MakeToken(data string) (string, error) {
	nowTime := time.Now()
	failTime := nowTime.Add(7 * (24 * time.Hour))
	claims := Claims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: failTime.Unix(), // 失效时间 一周
			Issuer:    "ROOT",          // 签发人
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(baseLoad.TokenKey)
	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

// ParsingToken 解析token
func ParsingToken(token string) (string, error) {
	parseToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return baseLoad.TokenKey, nil
	})
	if parseToken != nil {
		if claims, ok := parseToken.Claims.(*Claims); ok && parseToken.Valid {
			return claims.Data, nil
		}
	}
	return "", err
}
