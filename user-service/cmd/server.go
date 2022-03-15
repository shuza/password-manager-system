package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"user-service/internal/app/server"
	"user-service/internal/app/token"
	"user-service/internal/app/user"
	"user-service/internal/pkg/postgres"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Start server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := postgres.New(&postgres.Config{
			/*Host:     "postgresql",
			Port:     "5432",
			User:     "foobar",
			Password: "foobar",
			Name:     "learn-db",*/
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		})
		if err != nil {
			panic(err)
		}
		s := server.NewServer(os.Getenv("APP_PORT"),
			user.NewService(user.NewRepository(db)),
			token.NewService(),
		)

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sig
			if err := s.Shutdown(); err != nil {
				log.Error().Err(err).Msg("error during server shutdown")
			}
		}()

		return s.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
