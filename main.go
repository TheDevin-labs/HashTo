package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
	"time"
)

const version = "1.0.0"

func newHasher(algo string) (hash.Hash, error) {
	switch strings.ToLower(algo) {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	}
	return nil, fmt.Errorf("unsupported algorithm: %s", algo)
}

func detectLocale() string {
	if v := os.Getenv("LC_ALL"); v != "" {
		return v
	}
	if v := os.Getenv("LANG"); v != "" {
		return v
	}
	return "C"
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "hashto: "+err.Error())
	os.Exit(1)
}

func defaultOutput(input, ext string) string {
	if input == "" {
		return "stdin" + ext
	}
	return input + ext
}

func swapExt(input, newExt string) string {
	for _, old := range []string{".hash", ".json", ".yaml", ".yml"} {
		if strings.HasSuffix(input, old) {
			return strings.TrimSuffix(input, old) + newExt
		}
	}
	return input + newExt
}

func runHash(input, output, algo string) error {
	var reader io.Reader
	var sourceName string

	if input != "" {
		f, err := os.Open(input)
		if err != nil {
			return err
		}
		defer f.Close()
		reader = f
		sourceName = input
	} else {
		reader = os.Stdin
		sourceName = ""
	}

	h, err := newHasher(algo)
	if err != nil {
		return err
	}

	written, err := io.Copy(h, reader)
	if err != nil {
		return err
	}

	hf := HashFile{
		Algorithm:  strings.ToLower(algo),
		Digest:     h.Sum(nil),
		Timestamp:  time.Now().Unix(),
		Locale:     detectLocale(),
		Source:     sourceName,
		SourceSize: written,
	}

	if output == "" {
		output = defaultOutput(sourceName, ".hash")
	}

	if err := WriteHashFile(output, hf); err != nil {
		return err
	}

	fmt.Printf("written %s (%s, %x)\n", output, hf.Algorithm, hf.Digest)
	return nil
}

func runToJSON(input, output string) error {
	if input == "" {
		return fmt.Errorf("--input is required for --to-json")
	}
	hf, err := ReadHashFile(input)
	if err != nil {
		return err
	}
	if output == "" {
		output = swapExt(input, ".json")
	}
	return WriteJSONRecord(output, ToRecord(hf))
}

func runToYAML(input, output string) error {
	if input == "" {
		return fmt.Errorf("--input is required for --to-yaml")
	}
	hf, err := ReadHashFile(input)
	if err != nil {
		return err
	}
	if output == "" {
		output = swapExt(input, ".yaml")
	}
	return WriteYAMLRecord(output, ToRecord(hf))
}

func runFromJSON(input, output string) error {
	if input == "" {
		return fmt.Errorf("--input is required for --from-json")
	}
	r, err := ReadJSONRecord(input)
	if err != nil {
		return err
	}
	hf, err := FromRecord(r)
	if err != nil {
		return err
	}
	if output == "" {
		output = swapExt(input, ".hash")
	}
	return WriteHashFile(output, hf)
}

func runFromYAML(input, output string) error {
	if input == "" {
		return fmt.Errorf("--input is required for --from-yaml")
	}
	r, err := ReadYAMLRecord(input)
	if err != nil {
		return err
	}
	hf, err := FromRecord(r)
	if err != nil {
		return err
	}
	if output == "" {
		output = swapExt(input, ".hash")
	}
	return WriteHashFile(output, hf)
}

func usage() {
	fmt.Fprintf(os.Stderr, `hashto - convert anything into a hash, and a hash into anything

Usage:
  hashto --hash [--input FILE] [--output FILE] [--algo md5|sha1|sha256|sha512]
  hashto --to-json   --input FILE.hash [--output FILE.json]
  hashto --to-yaml   --input FILE.hash [--output FILE.yaml]
  hashto --from-json --input FILE.json [--output FILE.hash]
  hashto --from-yaml --input FILE.yaml [--output FILE.hash]
  hashto --table-hash FILE.hash

Flags:
  -i, --input      input file (reads stdin when omitted, --hash only)
  -o, --output     output file ("-" for stdout on json/yaml conversions)
  -a, --algo       md5, sha1, sha256, sha512 (default sha256)
  -h, --help       show this help
  -v, --version    show version

Examples:
  hashto --hash --input photo.png --algo sha256
  cat notes.txt | hashto --hash --output notes.hash
  hashto --to-yaml --input notes.hash
  hashto --from-yaml --input notes.yaml --output rebuilt.hash
  hashto --table-hash notes.hash
`)
}

func main() {
	var input, output, algo, tableHash string
	var doHash, doToJSON, doToYAML, doFromJSON, doFromYAML, showHelp, showVersion bool

	flag.StringVar(&input, "input", "", "input file")
	flag.StringVar(&input, "i", "", "input file")
	flag.StringVar(&output, "output", "", "output file")
	flag.StringVar(&output, "o", "", "output file")
	flag.StringVar(&algo, "algo", "sha256", "hash algorithm")
	flag.StringVar(&algo, "a", "sha256", "hash algorithm")
	flag.StringVar(&tableHash, "table-hash", "", "show information about a .hash file")

	flag.BoolVar(&doHash, "hash", false, "compute a hash")
	flag.BoolVar(&doToJSON, "to-json", false, "convert a .hash file to JSON")
	flag.BoolVar(&doToYAML, "to-yaml", false, "convert a .hash file to YAML")
	flag.BoolVar(&doFromJSON, "from-json", false, "convert a JSON file to .hash")
	flag.BoolVar(&doFromYAML, "from-yaml", false, "convert a YAML file to .hash")

	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")

	flag.Usage = usage
	flag.Parse()

	if showHelp {
		usage()
		return
	}
	if showVersion {
		fmt.Println("hashto version " + version)
		return
	}

	if tableHash != "" {
		if err := PrintTable(tableHash); err != nil {
			fail(err)
		}
		return
	}

	if input == "" && flag.NArg() > 0 {
		input = flag.Arg(0)
	}

	switch {
	case doHash:
		if err := runHash(input, output, algo); err != nil {
			fail(err)
		}
	case doToJSON:
		if err := runToJSON(input, output); err != nil {
			fail(err)
		}
	case doToYAML:
		if err := runToYAML(input, output); err != nil {
			fail(err)
		}
	case doFromJSON:
		if err := runFromJSON(input, output); err != nil {
			fail(err)
		}
	case doFromYAML:
		if err := runFromYAML(input, output); err != nil {
			fail(err)
		}
	default:
		usage()
		os.Exit(1)
	}
}
