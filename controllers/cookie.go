package controllers

import (
	"net/http"
	"time"

	"google.golang.org/appengine"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util/passenger"
	"golang.org/x/net/context"
)

func init() {
	router.Handle("/cookie", ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
		if r.Method == "PUT" {
			return putCookie(ctx, w, r)
		}

		if r.Method == "DELETE" {
			return deleteCookie(ctx, w, r)
		}

		w.Header().Set("Allow", "PUT, DELETE")
		return http.StatusMethodNotAllowed, nil
	}))
}

func putCookie(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	p, ok := passenger.FromContext(ctx)
	if !ok {
		return http.StatusUnauthorized, nil
	}

	token := &model.Token{
		Description: "Login from " + r.RemoteAddr,
	}

	value, err := p.IssueToken(ctx, token)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    value,
		Secure:   !appengine.IsDevAppServer(),
		HttpOnly: true,
		Expires:  token.Expiry,
	})

	w.Write([]byte("OK"))
	return http.StatusOK, nil
}

func deleteCookie(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "deleted",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Time{},
	})

	w.Write([]byte("OK"))
	return http.StatusOK, nil
}
