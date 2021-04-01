package server

import (
	"context"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/config"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/handler"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/model"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/router"
	"github.com/saman2000hoseini/virtual-box_management_server/pkg/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
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
