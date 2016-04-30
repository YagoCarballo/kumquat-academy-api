package tools

import (
	"testing"

	. "github.com/franela/goblin"
	"crypto/rsa"
	"os"
)

const TMP_PRIVATE_KEY_PATH = "testKey.pem"
const TMP_PUBLIC_KEY_PATH = "testKey.pub"

func Test_Tools(t *testing.T) {
	g := Goblin(t)
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey

	g.Describe("When loading the Private Key File, ", func() {
		g.It("static file should not exist", func () {
			exists := ShouldServeStatic("-unexistent-file-")
			g.Assert(exists).IsFalse()
		})

		g.It("loading an Invalid Path does not return an error and creates the file.", func() {
			// Variable is empty
			g.Assert(privateKey == nil).IsTrue()

			var err error;
			privateKey, publicKey, err = LoadKey(TMP_PRIVATE_KEY_PATH, TMP_PUBLIC_KEY_PATH)
			g.Assert(err == nil).IsTrue()
			g.Assert(privateKey != nil).IsTrue()
			g.Assert(publicKey != nil).IsTrue()
		})

		g.It("static file should exist", func () {
			exists := ShouldServeStatic(TMP_PRIVATE_KEY_PATH)
			g.Assert(exists).IsTrue()
		})

		g.It("loading a Valid Path reads the file.", func() {
			existingPrivateKey, existingPublicKey, err := LoadKey(TMP_PRIVATE_KEY_PATH, TMP_PUBLIC_KEY_PATH)
			g.Assert(err == nil).IsTrue()
			g.Assert(existingPrivateKey != nil).IsTrue()
			g.Assert(existingPublicKey != nil).IsTrue()

			// Public Keys are the same
			g.Assert(privateKey != nil).IsTrue()
			g.Assert(privateKey.E == existingPrivateKey.E).IsTrue()
			g.Assert(publicKey.E == existingPublicKey.E).IsTrue()
		})
	})

	// Remove the Temporary Key
	os.Remove(TMP_PRIVATE_KEY_PATH)
	os.Remove(TMP_PUBLIC_KEY_PATH)
}
