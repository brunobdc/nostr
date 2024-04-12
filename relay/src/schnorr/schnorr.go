package schnorr

// #include <secp256k1/include/secp256k1.h>
// #include <secp256k1/include/secp256k1_schnorrsig.h>
// #cgo LDFLAGS: ${SRCDIR}/secp256k1/.libs/libsecp256k1.a
import "C"

func Verify(sig []byte, msg []byte, pk []byte) bool {
	var publicKey *C.secp256k1_xonly_pubkey
	_, err := C.secp256k1_xonly_pubkey_parse(C.secp256k1_context_static, publicKey, (*C.uchar)(&pk[0]))
	if err != nil {
		panic(err)
	}
	result, err := C.secp256k1_schnorrsig_verify(C.secp256k1_context_static, (*C.uchar)(&sig[0]), (*C.uchar)(&msg[0]), C.size_t(len(msg)), publicKey)
	if err != nil {
		panic(err)
	}
	if int(result) > 0 {
		return true
	}
	return false
}
