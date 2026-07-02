package main

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"
)

type HashRecord struct {
	Algorithm    string `json:"algorithm"`
	DigestHex    string `json:"digest_hex"`
	DigestBase64 string `json:"digest_base64"`
	Timestamp    int64  `json:"timestamp"`
	Created      string `json:"created"`
	Locale       string `json:"locale"`
	Source       string `json:"source"`
	SourceSize   int64  `json:"source_size"`
}

func ToRecord(hf HashFile) HashRecord {
	return HashRecord{
		Algorithm:    hf.Algorithm,
		DigestHex:    hex.EncodeToString(hf.Digest),
		DigestBase64: base64.StdEncoding.EncodeToString(hf.Digest),
		Timestamp:    hf.Timestamp,
		Created:      time.Unix(hf.Timestamp, 0).UTC().Format(time.RFC3339),
		Locale:       hf.Locale,
		Source:       hf.Source,
		SourceSize:   hf.SourceSize,
	}
}

func FromRecord(r HashRecord) (HashFile, error) {
	digest, err := hex.DecodeString(r.DigestHex)
	if err != nil {
		return HashFile{}, errors.New("invalid digest_hex: " + err.Error())
	}
	return HashFile{
		Algorithm:  r.Algorithm,
		Digest:     digest,
		Timestamp:  r.Timestamp,
		Locale:     r.Locale,
		Source:     r.Source,
		SourceSize: r.SourceSize,
	}, nil
}
