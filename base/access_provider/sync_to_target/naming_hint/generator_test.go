package naming_hint

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

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
	constraints := AllowedCharacters{
		UpperCaseLetters:  false,
		LowerCaseLetters:  true,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueGenerator{
		logger:         logger,
		constraints:    &constraints,
		splitCharacter: '_',
		existingNames:  map[string]uint{},
		translator:     &translatorMock{},
	}

	names := map[*sync_to_target.AccessProvider]map[string]string{
		&sync_to_target.AccessProvider{
			Id:         "SomeID",
			NamingHint: "the_first_access_provider",
			Access: []*sync_to_target.Access{
				{
					Id:         "AccessId1",
					ActualName: nil,
				},
			},
		}: {"AccessId1": "the_first_access_provider"},
		&sync_to_target.AccessProvider{
			Id:         "SomeID2",
			NamingHint: "second_access_provider",
			Access: []*sync_to_target.Access{
				{
					Id:         "AccessId2",
					ActualName: nil,
				},
			},
		}: {"AccessId2": "second_access_provider"},
		&sync_to_target.AccessProvider{
			Id:         "SomeID3",
			NamingHint: "and_the_last_access_provider",
			Access: []*sync_to_target.Access{
				{
					Id:         "AccessId3",
					ActualName: nil,
				},
			},
		}: {"AccessId3": "and_the_last_access_provid"},
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

	constraints := AllowedCharacters{
		UpperCaseLetters:  false,
		LowerCaseLetters:  true,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueGenerator{
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
			Access: []*sync_to_target.Access{
				{
					Id:         fmt.Sprintf("AccessId%d", id),
					ActualName: nil,
				},
			},
		}
	}

	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "the_same_name", output[ap.Access[0].Id])

	//WHEN + THEN
	for i := 0; i < 24; i++ {
		ap = accessProviderGenerator(i + 1)
		output, err = generator.Generate(ap)
		name := output[ap.Access[0].Id]

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
	upperCaseRegex := regexp.MustCompile("[A-Z0-9_]+")

	constraints := AllowedCharacters{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueGenerator{
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
			Access: []*sync_to_target.Access{
				{
					Id:         fmt.Sprintf("AccessId%d", id),
					ActualName: nil,
				},
			},
		}
	}

	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME", output[ap.Access[0].Id])

	//WHEN + THEN
	for i := 0; i < 24; i++ {
		ap = accessProviderGenerator(i + 1)
		output, err = generator.Generate(ap)
		name := output[ap.Access[0].Id]

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

	constraints := AllowedCharacters{
		UpperCaseLetters:  false,
		LowerCaseLetters:  true,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         16,
	}
	generator := uniqueGenerator{
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
			Access: []*sync_to_target.Access{
				{
					Id:         fmt.Sprintf("AccessId%d", id),
					ActualName: nil,
				},
			},
		}
	}

	ap := accessProviderGenerator(0, "abcdefghij")
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "abcdefghij", output[ap.Access[0].Id])

	//WHEN + THEN
	for i := 0; i < 24; i++ {
		ap = accessProviderGenerator(i+1, "abcdefghijkl")
		output, err = generator.Generate(ap)
		name := output[ap.Access[0].Id]

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
	constraints := AllowedCharacters{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueGenerator{
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
			Access: []*sync_to_target.Access{
				{
					Id:         fmt.Sprintf("AccessId%d", id),
					ActualName: actualName,
				},
			},
		}
	}

	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME__0", output[ap.Access[0].Id])

	ap = accessProviderGenerator(1)
	output, err = generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME__2", output[ap.Access[0].Id])

	ap = accessProviderGenerator(2)
	output, err = generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME__4", output[ap.Access[0].Id])

	ap = accessProviderGenerator(3)
	output, err = generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME__5", output[ap.Access[0].Id])
}

func TestUniqueGenerator_Generate_ActualNamesNotEqualToNameHint(t *testing.T) {
	constraints := AllowedCharacters{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueGenerator{
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
			Access: []*sync_to_target.Access{
				{
					Id:         fmt.Sprintf("AccessId%d", id),
					ActualName: &actualName,
				},
			},
		}
	}

	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_NAME_HINT_TO_USE", output[ap.Access[0].Id])
}

func TestUniqueGenerator_Generate_MultipleAccessElements(t *testing.T) {
	constraints := AllowedCharacters{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueGenerator{
		logger:         logger,
		constraints:    &constraints,
		splitCharacter: '_',
		existingNames:  map[string]uint{},
		translator:     &translatorMock{},
	}

	ap := &sync_to_target.AccessProvider{
		Id:         fmt.Sprintf("SomeId"),
		NamingHint: "THE_NAME_HINT_TO_USE",
		Access: []*sync_to_target.Access{
			{
				Id:         "BD",
				ActualName: nil,
			},
			{
				Id:         "AB",
				ActualName: nil,
			},
		},
	}

	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_NAME_HINT_TO_USE", output["AB"])
	assert.Equal(t, "THE_NAME_HINT_TO_USE__0", output["BD"])
}

func TestUniqueGenerator_Generate_InvalidActualName(t *testing.T) {
	constraints := AllowedCharacters{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		SpecialCharacters: "_-!",
		Numbers:           true,
		MaxLength:         32,
	}
	generator := uniqueGenerator{
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
		Access: []*sync_to_target.Access{
			{
				Id:         "BD",
				ActualName: &actualNamme,
			},
		},
	}

	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_NAME_HINT_TO_USE", output["BD"])
}

func TestUniqueGeneratorIT_Generate(t *testing.T) {
	//Given
	constraints := AllowedCharacters{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		Numbers:           true,
		MaxLength:         32,
		SpecialCharacters: "_-@#$",
	}

	generator, err := NewUniqueGenerator(logger, "", &constraints)

	assert.NoError(t, err)

	names := map[*sync_to_target.AccessProvider]map[string]string{
		&sync_to_target.AccessProvider{
			Id:         "SomeID",
			NamingHint: "the_first_access_provider",
			Access: []*sync_to_target.Access{
				{
					Id:         "AccessId1",
					ActualName: nil,
				},
			},
		}: {"AccessId1": "THE_FIRST_ACCESS_PROVIDER"},
		&sync_to_target.AccessProvider{
			Id:         "SomeID2",
			NamingHint: "second_access_provider",
			Access: []*sync_to_target.Access{
				{
					Id:         "AccessId2",
					ActualName: nil,
				},
			},
		}: {"AccessId2": "SECOND_ACCESS_PROVIDER"},
		&sync_to_target.AccessProvider{
			Id:         "SomeID3",
			NamingHint: "and_the_last_access_provider",
			Access: []*sync_to_target.Access{
				{
					Id:         "AccessId3",
					ActualName: nil,
				},
			},
		}: {"AccessId3": "AND_THE_LAST_ACCESS_PROVID"},
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

	constraints := AllowedCharacters{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		Numbers:           true,
		MaxLength:         32,
		SpecialCharacters: "_-@#$",
	}

	generator, err := NewUniqueGenerator(logger, "", &constraints)

	assert.NoError(t, err)

	existingNames := map[string]struct{}{}

	accessProviderGenerator := func(id int) *sync_to_target.AccessProvider {
		return &sync_to_target.AccessProvider{
			Id:         fmt.Sprintf("SomeID%d", id),
			NamingHint: "the_same_name",
			Access: []*sync_to_target.Access{
				{
					Id:         fmt.Sprintf("AccessId%d", id),
					ActualName: nil,
				},
			},
		}
	}

	//WHEN + THEN
	ap := accessProviderGenerator(0)
	output, err := generator.Generate(ap)
	assert.NoError(t, err)
	assert.Equal(t, "THE_SAME_NAME", output[ap.Access[0].Id])

	//WHEN + THEN
	for i := 0; i < 24; i++ {
		ap = accessProviderGenerator(i + 1)
		output, err = generator.Generate(ap)
		name := output[ap.Access[0].Id]

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

func TestUniqueGeneratorIT_Generate_WithPrefix(t *testing.T) {
	//Given
	constraints := AllowedCharacters{
		UpperCaseLetters:  true,
		LowerCaseLetters:  false,
		Numbers:           true,
		MaxLength:         32,
		SpecialCharacters: "_-@#$",
	}

	generator, err := NewUniqueGenerator(logger, "prefix_", &constraints)

	assert.NoError(t, err)

	names := map[*sync_to_target.AccessProvider]map[string]string{
		&sync_to_target.AccessProvider{
			Id:         "SomeID",
			NamingHint: "the_first_access_provider",
			Access: []*sync_to_target.Access{
				{
					Id:         "AccessId1",
					ActualName: nil,
				},
			},
		}: {"AccessId1": "PREFIX_THE_FIRST_ACCESS_PR"},
		&sync_to_target.AccessProvider{
			Id:         "SomeID2",
			NamingHint: "second_access_provider",
			Access: []*sync_to_target.Access{
				{
					Id:         "AccessId2",
					ActualName: nil,
				},
			},
		}: {"AccessId2": "PREFIX_SECOND_ACCESS_PROVI"},
		&sync_to_target.AccessProvider{
			Id:         "SomeID3",
			NamingHint: "and_the_last_access_provider",
			Access: []*sync_to_target.Access{
				{
					Id:         "AccessId3",
					ActualName: nil,
				},
			},
		}: {"AccessId3": "PREFIX_AND_THE_LAST_ACCESS"},
	}

	//WHEN + THEN
	for input, expectedOutput := range names {
		output, err := generator.Generate(input)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output)
	}
}
