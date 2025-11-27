package helpers

import (
	"strings"

	pdf "github.com/ledongthuc/pdf"
)

// อ่าน text จากไฟล์ PDF
func ReadPdfText(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var out strings.Builder
	totalPage := r.NumPage()
	// loop ทุกหน้า
	for i := 1; i <= totalPage; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		out.WriteString(text + "\n")
	}
	return out.String(), nil
}
