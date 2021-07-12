package storageInterface

import (
	"errors"
	"golang.org/x/crypto/sha3"
)

const (
	domainNone  = 1
	domainSHA3  = 0x06
	domainSHAKE = 0x1f
)

const (
	// WORDWIDTH is in [0,6]
	MATRIXROWS = 5
	MATRIXCOLS = 5
	WORDWIDTH  = 6
	WORDBITS   = 1 << WORDWIDTH
	ROUND      = 12 + 2*WORDWIDTH
	OUTLENGTH  = 256
	CAPACITY   = OUTLENGTH << 1
	RATE       = MATRIXROWS*MATRIXCOLS*WORDBITS - CAPACITY
)

func Sha3Encode(encodeStr string, secretKey string) (string, error) {
	// A MAC with 32 bytes of output has 256-bit security strength -- if you use at least a 32-byte-long key.
	out := make([]byte, 0)
	hash := sha3.New384()
	// Write the key into the hash.
	writeSK, writeSKErr := hash.Write([]byte(secretKey))
	if writeSKErr != nil {
		return "", writeSKErr
	}
	// Now write the data.
	writeSign, writeSignErr := hash.Write([]byte(encodeStr))
	if writeSignErr != nil {
		return "", writeSignErr
	}
	// Read 32 bytes of output from the hash into h.
	out = hash.Sum(out)
	if writeSK != len(secretKey) || writeSign != len(encodeStr) || writeSK != SKLENGTH || len(out) != DIGESTLENGTH {
		return "", errors.New("wrong length with the sk")
	}
	return string(out), nil
}
