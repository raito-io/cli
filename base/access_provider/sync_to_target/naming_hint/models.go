package naming_hint

import (
	"fmt"
	"regexp"
	"strings"
)

const lowerCaseLetters = "a-z"
const upperCaseLetters = "A-Z"
const numbers = "0-9"

var splitCharacters = map[rune]struct{}{
	'_': {},
	'-': {},
	'#': {},
}

type AllowedCharacters struct {
	LowerCaseLetters  bool
	UpperCaseLetters  bool
	Numbers           bool
	SpecialCharacters string
	MaxLength         uint
}

func (a *AllowedCharacters) allowedCharacterRegex() (*regexp.Regexp, error) {
	var sb strings.Builder

	sb.WriteString("[")
	sb.WriteString(a.characterRegex())
	sb.WriteString("]")

	return regexp.Compile(sb.String())
}

func (a *AllowedCharacters) nonAllowedCharacterRegex() (*regexp.Regexp, error) {
	var sb strings.Builder

	sb.WriteString("[^")
	sb.WriteString(a.characterRegex())
	sb.WriteString("]+")

	return regexp.Compile(sb.String())
}

func (a *AllowedCharacters) nonAlphaNumericEndRegex() (*regexp.Regexp, error) {
	var sb strings.Builder

	sb.WriteString("[^")

	if a.LowerCaseLetters {
		sb.WriteString(lowerCaseLetters)
	}

	if a.UpperCaseLetters {
		sb.WriteString(upperCaseLetters)
	}

	if a.Numbers {
		sb.WriteString(numbers)
	}

	sb.WriteString("]+$")

	return regexp.Compile(sb.String())
}

func (a *AllowedCharacters) nonAlphaNumericBeginRegex() (*regexp.Regexp, error) {
	var sb strings.Builder

	sb.WriteString("^[^")

	if a.LowerCaseLetters {
		sb.WriteString(lowerCaseLetters)
	}

	if a.UpperCaseLetters {
		sb.WriteString(upperCaseLetters)
	}

	if a.Numbers {
		sb.WriteString(numbers)
	}

	sb.WriteString("]+")

	return regexp.Compile(sb.String())
}

func (a *AllowedCharacters) consecutiveSplitCharacter() (*regexp.Regexp, error) {
	var sb strings.Builder

	splitchar := a.splitCharacter()
	if splitchar == 0 {
		return nil, nil
	}

	sb.WriteString(quoteMeta(fmt.Sprintf("%c", splitchar)))
	sb.WriteString("{2,}")

	return regexp.Compile(sb.String())
}

func (a *AllowedCharacters) characterRegex() string {
	var sb strings.Builder

	if a.LowerCaseLetters {
		sb.WriteString(lowerCaseLetters)
	}

	if a.UpperCaseLetters {
		sb.WriteString(upperCaseLetters)
	}

	if a.Numbers {
		sb.WriteString(numbers)
	}

	specialChars := quoteMeta(a.SpecialCharacters)

	sb.WriteString(specialChars)

	return sb.String()
}

func (a *AllowedCharacters) splitCharacter() rune {
	for _, char := range a.SpecialCharacters {
		if _, found := splitCharacters[char]; found {
			return char
		}
	}

	return 0
}
