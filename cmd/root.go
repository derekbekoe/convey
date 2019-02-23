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
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	stan "github.com/nats-io/go-nats-streaming"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

const natsURL = "nats://localhost:4223"
const clusterID = "test-cluster"

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

// TODO-DEREK Add E2E encryption - https://www.pubnub.com/developers/tech/security/aes-encryption/

// TODO-DEREK Figure out how to handle message ordering

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "convey",
	Short: "A command-line tool that makes sharing pipes between machines easy.",
	Long:  `A command-line tool that makes sharing pipes between machines easy.`,
	Args:  positionalArgsValidator,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: RootCommandFunc,
}

// RootCommandFunc is a handler for the bare application
func RootCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		PublishModeFunc()
	} else if len(args) == 1 {
		SubscribeModeFunc(args[0])
	} else {
		log.Fatal("Too many args")
	}
}

func createChannelName() string {
	u1, err := uuid.NewV1()
	if err != nil {
		s := fmt.Sprintf("Failed to create channel name: %s\n", err)
		log.Fatal(s)
	}
	// Remove dashes from UUID to make copy-paste easier in terminal
	return strings.Replace(u1.String(), "-", "", -1)
}

// PublishModeFunc handles publishing messages to message service
func PublishModeFunc() {
	clientID := "convey-pub-1"
	sc, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL(natsURL))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	channelName := createChannelName()

	fmt.Println(channelName)
	log.Printf("Channel name %s\n", channelName)

	// Wait a bit as we don't support history yet
	// time.Sleep(10 * time.Second)

	// Simple Synchronous Publisher
	log.Printf("About to send messages\n")
	// sc.Publish(channelName, []byte("Hello World")) // does not return until an ack has been received from NATS Streaming

	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	for scanner.Scan() {
		line := scanner.Text()
		// TODO-DEREK Publish multiple messages using the same channel instead of this current approach
		sc.Publish(channelName, []byte(line))
	}

	// TODO-DEREK On Ctrl+C, send EOF as well.

	// TODO-DEREK Send proper EOF message
	sc.Publish(channelName, []byte("EOF"))

	// Close connection
	sc.Close()
}

// SubscribeModeFunc handles reading messages from the message service
func SubscribeModeFunc(channelName string) {
	log.Println("Subscribing to channel ", channelName)

	clientID := "convey-sub-1"
	sc, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL(natsURL))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	donePublish := make(chan bool)

	// Simple Async Subscriber
	sub, _ := sc.Subscribe(channelName, func(m *stan.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
		fmt.Println(string(m.Data))
	}, stan.DeliverAllAvailable())

	<-donePublish

	// Unsubscribe
	sub.Unsubscribe()

	// Close connection
	sc.Close()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.convey.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
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
			fmt.Println(err)
			os.Exit(1)
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
