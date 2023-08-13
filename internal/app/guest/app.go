package guest

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"net/http"
)

type Guest struct {
	sessionManager *scs.SessionManager
}

func New(sessionManager *scs.SessionManager) *Guest {
	return &Guest{
		sessionManager: sessionManager,
	}
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (g *Guest) Login(ctx context.Context, r *LoginRequest) (bool, error) {
	if req, ok := ctx.Value("request").(*http.Request); ok {
		g.sessionManager.Put(req.Context(), "role", "user")
		g.sessionManager.Put(req.Context(), "user_name", r.Name)

		return true, nil
	}
	return true, nil
}
