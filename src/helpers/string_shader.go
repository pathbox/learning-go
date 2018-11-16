package helpers

import (
	"fmt"
	"math"
	"strings"
)

// StringShader returns the string with shader etc: UserName BusinessCode IdentityNo for fuzzy
func StringShader(content, contentType string) string {
	if content == "" {
		return ""
	}
	if contentType == "" {
		return content
	}

	var shaderResult string
	var offSet int
	runeStr := []rune(content)
	switch contentType {
	case "UserName":
		shaderResult = fmt.Sprintf("%s**", string(runeStr[:1]))
	case "BusinessCode":
		shaderResult = fmt.Sprintf("%s*********", string(runeStr[:6]))
	case "CreditCode":
		shaderResult = fmt.Sprintf("%s**********", string(runeStr[:8]))
	case "IdentityNo":
		shaderResult = fmt.Sprintf("%s**********%s", string(runeStr[:1]), string(runeStr[len(runeStr)-1:]))
	case "BankAccount":
		offSet = MinInt(len(runeStr), 4)
		shaderResult = fmt.Sprintf("**************%s", string(runeStr[len(runeStr)-offSet:]))
	case "AlipayAccount":
		indexOfAt := strings.LastIndex(content, "@")
		if indexOfAt == -1 {
			offSet = MinInt(len(runeStr), 4)
			shaderResult = fmt.Sprintf("*******%s", string(runeStr[len(runeStr)-offSet:]))
		} else {
			offSet = MinInt(indexOfAt, 2)
			shaderResult = fmt.Sprintf("*********%s", string(runeStr[indexOfAt-offSet:]))
		}
	}

	return shaderResult
}

func MinInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
