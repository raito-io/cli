package naming_hint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConstraints_validConstraints(t *testing.T) {
	constraints := NamingConstraints{
		LowerCaseLetters:  true,
		UpperCaseLetters:  false,
		Numbers:           false,
		SpecialCharacters: "-_#%$!",
		MaxLength:         128,
	}

	assert.NoError(t, validateConstraints(&constraints))

	translator, err := NewNameHintTranslator(&constraints)
	assert.NoError(t, err)
	assert.NotNil(t, translator)

}

func TestValidateConstraints_noAlphabeticCharacters(t *testing.T) {
	constraints := NamingConstraints{
		LowerCaseLetters:  false,
		UpperCaseLetters:  false,
		Numbers:           true,
		SpecialCharacters: "-_#%$!",
		MaxLength:         128,
	}

	assert.Error(t, validateConstraints(&constraints))

	translator, err := NewNameHintTranslator(&constraints)
	assert.Error(t, err)
	assert.Nil(t, translator)
}

func TestNameHintTranslator_Translate_AllowAll(t *testing.T) {
	constraints := NamingConstraints{
		LowerCaseLetters:  true,
		UpperCaseLetters:  true,
		Numbers:           true,
		SpecialCharacters: "_-!@#$%^&*()+=,.?/<> ",
		MaxLength:         128,
	}

	testParameters := map[string]string{
		"FirstStringWithUpperAndLowerCasing":                      "FirstStringWithUpperAndLowerCasing",
		"SNAKE_CASE_STRING":                                       "SNAKE_CASE_STRING",
		"second_snake_case_string":                                "second_snake_case_string",
		"@_Str!ng#With_$pecial.Charact*rs like ?/<> for example&": "Str!ng#With_$pecial.Charact*rs like ?/<> for example",
		"A string with spaces":                                    "A string with spaces",
		"A7tringW1th4Numbers2":                                    "A7tringW1th4Numbers2",
	}

	executeNameHintTranslatorTest(t, &constraints, testParameters)
}

func TestNameHintTranslator_TranslateToSnakeCase_UpperCharacter(t *testing.T) {
	constraints := NamingConstraints{
		LowerCaseLetters:  false,
		UpperCaseLetters:  true,
		Numbers:           true,
		SpecialCharacters: "_! ",
		MaxLength:         128,
	}

	testParameters := map[string]string{
		"FirstStringWithUpperAndLowerCasing":                     "FIRST_STRING_WITH_UPPER_AND_LOWER_CASING",
		"SNAKE_CASE_STRING":                                      "SNAKE_CASE_STRING",
		"second_snake_case_string":                               "SECOND_SNAKE_CASE_STRING",
		"@_Str!ng#With_$pecial.Charact*rs like ?/<> for example": "STR!NG_WITH_PECIAL_CHARACT_RS LIKE _ FOR EXAMPLE",
		"A string with spaces":                                   "A STRING WITH SPACES",
		"A7tringW1th4Numbers2":                                   "A7TRING_W1TH4_NUMBERS2",
	}

	executeNameHintTranslatorTest(t, &constraints, testParameters)
}

func TestNameHintTranslator_TranslateToSnakeCase_LowerCharacter(t *testing.T) {
	constraints := NamingConstraints{
		LowerCaseLetters:  true,
		UpperCaseLetters:  false,
		Numbers:           true,
		SpecialCharacters: "_!",
		MaxLength:         128,
	}

	testParameters := map[string]string{
		"FirstStringWithUpperAndLowerCasing":                     "first_string_with_upper_and_lower_casing",
		"SNAKE_CASE_STRING":                                      "snake_case_string",
		"second_snake_case_string":                               "second_snake_case_string",
		"@_Str!ng#With_$pecial.Charact*rs like ?/<> for example": "str!ng_with_pecial_charact_rs_like_for_example",
		"A string with spaces":                                   "a_string_with_spaces",
		"A7tringW1th4Numbers2":                                   "a7tring_w1th4_numbers2",
	}

	executeNameHintTranslatorTest(t, &constraints, testParameters)
}

func TestNameHintTranslator_TranslateWithoutSpecialCharacter(t *testing.T) {
	constraints := NamingConstraints{
		LowerCaseLetters:  true,
		UpperCaseLetters:  false,
		Numbers:           true,
		SpecialCharacters: "",
		MaxLength:         128,
	}

	testParameters := map[string]string{
		"FirstStringWithUpperAndLowerCasing":                     "firststringwithupperandlowercasing",
		"SNAKE_CASE_STRING":                                      "snakecasestring",
		"second_snake_case_string":                               "secondsnakecasestring",
		"@_Str!ng#With_$pecial.Charact*rs like ?/<> for example": "strngwithpecialcharactrslikeforexample",
		"A string with spaces":                                   "astringwithspaces",
		"A7tringW1th4Numbers2":                                   "a7tringw1th4numbers2",
	}

	executeNameHintTranslatorTest(t, &constraints, testParameters)
}

func TestNameHintTranslator_TranslateWithoutSpecialCharacter_lowerCase(t *testing.T) {
	constraints := NamingConstraints{
		LowerCaseLetters:  true,
		UpperCaseLetters:  false,
		Numbers:           true,
		SpecialCharacters: "",
		MaxLength:         128,
	}

	testParameters := map[string]string{
		"FirstStringWithUpperAndLowerCasing":                     "firststringwithupperandlowercasing",
		"SNAKE_CASE_STRING":                                      "snakecasestring",
		"second_snake_case_string":                               "secondsnakecasestring",
		"@_Str!ng#With_$pecial.Charact*rs like ?/<> for example": "strngwithpecialcharactrslikeforexample",
		"A string with spaces":                                   "astringwithspaces",
		"A7tringW1th4Numbers2":                                   "a7tringw1th4numbers2",
	}

	executeNameHintTranslatorTest(t, &constraints, testParameters)
}

func TestNameHintTranslator_TranslateWithoutSpecialCharacter_toCamel(t *testing.T) {
	constraints := NamingConstraints{
		LowerCaseLetters:  true,
		UpperCaseLetters:  true,
		Numbers:           false,
		SpecialCharacters: "",
		MaxLength:         128,
	}

	testParameters := map[string]string{
		"FirstStringWithUpperAndLowerCasing":                     "FirstStringWithUpperAndLowerCasing",
		"SNAKE_CASE_STRING":                                      "SNAKECASESTRING",
		"second_snake_case_string":                               "secondSnakeCaseString",
		"@_Str!ng#With_$pecial.Charact*rs like ?/<> for example": "StrNgWithPecialCharactRsLikeForExample",
		"A string with spaces":                                   "AStringWithSpaces",
		"A7tringW1th4Numbers2":                                   "ATringWThNumbers",
	}

	executeNameHintTranslatorTest(t, &constraints, testParameters)
}

func executeNameHintTranslatorTest(t *testing.T, constraints *NamingConstraints, testParams map[string]string) {
	translator, err := NewNameHintTranslator(constraints)
	assert.NoError(t, err)

	for input, expectedOutput := range testParams {
		t.Run(fmt.Sprintf("Translate(%q)=%q", input, expectedOutput), func(t *testing.T) {
			output, err := translator.Translate(input)
			assert.NoError(t, err)
			assert.Equal(t, expectedOutput, output)
		})
	}
}
