package pwd_hash_test

import (
	"crypto/md5"
	"crypto/sha1"
	"github.com/herald-it/goncord/pwd_hash"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const pwd = "my_super_crypto_password"

func TestCorrectHashFunc(t *testing.T) {
	Convey("Test correct hash func", t, func() {
		sha1_sum := sha1.Sum([]byte(pwd))
		vanilla_sum := md5.Sum(sha1_sum[:])
		vanilla_sum_slice := vanilla_sum[:]

		my_pwd_hash := pwd_hash.Sum([]byte(pwd))

		Convey("Equal self hash and original hash", func() {
			So(string(vanilla_sum_slice), ShouldEqual, string(my_pwd_hash))
		})

		hash := pwd_hash.New()
		hash.Write([]byte(pwd))

		Convey("Test hash sum method", func() {
			hash_sum := hash.Sum(nil)
			So(string(hash_sum), ShouldEqual, string(vanilla_sum_slice))
		})
	})
}

func TestCorrectWriteFunc(t *testing.T) {
	Convey("test correct write function", t, func() {
		sha1_sum := sha1.Sum([]byte(pwd))
		vanilla_sum := md5.Sum(sha1_sum[:])
		vanilla_sum_slice := vanilla_sum[:]

		hash := pwd_hash.New()
		hash.Write([]byte(pwd))

		hash_sum := hash.Sum(nil)

		So(string(hash_sum), ShouldEqual, string(vanilla_sum_slice))
	})
}

func TestCorrectData(t *testing.T) {
	Convey("Test constant data", t, func() {
		So(pwd_hash.Size, ShouldEqual, 16)
		So(pwd_hash.BlockSize, ShouldEqual, 64)
	})
}
