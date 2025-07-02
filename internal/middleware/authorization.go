package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dasunwickr/go-api/api"
	"github.com/dasunwickr/go-api/internal/tools"
	log "github.com/sirupsen/logrus"
)

var ErrUnauthorized = errors.New("invalid username or token")

func Authorization (next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username string = r.URL.Query().Get("username")
		var token = r.Header.Get("Authorization")
		var err error 

		if username == "" || token == "" {
			log.Error(ErrUnauthorized)
			api.RequestErrorHandler(w, ErrUnauthorized)
			return 
		}

		// Handle both "Bearer token" and "token" formats
		token = strings.TrimPrefix(token, "Bearer ")

		var database *tools.DatabaseInterface
		database, err = tools.NewDatabase()
		if err != nil {
			api.InternalErrorHandler(w)
			return
		}

		loginDetails := (*database).GetUserLoginDetails(username)

		if (loginDetails == nil || (token != (*loginDetails).AuthToken)){
			log.Error(ErrUnauthorized)
			api.RequestErrorHandler(w, ErrUnauthorized)
			return
		}

		next.ServeHTTP(w,r)
	})
}