package testutils

import (
	"encoding/base64"
	"strings"
)

// GenRandGalxeID generates a random Galxe ID
func GenRandGalxeID() string {

	idBytes := make([]byte, 16)
	_, err := random.Read(idBytes)
	if err != nil {
		panic(err)
	}
	// bit of a hack, but this is just for testing
	out := base64.RawStdEncoding.EncodeToString(idBytes)

	out = strings.ReplaceAll(out, "+", "a")
	out = strings.ReplaceAll(out, "/", "b")
	return out[:21]

}
