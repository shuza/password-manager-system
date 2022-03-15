package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "auth-server",
	Short: "auth-server",
	Long:  `auth-server`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	log.Info().Msgf("APP_PORT = %s", os.Getenv("APP_PORT"))
	log.Info().Msgf("DB_HOST = %s", os.Getenv("DB_HOST"))
	log.Info().Msgf("DB_PORT = %s", os.Getenv("DB_PORT"))
	log.Info().Msgf("DB_USER = %s", os.Getenv("DB_USER"))
	log.Info().Msgf("DB_PASSWORD = %s", os.Getenv("DB_PASSWORD"))
	log.Info().Msgf("DB_NAME = %s", os.Getenv("DB_NAME"))
}
