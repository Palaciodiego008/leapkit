package session

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

// ctxKey is the value used to store the session
// into the http.Request context.
var ctxKey sessionKey = "session"

// contextKey is the key type used to store the session
// into the http.Request context.
type sessionKey string

func New(secret, name string, options ...Option) *session {
	store := sessions.NewCookieStore([]byte(secret))

	// Run the options on the store
	for _, option := range options {
		option(store)
	}

	return &session{
		name:  name,
		store: store,
	}
}

type session struct {
	name  string
	store *sessions.CookieStore
}

// Register returns an *http.Request with the session set in its context and also
// a custom http.ResponseWriter implementation that will save the session after each HTTP call.
func (s *session) Register(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	session, _ := s.store.Get(r, s.name)

	// Look for a valuer in the context and set the values for flash
	// and session so that they can be used in other components of the request.
	vlr, ok := r.Context().Value("valuer").(interface{ Set(string, any) })
	if ok {
		vlr.Set("flash", flashHelper(session))
		vlr.Set("session", func() *sessions.Session { return session })
	}

	r = r.WithContext(context.WithValue(r.Context(), ctxKey, session))

	w = &saver{
		w:     w,
		req:   r,
		store: session,
	}

	return w, r
}
