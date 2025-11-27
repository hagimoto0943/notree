package main

import (
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractText(f *os.File) (string, error) {
	st, err := f.Stat()
	if err != nil {
		return "", err
	}

	r, err := pdf.NewReader(f, st.Size())
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		text, err := p.GetPlainText(nil)
		if err == nil {
			b.WriteString(text)
			b.WriteString("\n")
		}
	}
	return CleanText(b.String()), nil
}

func CleanText(text string) string {
	r := strings.NewReplacer(
		"\n\n", "PLACEHOLDER",
		"\n", "",
		"PLACEHOLDER", "\n\n",
	)
	return r.Replace(text)
}
