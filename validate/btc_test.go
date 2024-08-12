package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetValidBtcAddr(t *testing.T) {

	validAddress := []string{
		"bc1qvkh89cjz9jly7n9d0mszku720jtg9lr8c4eyz9",
		"bc1qm6mdp6yl2t7f0mge5ef83q7yjw5zr5heh7k3c7",
		"1DnAcVuacU8dcARDvyLD9yxAk2ZHCvcfnD",
		"3B8CnrpoiZxdPxZ4CxghXv45N7XET6T5uQ",
		"31kMmzFEM6pVBbcQHHsypd9T1CxHHoZtqH",
		"1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2",
		"3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
		"bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq",
		"bc1pcejefkj9658nyslnr0qup6dzzt5a0nawd535uymrqlhaqvztvves84wdaf",
		"bc1qeklep85ntjz4605drds6aww9u0qr46qzrv5xswd35uhjuj8ahfcqgf6hak",
		"bc1pxwww0ct9ue7e8tdnlmug5m2tamfn7q06sahstg39ys4c9f3340qqxrdu9k",
	}

	for _, addr := range validAddress {
		_, err := GetValidBtcAddr(addr)
		require.NoErrorf(t, err, "GetValidBtcAddr(%s) returned an error", addr)
	}

}

func Test_CheckValidSecp256k1PubKey(t *testing.T) {
	err := CheckValidSecp256k1PubKey("87176beec39cbbd2f1999209894684e1620bc39ebc2a704add0edb23d0207d7e")
	require.NoError(t, err)

	err = CheckValidSecp256k1PubKey("87176ceec39cbbd2f1999209894684e1620bc39ebc2a704add0edb23d0207d7e")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not on the secp256k1 curve")
}
