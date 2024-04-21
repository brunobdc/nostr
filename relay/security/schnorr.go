package security

// #include <secp256k1/include/secp256k1.h>
// #include <secp256k1/include/secp256k1_schnorrsig.h>
// #cgo LDFLAGS: ${SRCDIR}/secp256k1/.libs/libsecp256k1.a
import "C"
import (
	"crypto/rand"
	"errors"
)

type Signature interface {
	VerifySignature(sig [64]byte, msg [32]byte, pubKey [32]byte) (bool, error)
}

type SchnorrSignature struct{}

func (signature SchnorrSignature) VerifySignature(sig [64]byte, msg [32]byte, pubKey [32]byte) (bool, error) {
	ctx, err := C.secp256k1_context_create(C.SECP256K1_CONTEXT_NONE)
	if err != nil {
		return false, err
	}
	defer C.secp256k1_context_destroy(ctx)

	var random [32]byte
	lenRandom, err := rand.Read(random[:])
	if err != nil {
		return false, err
	}
	if lenRandom != 32 {
		return false, errors.New("failed to generate a valid random number")
	}
	valid, err := C.secp256k1_context_randomize(ctx, (*C.uchar)(&random[0]))
	if err != nil {
		return false, err
	}
	if valid == 0 {
		return false, errors.New("unexpected error in context randomize")
	}

	publicKey := new(C.secp256k1_xonly_pubkey)
	result, err := C.secp256k1_xonly_pubkey_parse(ctx, publicKey, (*C.uchar)(&pubKey[0]))
	if result == 0 {
		return false, errors.New("couldn't parse the public key")
	}
	if err != nil {
		return false, err
	}
	result, err = C.secp256k1_schnorrsig_verify(ctx, (*C.uchar)(&sig[0]), (*C.uchar)(&msg[0]), C.size_t(len(msg)), publicKey)
	if err != nil {
		return false, err
	}
	if int(result) > 0 {
		return true, nil
	}

	return false, nil
}

func SignMessage(privateKey [32]byte, msg [32]byte) [64]byte {
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

	var keyPair *C.secp256k1_keypair
	valid, err = C.secp256k1_keypair_create(ctx, keyPair, (*C.uchar)(&privateKey[0]))
	if err != nil {
		panic(err)
	}
	if valid == 0 {
		panic(errors.New("unexpected error in keypar create"))
	}

	var auxiliaryRand [32]byte
	lenAuxiliaryRand, err := rand.Read(auxiliaryRand[:])
	if err != nil {
		panic(err)
	}
	if lenAuxiliaryRand != 32 {
		panic(errors.New("failed to generate a valid random number"))
	}

	var signature [64]byte
	valid, err = C.secp256k1_schnorrsig_sign32(ctx, (*C.uchar)(&signature[0]), (*C.uchar)(&msg[0]), keyPair, (*C.uchar)(&auxiliaryRand[0]))
	if err != nil {
		panic(err)
	}
	if valid == 0 {
		panic(errors.New("unexpected while signing"))
	}

	return signature
}
