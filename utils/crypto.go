package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"runtime"

	"golang.org/x/crypto/argon2"
)

const (
	iteration     = 2
	memory        = 128 * 1024
	hashKeyLength = 128
	saltLength    = 32
)

type argon2idConf struct {
	time    uint32 // time represents the number of passed over the specified memory.
	memory  uint32 // cpu memory to be used.
	threads uint8  // threads for parallelism aspect of the algorithm.
	keyLen  uint32 // keyLen of the generate hash key.
	saltLen uint32 // saltLen the length of the salt used.
}

// NewArgon2idHash constructor function for
// Argon2idHash.
func newArgon2idHash(time, saltLen uint32, memory uint32, threads uint8, keyLen uint32) *argon2idConf {
	return &argon2idConf{
		time:    time,
		saltLen: saltLen,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
}

func RandomByte(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// GenerateHash using the password and provided salt.
// If not salt value provided fallback to random value
// generated of a given length.
func (a *argon2idConf) GenerateHash(password, salt []byte) ([]byte, error) {
	var err error
	// If salt is not provided generate a salt of the configured salt length.
	if len(salt) == 0 {
		salt, err = RandomByte(a.saltLen)
	}
	if err != nil {
		return nil, err
	}
	// Generate hash
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)
	// Return the generated hash and salt used for storage.
	return append(salt, hash...), nil
}

// Compare generated hash with store hash.
func (a *argon2idConf) Compare(passwordHash, password []byte) (bool, error) {
	salt := passwordHash[:a.saltLen]
	// Generate hash for comparison.
	hash, err := a.GenerateHash(password, salt)
	if err != nil {
		return false, err
	}
	// Compare the generated hash with the stored hash.
	return bytes.Equal(passwordHash, hash[a.saltLen:]), nil
}

func Hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return base64.URLEncoding.EncodeToString(hash[:])
}

var Argon2id = newArgon2idHash(iteration, saltLength, memory, uint8(runtime.NumCPU()), hashKeyLength)
