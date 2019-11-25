// Copyright © 2019 Derek Bekoe
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
	"encoding/hex"
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
	uuid "github.com/gofrs/uuid"
	homedir "github.com/mitchellh/go-homedir"
	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/sha3"
)

const (
	configKeyNatsURL       = "NatsURL"
	configKeyNatsClusterID = "NatsClusterID"
	configKeyUseShortName  = "UseShortName"
	configKeyFingerprint   = "Fingerprint"
)

// Path to config file set by user
var cfgFile string

// Whether to log verbose output
var verbose bool

// Whether non-tls connection should be used for NATS connection
var useUnsecure bool

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

// positionalArgsValidator validates the positional args
func positionalArgsValidator(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	} else if len(args) == 1 {
		// We do not check if it's a valid uuid as it could be a short name also
		return nil
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

// Channel ID is what we use as the NATS Streaming subject (channel name is just the user facing name)
func getChannelID(channelName string) string {
	hash := make([]byte, 64)
	fingerprint := viper.GetString(configKeyFingerprint)
	if fingerprint == "" {
		if useUnsecure {
			log.Printf("Allowing no fingerprint as user specified --unsecure.")
		} else {
			errorExit("No keyfile fingerprint found - Use 'convey configure' to set the keyfile.")
		}
	} else {
		if !IsValidFingerprint(fingerprint) {
			errorExit(InvalidFingerprintMsg)
		}
		log.Printf("Using fingerprint to generate channel id")
	}
	inputBytes := []byte(fingerprint + channelName)
	sha3.ShakeSum256(hash, inputBytes)
	return hex.EncodeToString(hash)
}

func connectToStan(clientID string) (stan.Conn, *nats.Conn) {

	natsURL := viper.GetString(configKeyNatsURL)
	natsClusterID := viper.GetString(configKeyNatsClusterID)

	if natsURL == "" || natsClusterID == "" {
		s := fmt.Sprintf("The configuration options '%s' and '%s' are not set. Use `convey configure` to set. Use `--help` for usage.",
			configKeyNatsURL,
			configKeyNatsClusterID)
		errorExit(s)
	}

	natsSecureOpt := nats.Secure()

	if useUnsecure {
		s := fmt.Sprintf("Using unsecure connection to server - %s", natsURL)
		log.Printf(s)
		natsSecureOpt = nil
	}

	natsConn, err := nats.Connect(natsURL, natsSecureOpt)
	if err != nil {
		s := fmt.Sprintf("Failed to connect to NATS server due to error - %s", err)
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

	// Check if should use short names.
	if useShortName {
		channelName = createChannelNameShort()
	} else {
		channelName = createChannelNameUUID()
	}

	channelID := getChannelID(channelName)

	log.Printf("Using friendly channel name - %s\n", channelName)
	log.Printf("Publishing to channel id - %s\n", channelID)

	// Print channel to console for user to copy
	fmt.Println(channelName)

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
			stanConn.Publish(channelID, []byte(line))
		}
		donePublish <- true
	}()

	<-donePublish

	stanConn.Publish(channelID, etx)
	stanConn.Close()
	natsConn.Close()
}

// subscribeModeFunc handles reading messages from the message service
func subscribeModeFunc(channelName string) {
	clientID := getClientID("convey-sub")
	stanConn, natsConn := connectToStan(clientID)

	channelID := getChannelID(channelName)

	log.Printf("Using friendly channel name - %s\n", channelName)
	log.Printf("Subscribing to channel - id %s\n", channelID)

	doneSubscribe := make(chan bool)

	sub, subErr := stanConn.Subscribe(channelID, func(m *stan.Msg) {
		if reflect.DeepEqual(m.Data, etx) {
			doneSubscribe <- true
		} else {
			fmt.Println(string(m.Data))
		}
	}, stan.DeliverAllAvailable())

	if subErr != nil {
		s := fmt.Sprintf("Failed to subscribe to channel name %s due to error %s", channelName, subErr)
		errorExit(s)
	}

	<-doneSubscribe

	sub.Unsubscribe()
	stanConn.Close()
	natsConn.Close()
}
