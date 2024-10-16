package bitcoin

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cockroachdb/errors"
	secp2561k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/samber/lo"
)

const ECDSASignatureLen = 65

// Based on https://github.com/BitonicNL/verify-signed-message/blob/main/internal/generic/verify.go
func RecoverPublicKey(message string, signatureDecoded []byte) (*secp2561k1.PublicKey, error) {
	// Ensure signature has proper length
	if len(signatureDecoded) != ECDSASignatureLen {
		return nil, errors.Errorf("invalid signature length: %d instead of %d", len(signatureDecoded), ECDSASignatureLen)
	}

	// Ensure signature has proper recovery flag
	recoveryFlag := int(signatureDecoded[0])
	if !lo.Contains[int](AllFlags(), recoveryFlag) {
		return nil, errors.Errorf("invalid recovery flag: %d", recoveryFlag)
	}

	// Reset recovery flag after obtaining keyID for Trezor
	if lo.Contains[int](TrezorFlags(), recoveryFlag) {
		signatureDecoded[0] = byte(27 + GetKeyID(recoveryFlag))
	}

	// Make and hash the message
	messageHash := chainhash.DoubleHashB([]byte(CreateMagicMessage(message)))

	// Recover the public key from signature and message hash
	publicKey, _, err := ecdsa.RecoverCompact(signatureDecoded, messageHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not recover pubkey")
	}

	return publicKey, nil
}

// ParsePublicKey parses a public key from a hex string
func ParsePublicKey(pubKeyHex string) (*secp2561k1.PublicKey, error) {
	if pubKeyHex[:2] == "0x" {
		pubKeyHex = pubKeyHex[2:]
	}
	pubKey, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode public key \"%s\"", pubKeyHex)
	}
	if len(pubKey) == 32 { // X only public key, convert to compressed public key
		pubKey = append([]byte{0x02}, pubKey...)
	}
	return secp2561k1.ParsePubKey(pubKey)
}

// ParseSignature is a convenience function to parse a signature from a hex or base64 string
func ParseSignature(signature string) ([]byte, error) {
	decodedSignature, err := hex.DecodeString(signature)
	if err == nil {
		return decodedSignature, nil
	}
	// likely base64
	return base64.StdEncoding.DecodeString(signature)
}
