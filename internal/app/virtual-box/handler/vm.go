package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/config"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/request"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/response"
	"github.com/sirupsen/logrus"
	"github.com/terra-farm/go-virtualbox"
)

type VMHandler struct {
	Cfg config.Config
}

func (h *VMHandler) GetAllStatus(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role != "admin" {
		return c.NoContent(http.StatusForbidden)
	}

	machines, err := virtualbox.ListMachines()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res := new(response.StatusResponse)

	for _, m := range machines {
		res.Status = append(res.Status, response.Status{
			Name:   m.Name,
			CPUs:   m.CPUs,
			Memory: m.Memory,
			State:  m.State,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *VMHandler) GetStatus(c echo.Context) error {
	id := c.Param("id")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role != "admin" && id != "VM1" {
		return c.NoContent(http.StatusForbidden)
	}

	machine, err := virtualbox.GetMachine(id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res := response.Status{
		Name:   id,
		CPUs:   machine.CPUs,
		Memory: machine.Memory,
		State:  machine.State,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *VMHandler) ChangeState(c echo.Context) error {
	id := c.QueryParam("id")
	state := c.QueryParam("state")
	if len(id) == 0 || len(state) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role != "admin" && id != "VM1" {
		return c.NoContent(http.StatusForbidden)
	}

	machine, err := virtualbox.GetMachine(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Machine not found")
	}

	switch state {
	case "start":
		err = machine.Start()
	case "stop":
		err = machine.Stop()
	case "restart":
		err = machine.Restart()
	case "reset":
		err = machine.Reset()
	case "refresh":
		err = machine.Refresh()
	case "delete":
		err = machine.Delete()
	case "poweroff":
		err = machine.Poweroff()
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *VMHandler) Modify(c echo.Context) error {
	req := new(request.ModifyRequest)

	if err := c.Bind(req); err != nil {
		logrus.Infof("clone: failed to bind: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("clone: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if req.CPUs+req.Memory <= 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role != "admin" && req.ID != "VM1" {
		return c.NoContent(http.StatusForbidden)
	}

	machine, err := virtualbox.GetMachine(req.ID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Machine not found")
	}

	if req.CPUs > 0 {
		machine.CPUs = req.CPUs
	}

	if req.Memory > 0 {
		machine.Memory = req.Memory
	}

	if err := machine.Modify(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *VMHandler) Clone(c echo.Context) error {
	req := new(request.CloneRequest)

	if err := c.Bind(req); err != nil {
		logrus.Infof("clone: failed to bind: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("clone: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role != "admin" && req.ID != "VM1" {
		return c.NoContent(http.StatusForbidden)
	}

	err := virtualbox.CloneMachine(req.ID, req.NewID, true)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *VMHandler) Exec(c echo.Context) error {
	req := new(request.ExecRequest)

	if err := c.Bind(req); err != nil {
		logrus.Infof("exec: failed to bind: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("exec: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	//exec.Command()

	return c.String(http.StatusOK, "")
}
