package bitcoin

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParsePublicKey(t *testing.T) {
	xOnlyPub := "87176beec39cbbd2f1999209894684e1620bc39ebc2a704add0edb23d0207d7e"

	pubKey, err := ParsePublicKey(xOnlyPub)
	require.NoError(t, err)

	require.Equal(t, xOnlyPub, pubKey.X().Text(16))

	pubKey2, err := ParsePublicKey(hex.EncodeToString(pubKey.SerializeCompressed()))
	require.NoError(t, err)

	require.Equal(t, xOnlyPub, pubKey2.X().Text(16))            // should still have same X value
	require.Equal(t, pubKey.Y().Text(16), pubKey2.Y().Text(16)) // should have same Y value

	pubKey3, err := ParsePublicKey(hex.EncodeToString(pubKey.SerializeUncompressed()))
	require.NoError(t, err)

	require.Equal(t, xOnlyPub, pubKey3.X().Text(16))            // should still have same X value
	require.Equal(t, pubKey.Y().Text(16), pubKey3.Y().Text(16)) // should have same Y value
}
