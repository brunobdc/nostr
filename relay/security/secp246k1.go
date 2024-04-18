package security

// #include <secp256k1/include/secp256k1.h>
// #include <secp256k1/include/secp256k1_schnorrsig.h>
// #cgo LDFLAGS: ${SRCDIR}/secp256k1/.libs/libsecp256k1.a
import "C"
import (
	"crypto/rand"
	"errors"
	"unsafe"
)

func RandomKeyPair() ([32]byte, [32]byte) {
	ctx, err := C.secp256k1_context_create(C.SECP256K1_CONTEXT_NONE)
	if err != nil {
		panic(err)
	}
	defer C.secp256k1_context_destroy(ctx)

	var random [32]byte
	lenRandom, err := rand.Read(random[:])
	if err != nil {
		panic(err)
	}
	if lenRandom != 32 {
		panic(errors.New("failed to generate a valid random number"))
	}
	valid, err := C.secp256k1_context_randomize(ctx, (*C.uchar)(&random[0]))
	if err != nil {
		panic(err)
	}
	if valid == 0 {
		panic(errors.New("unexpected error in context randomize"))
	}

	var privateKey [32]byte
	lenPrivateKey, err := rand.Read(privateKey[:])
	if err != nil {
		panic(err)
	}
	if lenPrivateKey != 32 {
		panic(errors.New("failed to generate a valid random number"))
	}

	var keyPair *C.secp256k1_keypair
	valid, err = C.secp256k1_keypair_create(ctx, keyPair, (*C.uchar)(&privateKey[0]))
	if err != nil {
		panic(err)
	}
	if valid == 0 {
		panic(errors.New("unexpected error in keypar create"))
	}

	var cPublicKey *C.secp256k1_xonly_pubkey
	valid, err = C.secp256k1_keypair_xonly_pub(ctx, cPublicKey, nil, keyPair)
	if err != nil {
		panic(err)
	}
	if valid == 0 {
		panic(errors.New("unexpected error in getting public key from key pair"))
	}

	var cSerializedPublicKey *C.uchar
	valid, err = C.secp256k1_xonly_pubkey_serialize(ctx, cSerializedPublicKey, cPublicKey)
	if err != nil {
		panic(err)
	}
	if valid == 0 {
		panic(errors.New("unexpected error in getting serialized public key"))
	}

	publicKey := C.GoBytes(unsafe.Pointer(cSerializedPublicKey), C.int(unsafe.Sizeof(cSerializedPublicKey)))

	return privateKey, [32]byte(publicKey)
}
