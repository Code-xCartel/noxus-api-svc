package _api

import (
	"context"
	"fmt"
	"github.com/Code-xCartel/noxus-api-svc/config"
	"github.com/Code-xCartel/noxus-api-svc/types/auth"
	"github.com/Code-xCartel/noxus-api-svc/utils"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.HandlerFunc

var chain = middlewareChain(
	corsMiddleware,
	requestLoggerMiddleware,
	requireAuthentication,
)

func middlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

func corsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// If the request is OPTIONS (preflight request), respond with 200 OK
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	}
}

func requestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func requireAuthentication(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if utils.IsExcludedFromAuth(r.RequestURI) {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("authorization token required"))
			return
		}
		st := strings.Split(token, " ")
		tokenMethod := st[0]
		tokenSignature := st[1]
		if tokenMethod != config.Envs.RequestMethod {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("incorrect request method"))
			return
		}
		parsedToken, err := jwt.ParseWithClaims(tokenSignature, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Envs.JWTSecretKey), nil
		})
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		claims, ok := parsedToken.Claims.(*auth.CustomClaims)
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unknown claims, cannot proceeed"))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "claims", claims)))
	}
}
