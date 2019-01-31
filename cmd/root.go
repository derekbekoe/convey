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
	pubnub "github.com/pubnub/go"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

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

func getPubNub() *pubnub.PubNub {
	config := pubnub.NewConfig()

	subKey := viper.GetString("SubscribeKey")
	pubKey := viper.GetString("PublishKey")

	if subKey == "" || pubKey == "" {
		log.Fatal("PubNub subscription and publish keys are required.")
	}

	config.SubscribeKey = subKey
	config.PublishKey = pubKey

	pn := pubnub.NewPubNub(config)
	return pn
}

// SubscribeModeFunc handles reading messages from the message service
func SubscribeModeFunc(channelName string) {
	log.Println("Subscribed to channel ", channelName)

	pn := getPubNub()

	// TODO-DEREK Get history since channel began instead of last 100 only - https://support.pubnub.com/support/solutions/articles/14000043629-how-do-i-page-through-stored-messages-
	res, status, _ := pn.History().Channel(channelName).Execute()

	log.Println(status)

	for _, element := range res.Messages {
		if _, isEOF := element.Message.(map[string]interface{})["EOF"]; isEOF {
			os.Exit(0)
		}
		fmt.Println(element.Message.(map[string]interface{})["msg"])
	}

	listener := pubnub.NewListener()
	doneConnect := make(chan bool)
	doneSubscribe := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNDisconnectedCategory:
					// This event happens when radio / connectivity is lost
					log.Println("Messaging service lost connectivity")
				case pubnub.PNConnectedCategory:
					// Connect event. You can do stuff like publish, and know you'll get it.
					// Or just use the connected event to confirm you are subscribed for
					// UI / internal notifications, etc
					log.Println("Messaging service connected")
					doneConnect <- true
				case pubnub.PNReconnectedCategory:
					// Happens as part of our regular operation. This event happens when
					// radio / connectivity is lost, then regained.
					log.Println("Messaging service regained connectivity")
				}
			case message := <-listener.Message:
				// Handle new message stored in message.message
				if msg, ok := message.Message.(map[string]interface{}); ok {
					if _, isEOF := msg["EOF"]; isEOF {
						doneSubscribe <- true
					} else {
						fmt.Println(msg["msg"])
						log.Printf("Got message: %s\n", msg["msg"])
					}
				}
				log.Println("message.Message", message.Message)
				log.Println("message.Timetoken", message.Timetoken)

			case <-listener.Presence:
				// handle presence
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{channelName}).
		Execute()

	<-doneConnect
	<-doneSubscribe
}

// PublishModeFunc handles publishing messages to message service
func PublishModeFunc() {

	// TODO-DEREK Do not subscribe if all you're doing is publishing

	pn := getPubNub()

	listener := pubnub.NewListener()
	doneConnect := make(chan bool)
	donePublish := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNDisconnectedCategory:
					// This event happens when radio / connectivity is lost
					log.Println("Messaging service lost connectivity")
				case pubnub.PNConnectedCategory:
					// Connect event. You can do stuff like publish, and know you'll get it.
					// Or just use the connected event to confirm you are subscribed for
					// UI / internal notifications, etc
					log.Println("Messaging service connected")
					doneConnect <- true
				case pubnub.PNReconnectedCategory:
					// Happens as part of our regular operation. This event happens when
					// radio / connectivity is lost, then regained.
					log.Println("Messaging service regained connectivity")
				}
			case message := <-listener.Message:
				// Handle new message stored in message.message
				if message.Channel != "" {
					// Message has been received on channel group stored in
					// message.Channel
				} else {
					// Message has been received on channel stored in
					// message.Subscription
				}
				if msg, ok := message.Message.(map[string]interface{}); ok {
					if _, isEOF := msg["EOF"]; isEOF {
						donePublish <- true
					} else {
						// fmt.Println(msg["msg"])
						// log.Printf("Got message: %s\n", msg["msg"])
					}
				}
				log.Println("message.Message", message.Message)
				log.Println("message.Timetoken", message.Timetoken)

			case <-listener.Presence:
				// handle presence
			}
		}
	}()

	pn.AddListener(listener)

	channelName := createChannelName()

	fmt.Println(channelName)
	log.Printf("Channel name %s\n", channelName)

	pn.Subscribe().
		Channels([]string{channelName}).
		Execute()

	<-doneConnect

	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	for scanner.Scan() {
		line := scanner.Text()
		msg := map[string]interface{}{
			"msg": line,
		}
		// TODO-DEREK Publish multiple messages using the same channel instead of this current approach
		_, _, err := pn.Publish().Channel(channelName).Message(msg).Execute()
		if err != nil {
			// Request processing failed.
			// Handle message publish error
			log.Printf("Failed to publish message: %s\n", err)
		}
	}

	// TODO-DEREK On Ctrl+C, send EOF as well.

	doneMsg := map[string]interface{}{
		"EOF": true,
	}
	pn.Publish().Channel(channelName).Message(doneMsg).Execute()

	<-donePublish
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.convey.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
