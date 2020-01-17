package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:  "rqo",
	Long: `rqlite-operator CLI`,
}

var (
	home     = os.Getenv("HOME")
	username string
	password string

	// ConfigFile is $HOME/.rqo/config.json per default
	// contains user, password and url of zapper
	configFile string
)

func main() {
	addCommands()

	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username, faculzapperive if you have a "+home+"/.rqo/config.json file")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password, faculzapperive if you have a "+home+"/.rqo/config.json file")
	rootCmd.PersistentFlags().StringVarP(&configFile, "configFile", "c", home+"/.rqo/config.json", "configuration file, default is "+home+"/.rqo/config.json")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	log.SetLevel(log.DebugLevel)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

//AddCommands adds child commands to the root command rootCmd.
func addCommands() {
	rootCmd.AddCommand(cmdCreate)
}

var cmdCreate = &cobra.Command{
	Use:   "create <rqlite cluster name>",
	Short: "Create a rqlite cluster",
	Run: func(cmd *cobra.Command, args []string) {
		create(args)
	},
}

// create creates a message in specified topic
func create(name []string) {
	readConfig()
	log.Debugf("ID Message Created: %d", 1)
}

/**
func getClient() *zapper.Client {
	tc, err := zapper.NewClient(zapper.Options{
		URL:      viper.GetString("url"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
		Referer:  "rqo.v0",
	})

	if err != nil {
		log.Fatalf("Error while create new zapper Client: %s", err)
	}

	zapper.DebugLogFunc = log.Debugf
	return tc
}
*/

// readConfig reads config in .rqo/config per default
func readConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		viper.ReadInConfig() // Find and read the config file
	}
}
