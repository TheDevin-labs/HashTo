package main

import "fmt"

func PrintTable(path string) error {
	hf, err := ReadHashFile(path)
	if err != nil {
		return err
	}
	r := ToRecord(hf)

	rows := [][2]string{
		{"File", path},
		{"Source", displayOrDash(r.Source)},
		{"Source Size", fmt.Sprintf("%d bytes", r.SourceSize)},
		{"Algorithm", r.Algorithm},
		{"Digest (hex)", r.DigestHex},
		{"Digest (base64)", r.DigestBase64},
		{"Created", r.Created},
		{"Locale/Charset", displayOrDash(r.Locale)},
	}

	labelWidth := 0
	for _, row := range rows {
		if len(row[0]) > labelWidth {
			labelWidth = len(row[0])
		}
	}

	separator := repeatStr("-", labelWidth+3)
	fmt.Println(separator)
	for _, row := range rows {
		fmt.Printf("%-*s : %s\n", labelWidth, row[0], row[1])
	}
	fmt.Println(separator)
	return nil
}

func displayOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func repeatStr(s string, n int) string {
	out := make([]byte, 0, n*len(s))
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}
