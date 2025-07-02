package handler

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	utils "github.com/quantinium3/houshou/app/api/internal"
	"github.com/quantinium3/houshou/app/api/internal/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserHandler struct {
	DB *db.Queries
}

func NewUserHandler(db *db.Queries) *UserHandler {
	return &UserHandler{
		DB: db,
	}
}

type UserSignUpRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Password  string `json:"password" validate:"required,min=8,max=32"`
	Email     string `json:"email" validate:"required,email"`
	SubDomain string `json:"sub_domain" validate:"required,min=3,max=255"`
}

type SignUpResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	SubDomain string `json:"sub_domain"`
}

type UserSignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *UserHandler) SignUp(c echo.Context) error {
	req := new(UserSignUpRequest)

	if err := c.Bind(req); err != nil {
		c.Logger().Errorf("SignUp: failed to bind user request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest,
			"request body invalid")
	}

	if err := c.Validate(req); err != nil {
		c.Logger().Errorf("SignUp: failed to validate user request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to validate request")
	}

	userExists, err := h.DB.UserExists(c.Request().Context(), req.Email)
	if err != nil {
		c.Logger().Errorf("SignUp: Failed to fetch user data from database"+
			": %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if user exists")
	}
	if userExists {
		c.Logger().Errorf("SignUp: User already exists: %v", req.Email)
		return echo.NewHTTPError(http.StatusConflict, "User already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Errorf("SignUp: failed to generate password: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate password")
	}

	user, err := h.DB.CreateUser(c.Request().Context(), db.CreateUserParams{
		ID:        utils.CreateNamespacedId("user"),
		Name:      req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		Subdomain: req.SubDomain,
	})
	if err != nil {
		c.Logger().Errorf("SignUp: failed to create user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	// todo: implement email verification

	c.Logger().Printf("SignUp: User created: %v", user)
	return c.JSON(http.StatusCreated, SignUpResponse{
		Username:  user.Name,
		Email:     user.Email,
		SubDomain: user.Subdomain,
	})
}

func (h *UserHandler) SignIn(c echo.Context) error {
	req := new(UserSignInRequest)
	if err := c.Bind(req); err != nil {
		c.Logger().Errorf("SignIn: failed to bind user request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind user request")
	}

	if err := c.Validate(req); err != nil {
		c.Logger().Errorf("SignIn: failed to validate user request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to validate request")
	}

	userExists, err := h.DB.UserExists(c.Request().Context(), req.Email)
	if err != nil {
		c.Logger().Errorf("SignIn: Failed to fetch user data from database: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if user exists")
	}

	if !userExists {
		c.Logger().Errorf("SignIn: User does not exist: %v", req.Email)
		return echo.NewHTTPError(http.StatusUnauthorized, "User does not exist")
	}

	user, err := h.DB.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		c.Logger().Errorf("SignIn: Failed to fetch user data from database: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user data")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.Logger().Errorf("SignIn: Invalid password: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid password")
	}

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.Logger().Errorf("SignIn: Failed to get cookie: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get cookie")
	}

	refreshToken, err := h.DB.GetRefreshTokenByToken(c.Request().Context(),
		cookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Logger().Errorf("SignIn: Refresh token not found")
			err := h.DB.DeleteAllRefreshTokens(c.Request().Context(), user.ID)
			if err != nil {
				c.Logger().Errorf("SignIn: Failed to delete refresh token: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete refresh token")
			}
		} else {
			c.Logger().Errorf("SignIn: Failed to get refresh token from"+
				" database: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError,
				"Failed to get refresh token from database")
		}
	} else if refreshToken.Token != user.ID {
		c.Logger().Errorf("SignIn: Refresh token does not match")
		err := h.DB.DeleteAllRefreshTokens(c.Request().Context(), user.ID)
		if err != nil {
			c.Logger().Errorf("SignIn: Failed to delete refresh token: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete refresh token")
		}
	}

	if err = h.DB.DeleteAllRefreshTokens(c.Request().Context(),
		user.ID); err != nil {
		c.Logger().Errorf("SignIn: Failed to delete refresh token: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete refresh token")
	}

	accessToken, err := utils.CreateAccessToken(user.ID)
	if err != nil {
		c.Logger().Errorf("SignIn: Failed to create access token: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create access token")
	}

	refreshTokenNew, err := utils.CreateRefreshToken(user.ID)
	if err != nil {
		c.Logger().Errorf("SignIn: Failed to create refresh token: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create refresh token")
	}

	refreshTokenNewDB, err := h.DB.CreateRefreshTokens(c.Request().Context(),
		db.CreateRefreshTokensParams{
			ID:     utils.CreateNamespacedId("refrehtoken"),
			Token:  refreshTokenNew,
			UserID: user.ID,
			ExpiresAt: pgtype.Timestamp{
				Time: time.Now().Add(30 * 24 * time.Hour),
				// 30 days todo: change in future to use config
				Valid: true,
			},
			CreatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
		})
	if err != nil {
		c.Logger().Errorf("SignIn: Failed to create refresh token in database"+
			": %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			"Failed to create refresh token in database")
	}

	newCookie := new(http.Cookie)
	newCookie.Name = "refresh_token"
	newCookie.Value = refreshTokenNewDB.Token
	newCookie.Path = "/"
	newCookie.Secure = true
	newCookie.HttpOnly = true
	// add the config to global
	newCookie.Expires = time.Now().Add(30 * 24 * time.
		Hour) // 30 days todo: change in future to use config
	c.SetCookie(newCookie)

	return c.JSON(http.StatusOK, SignInResponse{
		AccessToken: accessToken,
	})
}

func (h *UserHandler) SignOut(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.Logger().Errorf("SignOut: Failed to get cookie: %v", err)
		return echo.NewHTTPError(http.StatusNoContent, "Failed to get cookie")
	}

	if cookie.Value == "" {
		c.Logger().Errorf("SignOut: Cookie value is empty")
		return echo.NewHTTPError(http.StatusNoContent, "Cookie value is empty")
	}

	refreshTokenExist, err := h.DB.RefreshTokenExists(c.Request().Context(), cookie.Value)
	if err != nil {
		c.Logger().Errorf("SignOut: Failed to check if refresh token exists: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if refresh token exists")
	}

	if !refreshTokenExist {
		newCookie := new(http.Cookie)
		newCookie.Name = "refresh_token"
		newCookie.Value = ""
		newCookie.Expires = time.Unix(0, 0)
		newCookie.HttpOnly = true
		newCookie.Secure = true
		newCookie.Path = "/"
		c.SetCookie(newCookie)
		return c.JSON(http.StatusNoContent, SignInResponse{})
	}

	if err := h.DB.DeleteAllRefreshTokens(c.Request().Context(),
		cookie.Value); err != nil {
		c.Logger().Errorf("SignOut: Failed to delete refresh token: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete refresh token")
	}

	newCookie := new(http.Cookie)
	newCookie.Name = "refresh_token"
	newCookie.Value = ""
	newCookie.Expires = time.Unix(0, 0)
	newCookie.HttpOnly = true
	newCookie.Secure = true
	newCookie.Path = "/"
	c.SetCookie(newCookie)
	return c.JSON(http.StatusNoContent, SignInResponse{})
}
func (h *UserHandler) RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.Logger().Errorf("RefreshToken: Failed to get cookie: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Failed to get cookie")
	}

	if cookie.Value == "" {
		c.Logger().Errorf("RefreshToken: Cookie value is empty")
		return echo.NewHTTPError(http.StatusUnauthorized, "Cookie value is empty")
	}

	refreshTokenExist, err := h.DB.RefreshTokenExists(c.Request().Context(), cookie.Value)
	if err != nil {
		c.Logger().Errorf("RefreshToken: Failed to check if refresh token exists: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if refresh token exists")
	}

	newCookie := new(http.Cookie)
	newCookie.Name = "refresh_token"
	newCookie.Value = ""
	newCookie.Expires = time.Unix(0, 0)
	newCookie.HttpOnly = true
	newCookie.Secure = true
	newCookie.Path = "/"
	c.SetCookie(newCookie)

	if !refreshTokenExist {

	}
}
