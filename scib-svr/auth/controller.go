package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"scib-svr/logging"
	"strings"
)

const (
	userCollection = "users"
	RequestUri = "/users"
)

type Controller struct {
	s *Service
	l logging.Logger
}

func NewController(service *Service, logger logging.Logger) *Controller {
	return &Controller{
		service, logger,
	}
}

func (c *Controller) Authenticate(h httprouter.Handle, rolesRequired []string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenString, err := request.OAuth2Extractor.ExtractToken(r)
		if err != nil {
			c.l.Error(r.Context(), "%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		c.l.Info(r.Context(), "%v", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("wrong signing method: %v", token.Header["alg"])
			} else {
				return []byte(os.Getenv("SECURITY_CLIENT_SECRET")), nil
			}
		})
		if err != nil {
			c.l.Error(r.Context(), "%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			roles := extractRoles(claims)
			if containsAny(roles, rolesRequired) {
				r.Header.Add("roles", strings.Join(roles, ","))
				h(w, r, ps)
				return
			} else {
				c.l.Error(r.Context(), "%v", fmt.Errorf("insufficient rights"))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			c.l.Error(r.Context(), "%v", fmt.Errorf("request from %s could not be authenticated", r.RemoteAddr))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}

func extractRoles(claims jwt.MapClaims) []string {
	var roles []string
	for _, claim := range claims["roles"].([]interface{}) {
		roles = append(roles, fmt.Sprint(claim))
	}
	return roles
}

func containsAny(roles []string, requiredRoles []string) bool {
	for _, r := range roles {
		for _, rr := range requiredRoles {
			if r == rr {
				return true
			}
		}
	}
	return false
}

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	password := ps.ByName("password")
	user, err := c.s.SignIn(r.Context(), username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, "Authentication failed.")
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"roles": user.Roles,
		})
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, token)
	}
}

func (c *Controller) Save(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user = &User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err == nil {
		createdCnt, err := c.s.Save(r.Context(), user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "internal server error: %v", err)
			return
		}
		if createdCnt > 0 {
			w.WriteHeader(http.StatusCreated)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users, err := c.s.queryUsers(r.Context(), bson.M{})
	if err != nil {
		_, _ = fmt.Fprintf(w, "an error occured: %+v", err.Error())
	} else {
		bytes, err := json.Marshal(users)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(bytes)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "an error occured: %+v", err.Error())
		}
	}
}
