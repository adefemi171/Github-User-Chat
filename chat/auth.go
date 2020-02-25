package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
)

// The authyHandler not only implements the serveHTTP method but also
// stores(wraps) http.Handler in the `next` field
type authyHandler struct {
	next http.Handler
}

// ServeHTTP here will search for a special cookie called auth
// and it will use the Header and WriteHeader method on http.ResponseWriter
// to redirect the user to a login page only if the cookie is not missing
func (h *authyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// user not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Successfully Authenticated - call the next handler
	h.next.ServeHTTP(w, r)
}

// MustAuthy simply creates authyHandler that wraps any
// other http.Handler
func MustAuthy(handler http.Handler) http.Handler {
	return &authyHandler{next: handler}
}

// loginHandler handles the third-party login process.
// format: /authy/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// break the path into segments using strings.Split
	// before pulling out the value for action and provider
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		// getting provider object that matches the object specified in the URL
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		// GetBeginAuthURL gets the location where we send users in order to start authorization.
		// GetBeginAuthURL(nil, nil) are for state map and map of additional options
		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting GetBeginAuthURL %s:%s", provider, err), http.StatusInternalServerError)
			return
		}
		// redirect users browser to the returned URL if code gets no error from GetBeginAuthURL
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		// log.Println("TODO handle login for ", provider)
	// When the authentication provider redirects the users after the must have granted permission
	// the  URL specifies that its a callback action.
	// we parse RawQuery from the request into objx.Map and the CompleteAuth uses the value to complete the OAuth2
	// provider handshake
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		// getting credentials
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error trying to complete auth for provider %s: %s", provider, err), http.StatusInternalServerError)
			return
		}
		// getting the user
		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error trying to get user %s: %s", provider, err), http.StatusInternalServerError)
			return
		}
		// using authCookieValue to save some data
		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
