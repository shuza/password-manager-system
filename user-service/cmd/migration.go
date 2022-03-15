package cmd

import (
	"os"
	postgres2 "user-service/internal/pkg/postgres"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Start DB migrations",
	Long:  `It will run migrations according to db/migrations scripts on Reporting database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return postgres2.RunDatabaseMigration(&postgres2.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		})
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Start DB rollback",
	Long:  "It will rollback one step",
	RunE: func(cmd *cobra.Command, args []string) error {
		return postgres2.RollbackLatestMigration(&postgres2.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		})
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
}
