package naming_hint

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raito-io/cli/base/access_provider/sync_to_target"
)

type translatorMock struct {
}

var logger = hclog.L()

func (m *translatorMock) Translate(input string) (string, error) {
	return input, nil
}

func TestUniqueGenerator_Generate_NoDuplicatedNames(t *testing.T) {
	//Given
	constraints := NamingConstraints{
		UpperCaseLetters:  false,
		LowerCaseLetters:  true,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueNameGenerator{
		logger:         logger,
		constraints:    &constraints,
		splitCharacter: '_',
		existingNames:  map[string]uint{},
		translator:     &translatorMock{},
	}

	names := map[*sync_to_target.AccessProvider]string{
		&sync_to_target.AccessProvider{
			Id:         "SomeID",
			NamingHint: "the_first_access_provider",
			ActualName: nil,
		}: "the_first_access_provider",
		&sync_to_target.AccessProvider{
			Id:         "SomeID2",
			NamingHint: "second_access_provider",
			ActualName: nil,
		}: "second_access_provider",
		&sync_to_target.AccessProvider{
			Id:         "SomeID3",
			NamingHint: "and_the_last_access_provider",
			ActualName: nil,
		}: "and_the_last_access_provid",
	}

	//WHEN + THEN
	for input, expectedOutput := range names {
		output, err := generator.Generate(input)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output)
	}
}

func TestUniqueGenerator_Generate_DuplicatedNames(t *testing.T) {
	//Given
	lowerCaseRegex := regexp.MustCompile("[a-z0-9_]+")

	constraints := NamingConstraints{
		UpperCaseLetters:  false,
		LowerCaseLetters:  true,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueNameGenerator{
		logger:         logger,
		constraints:    &constraints,
		splitCharacter: '_',
		existingNames:  map[string]uint{},
		translator:     &translatorMock{},
	}

	existingNames := map[string]struct{}{}

	accessProviderGenerator := func(id int) *sync_to_target.AccessProvider {
		return &sync_to_target.AccessProvider{
			Id:         fmt.Sprintf("SomeID%d", id),
			NamingHint: "the_same_name",
			ActualName: nil,
		}
	}

	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "the_same_name", output)

	//WHEN + THEN
	for i := 0; i < 24; i++ {
		ap = accessProviderGenerator(i + 1)
		output, err = generator.Generate(ap)
		name := output

		assert.NoError(t, err)
		assert.Equal(t, "the_same_name__", name[:15])

		_, found := existingNames[name]
		assert.False(t, found)

		if found {
			return
		}

		assert.True(t, lowerCaseRegex.MatchString(name))

		existingNames[name] = struct{}{}

	}
}

func TestUniqueGenerator_Generate_DuplicatedNames_uppercase(t *testing.T) {
	//Given
	upperCaseRegex := regexp.MustCompile("^[A-Z0-9_]+$")

	constraints := NamingConstraints{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueNameGenerator{
		logger:         logger,
		constraints:    &constraints,
		splitCharacter: '_',
		existingNames:  map[string]uint{},
		translator:     &translatorMock{},
	}

	existingNames := map[string]struct{}{}

	accessProviderGenerator := func(id int) *sync_to_target.AccessProvider {
		return &sync_to_target.AccessProvider{
			Id:         fmt.Sprintf("SomeID%d", id),
			NamingHint: "THE_SAME_NAME",
			ActualName: nil,
		}
	}

	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME", output)

	//WHEN + THEN
	for i := 0; i < 24; i++ {
		ap = accessProviderGenerator(i + 1)
		output, err = generator.Generate(ap)
		name := output

		assert.NoError(t, err)
		assert.Equal(t, "THE_SAME_NAME__", name[:15])

		_, found := existingNames[name]
		assert.False(t, found)

		if found {
			return
		}

		assert.True(t, upperCaseRegex.MatchString(name))

		existingNames[name] = struct{}{}

	}
}

func TestUniqueGenerator_Generate_LongAndAlreadyExistingNames(t *testing.T) {
	//Given
	lowerCaseRegex := regexp.MustCompile("[a-z0-9_]+")

	constraints := NamingConstraints{
		UpperCaseLetters:  false,
		LowerCaseLetters:  true,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         16,
	}
	generator := uniqueNameGenerator{
		logger:         logger,
		constraints:    &constraints,
		splitCharacter: '_',
		existingNames:  map[string]uint{},
		translator:     &translatorMock{},
	}

	existingNames := map[string]struct{}{}

	accessProviderGenerator := func(id int, namingHint string) *sync_to_target.AccessProvider {
		return &sync_to_target.AccessProvider{
			Id:         fmt.Sprintf("SomeID%d", id),
			NamingHint: namingHint,
			ActualName: nil,
		}
	}

	ap := accessProviderGenerator(0, "abcdefghij")
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "abcdefghij", output)

	//WHEN + THEN
	for i := 0; i < 24; i++ {
		ap = accessProviderGenerator(i+1, "abcdefghijkl")
		output, err = generator.Generate(ap)
		name := output

		assert.NoError(t, err)
		assert.Equal(t, "abcdefghij__", name[:12])

		_, found := existingNames[name]
		assert.False(t, found)

		if found {
			return
		}

		assert.True(t, lowerCaseRegex.MatchString(name))

		existingNames[name] = struct{}{}
	}
}

