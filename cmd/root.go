package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Capricorn = &cobra.Command{
	Use:   "capricorn",
	Short: "Root command for Capricorn CLI Tools",
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("author", "Aaron Hill <68696c6c@gmail.com>")
	viper.SetDefault("license", "MIT")
}
