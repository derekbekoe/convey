// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/docker/docker/pkg/namesgenerator"
	homedir "github.com/mitchellh/go-homedir"
	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyNatsURL       = "NatsURL"
	configKeyNatsClusterID = "NatsClusterID"
	configKeyUseShortName  = "UseShortName"
	demoNatsURL            = "tls://demo.nats.io:4443"
	demoNatsClusterID      = "convey-demo-cluster"
)

// Path to config file set by user
var cfgFile string

// Whether to log verbose output
var verbose bool

// Whether non-tls connection should be used for NATS connection
var useUnsecure bool

// Whether the demo server and mode should be used
var useDemoMode bool

// etx is an identifier for End Of Text Sequence
var etx = []byte{3}

func errorExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "convey",
	Short: "A command-line tool that makes sharing pipes between machines easy.",
	Long:  `A command-line tool that makes sharing pipes between machines easy.`,
	Args:  positionalArgsValidator,
	Run:   RootCommandFunc,
}

// RootCommandFunc is a handler for the bare application
func RootCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		publishModeFunc()
	} else if len(args) == 1 {
		subscribeModeFunc(args[0])
	} else {
		errorExit("Too many args")
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorExit(err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.convey.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&useUnsecure, "unsecure", false, "use unsecured connection (for development purposes only)")
	rootCmd.PersistentFlags().BoolVar(&useDemoMode, "demo", false, "use demo mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			errorExit(err.Error())
		}

		// Search config in home directory with name ".convey" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".convey")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Config file:", viper.ConfigFileUsed())
	}
}

// positionalArgsValidator valids the positional args
func positionalArgsValidator(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	} else if len(args) == 1 {
		_, err := uuid.FromString(args[0])
		return err
	}
	return errors.New("Invalid positional arguments")
}

func createChannelNameUUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		s := fmt.Sprintf("Failed to create channel name: %s\n", err)
		errorExit(s)
	}
	// Remove dashes from UUID to make copy-paste easier in terminal
	return strings.Replace(u.String(), "-", "", -1)
}

func createChannelNameShort() string {
	return namesgenerator.GetRandomName(0)
}

func getClientID(prefix string) string {
	u, err := uuid.NewV4()
	if err != nil {
		s := fmt.Sprintf("Failed to create client ID: %s\n", err)
		errorExit(s)
	}
	return fmt.Sprintf("%s-%s", prefix, u.String())
}

func connectToStan(clientID string) (stan.Conn, *nats.Conn) {

	natsURL := viper.GetString(configKeyNatsURL)
	natsClusterID := viper.GetString(configKeyNatsClusterID)

	if useDemoMode {
		log.Printf("Using demo mode")
		natsURL = demoNatsURL
		natsClusterID = demoNatsClusterID
	}

	if natsURL == "" || natsClusterID == "" {
		s := fmt.Sprintf("The configuration options '%s' and '%s' are not set. Use `convey configure` to set. Use `--help` for usage.",
			configKeyNatsURL,
			configKeyNatsClusterID)
		errorExit(s)
	}

	natsSecureOpt := nats.Secure()

	if useUnsecure {
		log.Printf("Using unsecure connection to server")
		natsSecureOpt = nil
	}

	natsConn, err1 := nats.Connect(natsURL, natsSecureOpt)
	if err1 != nil {
		s := fmt.Sprintf("Failed to connect to NATS server due to error - %s", err1)
		errorExit(s)
	}

	stanConn, err := stan.Connect(
		natsClusterID,
		clientID,
		stan.NatsConn(natsConn),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			log.Printf("Lost connection due to error - %s", err)
		}))
	if err != nil {
		s := fmt.Sprintf("Failed to connect to streaming server due to error - %s", err)
		errorExit(s)
	}

	return stanConn, natsConn
}

// publishModeFunc handles publishing messages to message service
func publishModeFunc() {
	clientID := getClientID("convey-pub")
	stanConn, natsConn := connectToStan(clientID)

	useShortName := viper.GetBool(configKeyUseShortName)
	channelName := ""

	if useDemoMode && useShortName {
		log.Printf("Short names not allowed in demo mode to prevent possible conflicts")
	}

	// Check if should use short names. Short names not allowed in demo mode to prevent possible conflicts.
	if useShortName && !useDemoMode {
		channelName = createChannelNameShort()
	} else {
		channelName = createChannelNameUUID()
	}

	// Print channel to console for user to copy
	fmt.Println(channelName)
	log.Printf("Publishing to channel %s\n", channelName)

	donePublish := make(chan bool)

	// Handle Ctrl+C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		donePublish <- true
	}()

	go func() {
		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		for scanner.Scan() {
			line := scanner.Text()
			stanConn.Publish(channelName, []byte(line))
		}
		donePublish <- true
	}()

	<-donePublish

	stanConn.Publish(channelName, etx)
	stanConn.Close()
	natsConn.Close()
}

// subscribeModeFunc handles reading messages from the message service
func subscribeModeFunc(channelName string) {
	clientID := getClientID("convey-sub")
	stanConn, natsConn := connectToStan(clientID)

	log.Printf("Subscribing to channel %s\n", channelName)

	doneSubscribe := make(chan bool)

	sub, subErr := stanConn.Subscribe(channelName, func(m *stan.Msg) {
		if reflect.DeepEqual(m.Data, etx) {
			doneSubscribe <- true
		} else {
			fmt.Println(string(m.Data))
		}
	}, stan.DeliverAllAvailable())

	if subErr != nil {
		s := fmt.Sprintf("Failed to subscribe to channel %s due to error %s", channelName, subErr)
		errorExit(s)
	}

	<-doneSubscribe

	sub.Unsubscribe()
	stanConn.Close()
	natsConn.Close()
}
