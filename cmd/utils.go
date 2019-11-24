package cmd

import (
	"encoding/hex"
	"fmt"
)

// FingerprintByteLength is our fingerprint byte length
const FingerprintByteLength = 64

// InvalidFingerprintMsg is the user-friendly error message if fingerprint is invalid
var InvalidFingerprintMsg = fmt.Sprintf("The fingerprint in use is not %d bytes long and a valid hexidecimal string", FingerprintByteLength)

// IsValidFingerprint determines if a fingerprint is valid.
func IsValidFingerprint(fingerprint string) bool {
	if _, err := hex.DecodeString(fingerprint); err != nil || len(fingerprint) != FingerprintByteLength*2 {
		return false
	}
	return true
}