func TestUniqueGenerator_Generate_ActualNamesExist(t *testing.T) {
	constraints := NamingConstraints{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueNameGenerator{
		logger:         logger,
		constraints:    &constraints,
		splitCharacter: '_',
		existingNames:  map[string]uint{},
		translator:     &translatorMock{},
	}

	accessProviderGenerator := func(id int) *sync_to_target.AccessProvider {
		nameHint := "THE_SAME_NAME"
		var actualName *string

		if id < 3 {
			actualNamePointee := fmt.Sprintf("%s__%x", nameHint, id*2)
			actualName = &actualNamePointee
		}

		return &sync_to_target.AccessProvider{
			Id:         fmt.Sprintf("SomeID%d", id),
			NamingHint: nameHint,
			ActualName: actualName,
		}
	}

	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME__0", output)

	ap = accessProviderGenerator(1)
	output, err = generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME__2", output)

	ap = accessProviderGenerator(2)
	output, err = generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME__4", output)

	ap = accessProviderGenerator(3)
	output, err = generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME__5", output)
}

func TestUniqueGenerator_Generate_ActualNamesNotEqualToNameHint(t *testing.T) {
	constraints := NamingConstraints{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueNameGenerator{
		logger:         logger,
		constraints:    &constraints,
		splitCharacter: '_',
		existingNames:  map[string]uint{},
		translator:     &translatorMock{},
	}

	accessProviderGenerator := func(id int) *sync_to_target.AccessProvider {
		actualName := "ORIGINAL_NAME__3"
		return &sync_to_target.AccessProvider{
			Id:         fmt.Sprintf("SomeID%d", id),
			NamingHint: "THE_NAME_HINT_TO_USE",
			ActualName: &actualName,
		}
	}

	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_NAME_HINT_TO_USE", output)
}

func TestUniqueGenerator_Generate_InvalidActualName(t *testing.T) {
	constraints := NamingConstraints{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueNameGenerator{
		logger:         logger,
		constraints:    &constraints,
		splitCharacter: '_',
		existingNames:  map[string]uint{},
		translator:     &translatorMock{},
	}

	actualNamme := "THE_NAME_HINT_TO_USE__RTE"
	ap := &sync_to_target.AccessProvider{
		Id:         fmt.Sprintf("SomeId"),
		NamingHint: "THE_NAME_HINT_TO_USE",
		ActualName: &actualNamme,
	}

	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_NAME_HINT_TO_USE", output)
}

func TestUniqueGeneratorIT_Generate(t *testing.T) {
	//Given
	constraints := NamingConstraints{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		Numbers:           true,
		MaxLength:         32,
		SpecialCharacters: "_-@#$",
	}

	generator, err := NewUniqueNameGenerator(logger, "", &constraints)

	assert.NoError(t, err)

	names := map[*sync_to_target.AccessProvider]string{
		&sync_to_target.AccessProvider{
			Id:         "SomeID",
			NamingHint: "the_first_access_provider",
			ActualName: nil,
		}: "THE_FIRST_ACCESS_PROVIDER",
		&sync_to_target.AccessProvider{
			Id:         "SomeID2",
			NamingHint: "second_access_provider",
			ActualName: nil,
		}: "SECOND_ACCESS_PROVIDER",
		&sync_to_target.AccessProvider{
			Id:         "SomeID3",
			NamingHint: "and_the_last_access_provider",
			ActualName: nil,
		}: "AND_THE_LAST_ACCESS_PROVID",
	}

	//WHEN + THEN
	for input, expectedOutput := range names {
		output, err := generator.Generate(input)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output)
	}
}

func TestUniqueGeneratorIT_Generate_DuplicatedNames_uppercase(t *testing.T) {
	//Given
	upperCaseRegex := regexp.MustCompile("[A-Z0-9_]+")

	constraints := NamingConstraints{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		Numbers:           true,
		MaxLength:         32,
		SpecialCharacters: "_-@#$",
	}

	generator, err := NewUniqueNameGenerator(logger, "", &constraints)

	assert.NoError(t, err)

	existingNames := map[string]struct{}{}

	accessProviderGenerator := func(id int) *sync_to_target.AccessProvider {
		return &sync_to_target.AccessProvider{
			Id:         fmt.Sprintf("SomeID%d", id),
			NamingHint: "the_same_name",
			ActualName: nil,
		}
	}

	//WHEN + THEN
	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME", output)

	//WHEN + THEN
	for i := 0; i < 24; i++ {
		ap = accessProviderGenerator(i + 1)
		output, err = generator.Generate(ap)
		name := output

		assert.NoError(t, err)
		assert.Equal(t, "THE_SAME_NAME__", name[:15])

		_, found := existingNames[name]
		assert.False(t, found)

		if found {
			return
		}

		assert.True(t, upperCaseRegex.MatchString(name))

		existingNames[name] = struct{}{}
	}
}

func TestUniqueGeneratorIT_Generate_DuplicatedNames_uppercase_ImportedAccessProvider(t *testing.T) {
	//Given
	constraints := NamingConstraints{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		Numbers:           true,
		MaxLength:         32,
		SpecialCharacters: "_-@#$",
	}

	generator, err := NewUniqueNameGenerator(logger, "", &constraints)

	assert.NoError(t, err)

	accessProviderGenerator := func(id int, namingHint string) *sync_to_target.AccessProvider {
		return &sync_to_target.AccessProvider{
			Id:         fmt.Sprintf("SomeID%d", id),
			NamingHint: namingHint,
			ActualName: nil,
		}
	}

	importedAccessProviderGenerator := func(id int, actualName string) *sync_to_target.AccessProvider {
		return &sync_to_target.AccessProvider{
			Id:         fmt.Sprintf("ExternalId%d", id),
			NamingHint: actualName,
			ActualName: &actualName,
		}
	}
	ap := accessProviderGenerator(0, "the_same_name")
	output, err := generator.Generate(ap)
	require.NoError(t, err)
	require.Equal(t, "THE_SAME_NAME", output)

	// IMPORTED ACCESS PROVIDER
	importedActualName := "Some Actual Name w^th Non-generated  char"
	ap = importedAccessProviderGenerator(0, importedActualName)
	output, err = generator.Generate(ap)
	require.NoError(t, err)
	assert.Equal(t, importedActualName, output)

	// Not able to import the same name
	ap = importedAccessProviderGenerator(1, "THE_SAME_NAME")
	output, err = generator.Generate(ap)
	require.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME__0", output)

	// generate new name for non imported access provider
	ap = accessProviderGenerator(1, "Some Actual Name w^th Non-generated  char")
	output, err = generator.Generate(ap)
	require.NoError(t, err)
	assert.Equal(t, "SOME_ACTUAL_NAME_W_TH_NON-", output)

}

func TestUniqueGeneratorIT_Generate_WithPrefix(t *testing.T) {
	//Given
	constraints := NamingConstraints{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		Numbers:           true,
		MaxLength:         32,
		SpecialCharacters: "_-@#$",
	}

	generator, err := NewUniqueNameGenerator(logger, "prefix_", &constraints)

	assert.NoError(t, err)

	names := map[*sync_to_target.AccessProvider]string{
		&sync_to_target.AccessProvider{
			Id:         "SomeID",
			NamingHint: "the_first_access_provider",
			ActualName: nil,
		}: "PREFIX_THE_FIRST_ACCESS_PR",
		&sync_to_target.AccessProvider{
			Id:         "SomeID2",
			NamingHint: "second_access_provider",
			ActualName: nil,
		}: "PREFIX_SECOND_ACCESS_PROVI",
		&sync_to_target.AccessProvider{
			Id:         "SomeID3",
			NamingHint: "and_the_last_access_provider",
			ActualName: nil,
		}: "PREFIX_AND_THE_LAST_ACCESS",
	}

	//WHEN + THEN
	for input, expectedOutput := range names {
		output, err := generator.Generate(input)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output)
	}
}
