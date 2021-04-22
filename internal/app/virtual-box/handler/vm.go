package handler

import (
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/config"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/request"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/response"
	"github.com/sirupsen/logrus"
	"github.com/terra-farm/go-virtualbox"
)

const (
	vmCommand    = "VBoxManage"
	modify       = "modifyvm"
	cpus         = "--cpus"
	memory       = "--memory"
	guestControl = "guestcontrol"
	run          = "run"
	exe          = "--exe"
	username     = "--username"
	password     = "--password"
	copyto       = "copyto"
	copyfrom     = "copyfrom"
	target       = "--target-directory"
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
	req := new(request.ChangeStateRequest)

	if err := c.Bind(req); err != nil {
		logrus.Infof("change state: failed to bind: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("change state: failed to validate: %s", err.Error())
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

	switch req.State {
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

	outputs := []string{}

	if req.CPUs > 0 {
		cmd := exec.Command(vmCommand, modify, req.ID, cpus, strconv.Itoa(int(req.CPUs)))
		logrus.Infof("executing %s", cmd)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		outputs = append(outputs, string(output))
	}

	if req.Memory > 0 {
		cmd := exec.Command(vmCommand, modify, req.ID, memory, strconv.Itoa(int(req.Memory)))

		logrus.Infof("executing %s", cmd)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		outputs = append(outputs, string(output))
	}

	return c.String(http.StatusOK, strings.Join(outputs, "\n"))
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

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role != "admin" && req.ID != "VM1" {
		return c.NoContent(http.StatusForbidden)
	}

	var cmd *exec.Cmd
	if req.Username != nil {
		cmd = exec.Command(vmCommand, guestControl, req.ID, run, exe, req.Command,
			username, *req.Username, password, *req.Password)
	} else {
		cmd = exec.Command(vmCommand, guestControl, req.ID, run, exe, req.Command)
	}

	logrus.Infof("executing %s", cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, string(output))
}

func (h *VMHandler) Upload(c echo.Context) error {
	req := new(request.UploadRequest)

	if err := c.Bind(req); err != nil {
		logrus.Infof("exec: failed to bind: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("exec: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role != "admin" && req.ID != "VM1" {
		return c.NoContent(http.StatusForbidden)
	}

	var cmd *exec.Cmd
	if req.Username != nil {
		cmd = exec.Command(vmCommand, guestControl, req.ID, copyto, target, req.DstPath, req.SrcPath,
			username, *req.Username, password, *req.Password)
	} else {
		cmd = exec.Command(vmCommand, guestControl, req.ID, copyto, target, req.DstPath, req.SrcPath)
	}

	logrus.Infof("executing %s", cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, string(output))
}

func (h *VMHandler) Transfer(c echo.Context) error {
	req := new(request.TransferRequest)

	if err := c.Bind(req); err != nil {
		logrus.Infof("exec: failed to bind: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		logrus.Infof("exec: failed to validate: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role != "admin" && req.SrcPath != "VM1" && req.DstPath != "VM1" {
		return c.NoContent(http.StatusForbidden)
	}

	path := strings.Split(req.SrcPath, "/")

	var cmd *exec.Cmd
	if req.SrcUsername != nil {
		cmd = exec.Command(vmCommand, guestControl, req.Src, copyfrom, target, "storage/"+path[len(path)-1], req.SrcPath,
			username, *req.SrcUsername, password, *req.SrcPassword)
	} else {
		cmd = exec.Command(vmCommand, guestControl, req.Src, copyfrom, target, "storage/"+path[len(path)-1], req.SrcPath)
	}

	logrus.Infof("transfering %s to host", path)
	if err := cmd.Run(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if req.SrcUsername != nil {
		cmd = exec.Command(vmCommand, guestControl, req.Dst, copyto,
			target, req.DstPath, "storage/"+path[len(path)-1], username, *req.DstUsername, password, *req.DstPassword)
	} else {
		cmd = exec.Command(vmCommand, guestControl, req.Dst, copyto,
			target, req.DstPath, "storage/"+path[len(path)-1])
	}

	logrus.Infof("transfering %s to guest", path)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, string(output))
}
