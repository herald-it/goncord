package pwd_hash_test

import (
	"crypto/md5"
	"crypto/sha1"
	"github.com/herald-it/goncord/pwd_hash"
	"testing"
)

const pwd = "my_super_crypto_password"

func TestCorrectHashFunc(t *testing.T) {
	sha1_sum := sha1.Sum([]byte(pwd))
	vanilla_sum := md5.Sum(sha1_sum[:])
	vanilla_sum_slice := vanilla_sum[:]

	my_pwd_hash := pwd_hash.Sum([]byte(pwd))

	if string(vanilla_sum_slice) != string(my_pwd_hash) {
		t.Errorf("Original hash sum not equal. \norigin:  %v\ncurrent: %v",
			string(vanilla_sum_slice), string(my_pwd_hash))
	}

	hash := pwd_hash.New()
	hash.Write([]byte(pwd))

	hash_sum := hash.Sum(nil)
	if string(hash_sum) != string(vanilla_sum_slice) {
		t.Error("Not equal")
	}
}

func TestCorrectWriteFunc(t *testing.T) {
	sha1_sum := sha1.Sum([]byte(pwd))
	vanilla_sum := md5.Sum(sha1_sum[:])
	vanilla_sum_slice := vanilla_sum[:]

	hash := pwd_hash.New()
	hash.Write([]byte(pwd))

	hash_sum := hash.Sum(nil)
	if string(hash_sum) != string(vanilla_sum_slice) {
		t.Error("Not equal")
	}
}

func TestCorrectData(t *testing.T) {
	if pwd_hash.Size != 16 {
		t.Error("Size not equal")
	}

	if pwd_hash.BlockSize != 64 {
		t.Error("Block size")
	}
}
