// Copyright © 2017 Igor Maculan <n3wtron@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/foofilers/confHub/server"
	"github.com/Sirupsen/logrus"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "confHub",
	Short: "Configuration Hub Daemon",
	Run: func(cmd *cobra.Command, args []string) {
		initLog(cmd)
		listen := cmd.Flag("listen").Value.String()
		if len(listen) == 0 {
			listen = viper.GetString("listen")
			if len(listen) == 0 {
				logrus.Fatal("[listen] parameter/configuration missing")
			}
		}
		server.Start(listen)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().StringP("listen", "b", "0.0.0.0:8080", "listen address (format [addr]:[port])")
	RootCmd.Flags().StringP("logLevel", "l", "info", "logging level [debug,info,warning,error]")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.confHub.yaml)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initLog(cmd *cobra.Command) {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	switch cmd.Flag("logLevel").Value.String(){
	case "debug" :
		logrus.SetLevel(logrus.DebugLevel)
	case "info" :
		logrus.SetLevel(logrus.InfoLevel)
	case "warning", "warn" :
		logrus.SetLevel(logrus.WarnLevel)
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}
	logrus.Debug("Log initialized")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".confHub") // name of config file (without extension)
	viper.AddConfigPath(".")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Debug("Using config file:", viper.ConfigFileUsed())
	} else {
		logrus.Fatal(err)
	}
}
