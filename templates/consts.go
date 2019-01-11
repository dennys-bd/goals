package templates

func initConsts() {
	Templates["consts"] = `package lib

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	// ContextKeyAuth const for Context authorization
	ContextKeyAuth = contextKey("Authorization")
)
`
}
