package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"hotel-reservation/db"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrNotAuthorized()
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		if time.Now().Unix() > expires {
			return ErrTokenExpired()
		}

		userID := claims["userID"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrTokenExpired()
		}
		// Set the current authenticated user to the context
		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])

			return nil, ErrNotAuthorized()
		}

		secret := os.Getenv("JWT_SECRET")

		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, ErrNotAuthorized()
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, ErrNotAuthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrNotAuthorized()
	}

	return claims, nil
}
