package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func isArmorProtected(f *os.File) bool {
	_, err := armor.Decode(f)
	f.Seek(0, io.SeekStart)
	return err == nil
}

func readPublicKeyRing(publicKeyRing string) (openpgp.EntityList, error) {
	keyRingReader, err := os.Open(publicKeyRing)
	if err != nil {
		return nil, fmt.Errorf("Failed to open public key '%s': %s", publicKeyRing, err)
	}
	var keyring openpgp.EntityList
	if isArmorProtected(keyRingReader) {
		keyring, err = openpgp.ReadArmoredKeyRing(keyRingReader)
	} else {
		keyring, err = openpgp.ReadKeyRing(keyRingReader)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to parse public key '%s': %s", publicKeyRing, err)
	}
	return keyring, nil
}

func verifySignature(file *os.File, signature *os.File, keyRing openpgp.EntityList) error {
	var err error
	if isArmorProtected(signature) {
		_, err = openpgp.CheckArmoredDetachedSignature(keyRing, file, signature)
	} else {
		_, err = openpgp.CheckDetachedSignature(keyRing, file, signature)
	}

	if err != nil {
		return fmt.Errorf("Verification of file failed: %s", err)
	}
	return nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "USAGE: %s file signature publickeyring\n", os.Args[0])
		os.Exit(1)
	}
	path := os.Args[1]
	signaturePath := os.Args[2]
	publicKeyRing := os.Args[3]

	keyRing, err := readPublicKeyRing(publicKeyRing)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file '%s': %s\n",
			path, err)
		os.Exit(1)
	}

	signature, err := os.Open(signaturePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open detached signature '%s': %s\n",
			signaturePath, err)
		os.Exit(1)
	}

	err = verifySignature(file, signature, keyRing)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// everything is fine
}
