package bitcoin

import "github.com/cockroachdb/errors"

func ParseWitnessSignatureBIP322(sigBytes []byte) ([][]byte, error) {
	if len(sigBytes) < 66 {
		return nil, errors.Errorf("signature length (%d) is invalid, must be at least 66 bytes", len(sigBytes))
	}

	// Extract the number of signatures and the size of each signature
	signatureCount := int(sigBytes[0])
	signatureSize := int(sigBytes[1])

	witnesses := make([][]byte, signatureCount)

	if len(sigBytes) != 2+signatureCount*signatureSize {
		return nil, errors.Errorf("signature length (%d) is invalid, expected %d bytes", len(sigBytes), 2+signatureCount*signatureSize)
	}

	// Extract all the witness signatures encoded in the signature
	for i := 0; i < signatureCount; i++ {
		signature := make([]byte, signatureSize)
		copy(signature, sigBytes[2+i*signatureSize:2+(i+1)*signatureSize])
		witnesses[i] = signature
	}
	return witnesses, nil
}
