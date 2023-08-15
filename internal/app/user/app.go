package user

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type User struct {
	sessionManager *scs.SessionManager
}

func New(address string, sessionManager *scs.SessionManager) *User {
	return &User{
		sessionManager: sessionManager,
	}
}

func (u *User) Logout(ctx context.Context) error {
	if r, ok := ctx.Value("request").(*http.Request); ok {
		return u.sessionManager.Destroy(r.Context())
	}

	panic("no request provided")
}
