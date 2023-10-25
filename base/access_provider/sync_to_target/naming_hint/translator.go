package naming_hint

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var splitRegex = regexp.MustCompile("([a-z])([A-Z])")

type Translator interface {
	Translate(input string) (string, error)
}

type nameHintTranslator struct {
	allowedCharacters *NamingConstraints
}

func NewNameHintTranslator(constraints *NamingConstraints) (Translator, error) {
	err := validateConstraints(constraints)
	if err != nil {
		return nil, err
	}

	return &nameHintTranslator{allowedCharacters: constraints}, nil
}

func validateConstraints(allowedCharacters *NamingConstraints) error {
	if !allowedCharacters.LowerCaseLetters && !allowedCharacters.UpperCaseLetters {
		return errors.New("no support for non alphabetic constraints")
	}

	return nil
}

func (t *nameHintTranslator) Translate(input string) (string, error) {
	result := input
	splitChar := t.allowedCharacters.SplitCharacter()

	if !t.allowedCharacters.LowerCaseLetters || !t.allowedCharacters.UpperCaseLetters {
		//Remove invalid casing
		if splitChar != 0 {
			result = splitRegex.ReplaceAllString(result, fmt.Sprintf("${1}%c${2}", splitChar))
		}

		if !t.allowedCharacters.LowerCaseLetters {
			result = strings.ToUpper(result)
		} else {
			result = strings.ToLower(result)
		}
	} else if splitChar == 0 {
		result = makeCamel(result)
	}

	//Remove invalid characters
	invalidCharRegex, err := t.allowedCharacters.nonAllowedCharacterRegex()
	if err != nil {
		return "", err
	}

	if splitChar == 0 {
		result = invalidCharRegex.ReplaceAllString(result, "")
	} else {
		result = invalidCharRegex.ReplaceAllString(result, fmt.Sprintf("%c", splitChar))
	}

	//Start with alphanumeric character
	invalidStartCharRegex, err := t.allowedCharacters.nonAlphaNumericBeginRegex()
	if err != nil {
		return "", err
	}

	result = invalidStartCharRegex.ReplaceAllString(result, "")

	//End with alphanumeric character
	invalidLastCharRegex, err := t.allowedCharacters.nonAlphaNumericEndRegex()
	if err != nil {
		return "", err
	}

	if splitChar != 0 {
		//Remove consecutive splitCharacters
		consecutiveSplitCharRegex, err := t.allowedCharacters.consecutiveSplitCharacter()
		if err != nil {
			return "", err
		}

		result = consecutiveSplitCharRegex.ReplaceAllString(result, fmt.Sprintf("%c", splitChar))
	}

	result = invalidLastCharRegex.ReplaceAllString(result, "")

	return result, nil
}

func quoteMeta(input string) string {
	specialChars := regexp.QuoteMeta(input)
	specialChars = strings.ReplaceAll(specialChars, "-", "\\-")

	return strings.ReplaceAll(specialChars, "_", "\\_")
}

func makeCamel(s string) string {
	s = strings.TrimSpace(s)

	sb := strings.Builder{}

	upperNext := false

	for _, char := range s {
		charIsUpper := char >= 'A' && char <= 'Z'
		charIsLower := char >= 'a' && char <= 'z'
		alphabeticChar := charIsUpper || charIsLower

		if upperNext && charIsLower {
			char += 'A' - 'a'
		}

		if alphabeticChar || char >= '0' && char <= '9' {
			//Normal character
			sb.WriteRune(char)
			upperNext = false
		}

		if !alphabeticChar {
			upperNext = true
		}
	}

	return sb.String()
}
