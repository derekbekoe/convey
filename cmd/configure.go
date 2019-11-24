package cmd

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/sha3"
)

// The NATS URL passed in from command-line
var natsURL string

// The NATS cluster ID passed in from command-line
var natsClusterID string

// The URL or filepath to file
var keyFile string

// A known fingerprint in hex format
var knownFingerprintHex string

// Whether short channel names should be used instead of the standard uuid format
var useShortName bool

// Whether the current config file should be overwritten
var forceWrite bool

const shake256MinBytesRequired = 64
const fingerprintByteLength = 64

func generateFingerprint(keyFile string) string {
	hash := make([]byte, fingerprintByteLength)
	var inputBytes []byte
	if keyFile == "" {
		return ""
	}
	if strings.HasPrefix(keyFile, "http://") || strings.HasPrefix(keyFile, "https://") {
		// Read keyfile from url
		resp, err1 := http.Get(keyFile)
		if err1 != nil {
			errorExit(err1.Error())
		}
		defer resp.Body.Close()
		var err2 error
		inputBytes, err2 = ioutil.ReadAll(resp.Body)
		if err2 != nil {
			errorExit(err2.Error())
		}
	} else {
		// Read keyfile from local file
		var err1 error
		inputBytes, err1 = ioutil.ReadFile(keyFile)
		if err1 != nil {
			errorExit(err1.Error())
		}
	}

	// Generate the fingerprint
	inputBytesLen := len(inputBytes)
	if inputBytesLen < shake256MinBytesRequired {
		errMsg := fmt.Sprintf("Bad keyfile provided - At least %d bytes required but got %d byte(s).", shake256MinBytesRequired, inputBytesLen)
		errorExit(errMsg)
	}
	sha3.ShakeSum256(hash, inputBytes)
	return hex.EncodeToString(hash)
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.PersistentFlags().StringVar(&natsURL, "nats-url", "", "NATS server url")
	configureCmd.PersistentFlags().StringVar(&natsClusterID, "nats-cluster", "", "NATS cluster id")
	configureCmd.PersistentFlags().StringVar(&keyFile, "keyfile", "", "URL or local path to keyfile (at least 64 bytes is required)")
	configureCmd.PersistentFlags().StringVar(&knownFingerprintHex, "fingerprint", "", "If you know the fingerprint you want to use (SHAKE256 hex format), you can set it directly instead of using --keyfile")
	configureCmd.PersistentFlags().BoolVar(&useShortName, "short-names", false, "Use short channel names (channel conflicts could be more likely for a given keyfile)")
	configureCmd.PersistentFlags().BoolVar(&forceWrite, "overwrite", false, "Overwrite current configuration")
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Convey",
	Run:   ConfigureCommandFunc,
}

// ConfigureCommandFunc is a handler for the configure command
func ConfigureCommandFunc(cmd *cobra.Command, args []string) {
	var fingerprint string
	if knownFingerprintHex != "" {
		if keyFile != "" {
			errorExit("Specify either --fingerprint OR --keyfile, not both.")
		}
		if _, err := hex.DecodeString(knownFingerprintHex); err != nil || len(knownFingerprintHex) != fingerprintByteLength*2 {
			msg := fmt.Sprintf("The specified fingerprint is not %d bytes long and a valid hexidecimal string", fingerprintByteLength)
			errorExit(msg)
		}
		fingerprint = knownFingerprintHex
	} else {
		fingerprint = generateFingerprint(keyFile)
	}

	// Set config passed on the arguments passed in
	viper.Set(configKeyFingerprint, fingerprint)
	viper.Set(configKeyNatsURL, natsURL)
	viper.Set(configKeyNatsClusterID, natsClusterID)
	viper.Set(configKeyUseShortName, useShortName)

	// If a config file is found, read it in.
	configFileExists := false
	if err := viper.ReadInConfig(); err == nil {
		configFileExists = true
	}

	// If config file doesn't exist and it hasn't been set in viper, set it
	if !configFileExists && viper.ConfigFileUsed() == "" {
		home, err := homedir.Dir()
		if err != nil {
			errorExit(err.Error())
		}
		viper.SetConfigFile(path.Join(home, ".convey.yaml"))
	}

	configFilePath := viper.ConfigFileUsed()

	if forceWrite || !configFileExists {
		err := viper.WriteConfigAs(configFilePath)
		if err != nil {
			errorExit(err.Error())
		}
	} else {
		msg := fmt.Sprintf("Config file exists. Use --overwrite to overwrite the config file at %s", configFilePath)
		errorExit(msg)
	}
}
