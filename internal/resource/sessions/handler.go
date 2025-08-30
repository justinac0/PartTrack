package sessions

type Handler struct {
	store *SessionStore
}

func NewHandler() *Handler {
	return &Handler{
		store: NewStore(),
	}
}
