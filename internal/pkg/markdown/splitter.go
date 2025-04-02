package markdown

import "strings"

func Split(text string, maxLen int) []string {
	var parts []string
	var currentPart strings.Builder
	currentLen := 0
	inCodeBlock := false
	inCodeBlockLine := ""

	lines := strings.Split(text, "\n")
	for _, line := range lines {

		hasCodePrefix := strings.HasPrefix(line, "```")
		if hasCodePrefix {
			inCodeBlock = !inCodeBlock
			if inCodeBlock {
				inCodeBlockLine = line
			} else {
				inCodeBlockLine = ""
			}
		}
		lineLen := len(line) + 1

		if inCodeBlock {
			if currentLen+lineLen+5 > maxLen {
				currentPart.WriteString("\n```\n")
				parts = append(parts, currentPart.String())
				currentPart.Reset()

				currentLen = len(inCodeBlockLine) + 1
				currentPart.WriteString(inCodeBlockLine + "\n")
			}
		} else {
			if currentLen+lineLen > maxLen {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
				currentLen = 0
			}
		}
		currentLen += lineLen
		currentPart.WriteString(line + "\n")
	}
	if currentPart.Len() > 0 {
		parts = append(parts, strings.TrimRight(currentPart.String(), "\n"))
	}
	return parts
}
