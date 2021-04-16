package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4/middleware"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/config"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/handler"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/model"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/router"
	"github.com/saman2000hoseini/virtual-box_management_server/pkg/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	e := router.New(cfg)

	myDB, err := database.FirstSetup()
	if err != nil {
		logrus.Fatalf("failed to setup db: %s", err.Error())
	}

	userRepo := model.SQLUserRepo{DB: myDB}
	userHandler := handler.NewUserHandler(cfg, userRepo)

	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	vmHandler := &handler.VMHandler{Cfg: cfg}

	vm := e.Group("/vm", middleware.JWT([]byte(cfg.JWT.Secret)))
	vm.GET("/status", vmHandler.GetAllStatus)
	vm.GET("/status/:id", vmHandler.GetStatus)
	vm.PUT("/alter", vmHandler.ChangeState)
	vm.PUT("/modify", vmHandler.Modify)
	vm.POST("/clone", vmHandler.Clone)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := e.Start(cfg.Server.Address); err != nil {
			logrus.Fatalf("failed to start virtual-box management server: %s", err.Error())
		}
	}()

	logrus.Info("virtual-box management server started!")

	s := <-sig

	logrus.Infof("signal %s received", s)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	e.Server.SetKeepAlivesEnabled(false)

	if err := e.Shutdown(ctx); err != nil {
		logrus.Errorf("failed to shutdown virtual-box management server: %s", err.Error())
	}
}

// Register registers server command for virtual-box binary.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run virtual-box server component",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
