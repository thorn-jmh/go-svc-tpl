package cmd

import (
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "go-svc-tpl/api"
    "go-svc-tpl/internal/dao"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "go-svc-tpl",
    Short: "a simple go service template",
    Long: `This is a simple go service template,
        which is used to build a go service quickly.`,

    RunE: func(cmd *cobra.Command, args []string) error {

        // set log level
        if viper.GetString("App.RunLevel") == "debug" {
            logrus.SetLevel(logrus.DebugLevel)
        }

        // init db
        dao.InitDB()
        // start server
        return api.StartServer()
    },
}
