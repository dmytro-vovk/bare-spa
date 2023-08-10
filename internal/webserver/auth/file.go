package auth

type PlainAuth struct {
	userName string
	password string
}

func New(username, password string) *PlainAuth {
	return &PlainAuth{userName: username, password: password}
}

func (a *PlainAuth) Valid(username string, password string) bool {
	return username == a.userName && password == a.password
}
