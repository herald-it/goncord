package pwd_hash

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
)

const Size = 16
const BlockSize = 64

type digest struct {
	bytes []byte
}

func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Reset() {
	d.bytes = []byte{}
}

func (d *digest) Size() int {
	return Size
}

func (d *digest) BlockSize() int {
	return BlockSize
}

func (d *digest) Sum(b []byte) []byte {
	return Sum(d.bytes)
}

// Implement writer interface
func (d *digest) Write(p []byte) (n int, err error) {
	d.bytes = append(d.bytes, p...)

	n = len(p)
	err = nil
	return
}

func Sum(in []byte) []byte {
	bytes := sha1.Sum(in)
	md5s := md5.Sum(bytes[:])

	return md5s[:]
}

func HashPassword(password string) string {
	return hex.EncodeToString(Sum([]byte(password)))
}
