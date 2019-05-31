package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

func populateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			fmt.Println("coming in this block")
			// check if username, password are user's credentials
			email, err := fetchEmailForUser(username, password)
			if err == nil && email != "" {
				ctx := r.Context()
				ctx = context.WithValue(ctx, emailKey, email)
				ctx = context.WithValue(ctx, categoryKey, "users")
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			Log.Errorln(fmt.Errorf(`cannot fetch email for "username"="%s" and "password"="%s"`, username, password))
			msg := map[string]interface{}{
				"message": "Unauthorized ok",
				"status":  http.StatusUnauthorized,
			}
			writeBack(w, msg, nil)
			return
		}
		w.WriteHeader(401)
		fmt.Fprintln(w, `{"message":"Not Logged In, No Token"}`)
	})
}
