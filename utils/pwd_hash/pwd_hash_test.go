package pwd_hash_test

import (
	"crypto/md5"
	"crypto/sha1"
	"testing"

	pwdHash "github.com/herald-it/goncord/utils/pwd_hash"
	. "github.com/smartystreets/goconvey/convey"
)

const pwd = "my_super_crypto_password"

func TestCorrectHashFunc(t *testing.T) {
	Convey("Test correct hash func", t, func() {
		sha1Sum := sha1.Sum([]byte(pwd))
		vanillaSum := md5.Sum(sha1Sum[:])
		vanillaSumSlice := vanillaSum[:]

		myPwdHash := pwdHash.Sum([]byte(pwd))

		Convey("Equal self hash and original hash", func() {
			So(string(vanillaSumSlice), ShouldEqual, string(myPwdHash))
		})

		hash := pwdHash.New()
		hash.Write([]byte(pwd))

		Convey("Test hash sum method", func() {
			hash_sum := hash.Sum(nil)
			So(string(hash_sum), ShouldEqual, string(vanillaSumSlice))
		})
	})
}

func TestCorrectWriteFunc(t *testing.T) {
	Convey("test correct write function", t, func() {
		sha1Sum := sha1.Sum([]byte(pwd))
		vanillaSum := md5.Sum(sha1Sum[:])
		vanillaSumSlice := vanillaSum[:]

		hash := pwdHash.New()
		hash.Write([]byte(pwd))

		hashSum := hash.Sum(nil)

		So(string(hashSum), ShouldEqual, string(vanillaSumSlice))
	})
}

func TestCorrectData(t *testing.T) {
	Convey("Test constant data", t, func() {
		So(pwdHash.Size, ShouldEqual, 16)
		So(pwdHash.BlockSize, ShouldEqual, 64)

		tmp := pwdHash.New()
		So(tmp.Size(), ShouldEqual, 16)
		So(tmp.BlockSize(), ShouldEqual, 64)
	})
}
