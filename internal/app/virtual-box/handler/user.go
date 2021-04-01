package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/config"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/model"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/request"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type UserHandler struct {
	Cfg      config.Config
	UserRepo model.UserRepo
}

func NewUserHandler(cfg config.Config, userRepo model.UserRepo) *UserHandler {
	return &UserHandler{
		Cfg:      cfg,
		UserRepo: userRepo,
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	req := new(request.UserRequest)

	if err := c.Bind(req); err != nil {
		logrus.Infof("register: failed to bind: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("register: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	user := model.NewUser(req.Username, req.Password, req.Role)

	if err := h.UserRepo.Save(user); err != nil {
		logrus.Infof("register: failed to save: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	token, err := h.generateJWT(*user)
	if err != nil {
		logrus.Infof("register: failed to generate jwt: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	req := request.UserRequest{}

	if err := c.Bind(&req); err != nil {
		logrus.Infof("login: failed to bind: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("login: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := h.UserRepo.Find(req.Username)
	if err != nil {
		logrus.Infof("login: failed to find: %s", err.Error())
		return c.NoContent(http.StatusForbidden)
	}

	if !user.CheckPassword(req.Password) {
		logrus.Infof("login: incorrect password: %s", err.Error())
		return c.NoContent(http.StatusForbidden)
	}

	token, err := h.generateJWT(user)
	if err != nil {
		logrus.Infof("login: failed to generate jwt: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *UserHandler) generateJWT(user model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["Username"] = user.Username
	claims["Role"] = user.Role
	claims["exp"] = time.Now().Add(h.Cfg.JWT.Expiration).Unix()

	return token.SignedString([]byte(h.Cfg.JWT.Secret))
}
