package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func WriteYAMLRecord(path string, r HashRecord) error {
	var b strings.Builder
	b.WriteString("algorithm: " + yamlScalar(r.Algorithm) + "\n")
	b.WriteString("digest_hex: " + yamlScalar(r.DigestHex) + "\n")
	b.WriteString("digest_base64: " + yamlScalar(r.DigestBase64) + "\n")
	b.WriteString("timestamp: " + strconv.FormatInt(r.Timestamp, 10) + "\n")
	b.WriteString("created: " + yamlScalar(r.Created) + "\n")
	b.WriteString("locale: " + yamlScalar(r.Locale) + "\n")
	b.WriteString("source: " + yamlScalar(r.Source) + "\n")
	b.WriteString("source_size: " + strconv.FormatInt(r.SourceSize, 10) + "\n")

	data := []byte(b.String())
	if path == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func yamlScalar(s string) string {
	if s == "" {
		return "\"\""
	}
	needsQuote := false
	for _, r := range s {
		if r == ':' || r == '#' || r == '\'' || r == '"' || r == '\n' {
			needsQuote = true
			break
		}
	}
	if !needsQuote {
		return s
	}
	escaped := strings.ReplaceAll(s, "\"", "\\\"")
	return "\"" + escaped + "\""
}

func unquoteYAML(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		inner := s[1 : len(s)-1]
		return strings.ReplaceAll(inner, "\\\"", "\"")
	}
	return s
}

func ReadYAMLRecord(path string) (HashRecord, error) {
	var r HashRecord
	data, err := os.ReadFile(path)
	if err != nil {
		return r, err
	}

	fields := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := unquoteYAML(strings.TrimSpace(line[idx+1:]))
		fields[key] = value
	}

	r.Algorithm = fields["algorithm"]
	r.DigestHex = fields["digest_hex"]
	r.DigestBase64 = fields["digest_base64"]
	r.Created = fields["created"]
	r.Locale = fields["locale"]
	r.Source = fields["source"]

	if ts, ok := fields["timestamp"]; ok && ts != "" {
		v, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return r, errors.New("invalid timestamp: " + err.Error())
		}
		r.Timestamp = v
	}
	if ss, ok := fields["source_size"]; ok && ss != "" {
		v, err := strconv.ParseInt(ss, 10, 64)
		if err != nil {
			return r, errors.New("invalid source_size: " + err.Error())
		}
		r.SourceSize = v
	}

	return r, nil
}
