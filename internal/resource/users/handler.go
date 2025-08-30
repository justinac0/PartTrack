package users

import (
	"PartTrack/internal/resource/sessions"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store *UserStore
}

func NewHandler() *Handler {
	return &Handler{
		store: NewStore(),
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// delete old session in db if exists, creates new one and sets browsers cookie
func recreateSession(c echo.Context, ctx context.Context, userId uint64) error {
	sessionStore := sessions.NewStore()
	err := sessionStore.Delete(ctx, userId)
	if err != nil {
		fmt.Println("[err]: ", err)
	}

	// TODO: write helper for session creation
	// create session
	expiry := time.Now().Add(time.Hour * 24).UTC()
	now := time.Now().UTC()

	sessionId := uuid.New()

	session, err := sessionStore.Create(ctx, sessions.Session{
		UserId:    userId,
		SessionId: sessionId.String(), // TODO: generate unique
		Expiry:    &expiry,
		Created:   &now,
	})
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	cookie := http.Cookie{
		Name:    "session",
		Value:   session.SessionId,
		Expires: *session.Expiry,
		// TODO: needs to support the following cookie attributes in prod
		// HttpOnly: true,
		// Secure: true,
	}

	c.SetCookie(&cookie)

	fmt.Println("created cookie")
	return nil
}

func (h *Handler) SignIn(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	user, err := h.store.GetByUsername(ctx, username)
	if err != nil {
		return c.String(http.StatusOK, "incorrect username/password")
	}

	if !verifyPassword(password, user.PasswordHash) {
		return c.String(http.StatusOK, "incorrect username/password")
	}

	err = recreateSession(c, ctx, user.Id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if user.Role == RoleAdmin {
		c.Response().Header().Set("HX-Redirect", "/admin")
		return c.NoContent(http.StatusOK)
	}

	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.NoContent(http.StatusOK)
}

func (h *Handler) SignOut(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	}

	cookie.Value = ""
	cookie.Expires = time.Now().UTC()

	c.SetCookie(cookie)

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func (h *Handler) Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")

	passHash, err := hashPassword(password)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	now := time.Now().UTC()
	user := User{
		Username:     username,
		Email:        email,
		PasswordHash: passHash,
		Role:         RoleGuest,
		CreatedAt:    &now,
	}

	_, err = h.store.Create(ctx, user)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.String(http.StatusCreated, "account created, you can signin now!")
}

func (h *Handler) GetUserById(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = h.store.GetOne(ctx, id)
	if err != nil {
		panic(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := h.store.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) WhoAmI(c echo.Context) error {
	if ValidateSession(c) != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusUnauthorized)
	}

	cookie, err := c.Cookie("session")
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusUnauthorized)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sessionStore := sessions.NewStore()
	session, err := sessionStore.GetBySessionId(ctx, cookie.Value)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusUnauthorized)
	}

	user, err := h.store.GetOne(ctx, session.UserId)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.String(http.StatusOK, fmt.Sprintf("%s [%s]", user.Username, strings.ToUpper(string(user.Role))))
}

func ValidateSession(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cookie, err := c.Cookie("session")
	if err != nil {
		return sessions.SessionCookieNotSet
	}

	sessionStore := sessions.NewStore()
	session, err := sessionStore.GetBySessionId(ctx, cookie.Value)
	if err != nil {
		return sessions.SessionNotFound
	}

	userStore := NewStore()
	user, err := userStore.GetOne(ctx, session.UserId)
	if err != nil {
		return err
	}

	fmt.Println(user.Id, session.UserId)
	if user.Id != session.UserId {
		return sessions.SessionNotFound
	}

	if cookie.Value != session.SessionId {
		return sessions.SessionIdInvalid
	}

	if time.Now().After(*session.Expiry) {
		return sessions.SessionExpired
	}

	return nil
}
