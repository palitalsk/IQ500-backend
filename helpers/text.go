package helpers

import (
	"regexp"
	"strings"
)

// ทำความสะอาดข้อความก่อนเอาไป embed
func CleanText(text string) string {
	// Remove unknown characters
	re := regexp.MustCompile(`[^\x00-\x7Fก-๙0-9\s\p{P}]`)
	text = re.ReplaceAllString(text, "")

	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	return text
}

// แบ่งข้อความออกเป็น chunk เพื่อเอาไป embed
func SplitTextIntoChunks(text string, chunkWords int) []string {
	words := strings.Fields(text)
	var chunks []string
	for i := 0; i < len(words); i += chunkWords {
		end := i + chunkWords
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))
	}
	return chunks
}
