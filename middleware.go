package auth

import (
	"net/http"

	"github.com/go-martini/martini"
)

func Keystone(validator TokenValidator, cache Cache) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		// Get token from headers
		token := req.Header.Get("X-Auth-Token")
		if token == "" {
			unauthorized(res)
		} else {
			// Get token from cache
			_, err := cache.Get(token)
			if err != nil {
				// Token is not in the cache -> proceed to validate
				if validator.ValidateToken(token) == false {
					// Invalid token -> fail request
					unauthorized(res)
				} else {
					// Valid token -> set in cache
					cache.Set(token, "authorized")
				}
			}

			// Add token to context
			c.Map(Token(token))
		}
	}
}

func unauthorized(res http.ResponseWriter) {
	http.Error(res, "Not Authorized", http.StatusUnauthorized)
}
