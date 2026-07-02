package main

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

const magic = "HTOH"
const formatVersion = 1

const (
	algoMD5 = 1
	algoSHA1 = 2
	algoSHA256 = 3
	algoSHA512 = 4
)

type HashFile struct {
	Algorithm  string
	Digest     []byte
	Timestamp  int64
	Locale     string
	Source     string
	SourceSize int64
}

func algoToID(name string) (byte, error) {
	switch name {
	case "md5":
		return algoMD5, nil
	case "sha1":
		return algoSHA1, nil
	case "sha256":
		return algoSHA256, nil
	case "sha512":
		return algoSHA512, nil
	}
	return 0, errors.New("unsupported algorithm: " + name)
}

func idToAlgo(id byte) (string, error) {
	switch id {
	case algoMD5:
		return "md5", nil
	case algoSHA1:
		return "sha1", nil
	case algoSHA256:
		return "sha256", nil
	case algoSHA512:
		return "sha512", nil
	}
	return "", errors.New("unknown algorithm id in hash file")
}

func WriteHashFile(path string, hf HashFile) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	algoID, err := algoToID(hf.Algorithm)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(magic); err != nil {
		return err
	}
	if err := binary.Write(f, binary.BigEndian, uint8(formatVersion)); err != nil {
		return err
	}
	if err := binary.Write(f, binary.BigEndian, algoID); err != nil {
		return err
	}
	if err := binary.Write(f, binary.BigEndian, uint16(len(hf.Digest))); err != nil {
		return err
	}
	if _, err := f.Write(hf.Digest); err != nil {
		return err
	}
	if err := binary.Write(f, binary.BigEndian, hf.Timestamp); err != nil {
		return err
	}
	localeBytes := []byte(hf.Locale)
	if err := binary.Write(f, binary.BigEndian, uint16(len(localeBytes))); err != nil {
		return err
	}
	if _, err := f.Write(localeBytes); err != nil {
		return err
	}
	sourceBytes := []byte(hf.Source)
	if err := binary.Write(f, binary.BigEndian, uint16(len(sourceBytes))); err != nil {
		return err
	}
	if _, err := f.Write(sourceBytes); err != nil {
		return err
	}
	if err := binary.Write(f, binary.BigEndian, hf.SourceSize); err != nil {
		return err
	}
	return nil
}

func ReadHashFile(path string) (HashFile, error) {
	var hf HashFile
	f, err := os.Open(path)
	if err != nil {
		return hf, err
	}
	defer f.Close()

	magicBuf := make([]byte, 4)
	if _, err := io.ReadFull(f, magicBuf); err != nil {
		return hf, err
	}
	if string(magicBuf) != magic {
		return hf, errors.New("not a valid hashto file: " + path)
	}

	var ver uint8
	if err := binary.Read(f, binary.BigEndian, &ver); err != nil {
		return hf, err
	}

	var algoID uint8
	if err := binary.Read(f, binary.BigEndian, &algoID); err != nil {
		return hf, err
	}
	algoName, err := idToAlgo(algoID)
	if err != nil {
		return hf, err
	}
	hf.Algorithm = algoName

	var digestLen uint16
	if err := binary.Read(f, binary.BigEndian, &digestLen); err != nil {
		return hf, err
	}
	digest := make([]byte, digestLen)
	if _, err := io.ReadFull(f, digest); err != nil {
		return hf, err
	}
	hf.Digest = digest

	if err := binary.Read(f, binary.BigEndian, &hf.Timestamp); err != nil {
		return hf, err
	}

	var localeLen uint16
	if err := binary.Read(f, binary.BigEndian, &localeLen); err != nil {
		return hf, err
	}
	localeBytes := make([]byte, localeLen)
	if _, err := io.ReadFull(f, localeBytes); err != nil {
		return hf, err
	}
	hf.Locale = string(localeBytes)

	var sourceLen uint16
	if err := binary.Read(f, binary.BigEndian, &sourceLen); err != nil {
		return hf, err
	}
	sourceBytes := make([]byte, sourceLen)
	if _, err := io.ReadFull(f, sourceBytes); err != nil {
		return hf, err
	}
	hf.Source = string(sourceBytes)

	if err := binary.Read(f, binary.BigEndian, &hf.SourceSize); err != nil {
		return hf, err
	}

	return hf, nil
}
