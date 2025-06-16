package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"restapi/utils"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}
		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, utils.UnexpectedSigningMethodError
			}
			return []byte(jwtSecret), nil
		})
		if errors.Is(err, jwt.ErrTokenExpired) {
			http.Error(w, utils.TokenExpiredError.Error(), utils.TokenExpiredError.GetStatusCode())
			return
		} else if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		} else if err != nil {
			http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
			return
		}
		if !parsedToken.Valid {
			http.Error(w, utils.InvalidLoginTokenError.Error(), utils.InvalidLoginTokenError.GetStatusCode())
			return
		}
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println(err)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims["uid"])
		ctx = context.WithValue(ctx, "expiresAt", claims["exp"])
		ctx = context.WithValue(ctx, "username", claims["user"])
		ctx = context.WithValue(ctx, "role", claims["role"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
