package middlewares

import (
	"net/http"

	conf "github.com/hackclub/hackatime/config"
)

const (
	xForwardedProtoHeader = "x-forwarded-proto"
)

func ForceSsl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if conf.Get().Env == "heroku" {
			if r.Header.Get(xForwardedProtoHeader) != "https" {
				sslUrl := "https://" + r.Host + r.RequestURI
				http.Redirect(w, r, sslUrl, http.StatusTemporaryRedirect)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
