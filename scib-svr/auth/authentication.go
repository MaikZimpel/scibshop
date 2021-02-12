package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"scib-svr/logging"
	"strings"
)

const userCollection = "users"

func Authenticate(h httprouter.Handle, logger logging.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenString, err := request.OAuth2Extractor.ExtractToken(r)
		if err != nil {
			logger.Error(r.Context(), "%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Info(r.Context(), "%v", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("wrong signing method: %v", token.Header["alg"])
			} else {
				return []byte(os.Getenv("SECURITY_CLIENT_SECRET")), nil
			}
		})
		if err != nil {
			logger.Error(r.Context(), "%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			r.Header.Add("roles", strings.Join(func() []string {
				var roles []string
				for _, claim := range claims["roles"].([]interface{}) {
					roles = append(roles, fmt.Sprint(claim))
				}
				return roles
			}(), ","))
			h(w, r, ps)
			return
		} else {
			logger.Error(r.Context(), "%v", fmt.Errorf("request from %s could not be authenticated", r.RemoteAddr))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}