package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFolder string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kpas",
	Short: "kpas is a CLI tool for generating ready to use Kubernetes clusters",
	Long: `kpas is a CLI for bootstrapping and managing Kubernetes clusters
including standard components such as Ingress, Automatic SSL, Git
hosting and CI.

kpas currently supports Hetzner Cloud, GCP and AWS natively, as well
as working with any VPS provider.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFolder, "config", "", "config folder (default is $HOME/.kpas2/)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func configBasePath() string {
	if cfgFolder != "" {
		return cfgFolder
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		return filepath.Join(home, ".kpas2")
	}
}

func repoPath() string {
	return filepath.Join(configBasePath(), "repos")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFolder != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFolder)
	} else {
		// Make sure the config directory exists
		if _, err := os.Stat(configBasePath()); os.IsNotExist(err) {
			os.Mkdir(configBasePath(), 0700)
		}

		// Make sure the repo folder exists
		if _, err := os.Stat(repoPath()); os.IsNotExist(err) {
			os.Mkdir(repoPath(), 0700)
		}

		// Set base path to the default and set the config file name to be
		viper.AddConfigPath(configBasePath())
		viper.SetConfigName("kpas.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
