package compiler_test

import (
	"fmt"
	"testing"

	"github.com/nspcc-dev/neo-go/pkg/core/interop/interopnames"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// In this test we only check that needed interop
// is called with the provided arguments in the right order.
func TestVerifyGood(t *testing.T) {
	msg := []byte("test message")
	pub, sig := signMessage(t, msg)
	src := getVerifyProg(pub, sig, msg)

	v, p := vmAndCompileInterop(t, src)
	p.interops[interopnames.ToID([]byte(interopnames.NeoCryptoVerifyWithECDsaSecp256r1))] = func(v *vm.VM) error {
		assert.Equal(t, msg, v.Estack().Pop().Bytes())
		assert.Equal(t, pub, v.Estack().Pop().Bytes())
		assert.Equal(t, sig, v.Estack().Pop().Bytes())
		v.Estack().PushVal(true)
		return nil
	}

	require.NoError(t, v.Run())
}

func signMessage(t *testing.T, msg []byte) ([]byte, []byte) {
	key, err := keys.NewPrivateKey()
	require.NoError(t, err)

	sig := key.Sign(msg)
	pub := key.PublicKey().Bytes()

	return pub, sig
}

func getVerifyProg(pub, sig, msg []byte) string {
	pubS := fmt.Sprintf("%#v", pub)
	sigS := fmt.Sprintf("%#v", sig)
	msgS := fmt.Sprintf("%#v", msg)

	return `
		package hello

		import "github.com/nspcc-dev/neo-go/pkg/interop/crypto"

		func Main() bool {
			pub := ` + pubS + `
			sig := ` + sigS + `
			msg := ` + msgS + `
			return crypto.ECDsaSecp256r1Verify(msg, pub, sig)
		}
	`
}
