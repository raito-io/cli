package naming_hint

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/raito-io/cli/base/access_provider/sync_to_target"
)

type UniqueGenerator interface {

	// Generate create unique and consistent name that can be used in the target data source to create access elements.
	// The argument ap is an access provider pointer. Each next ap that is provided to the class should have an ascending ID to guarantee a deterministic ID generation
	// The output is a map that will point AccessIDs to the unique generate name
	Generate(ap *sync_to_target.AccessProvider) (map[string]string, error)
}

type uniqueGenerator struct {
	logger         hclog.Logger
	prefix         string
	constraints    *AllowedCharacters
	splitCharacter rune
	translator     Translator
	existingNames  map[string]uint
}

// NewUniqueGenerator will create an implementation of the UniqueGenerator interface. The UniqueGenerator will ensure the constraints provided in the first argument
func NewUniqueGenerator(logger hclog.Logger, prefix string, constraints *AllowedCharacters) (UniqueGenerator, error) {
	if constraints.splitCharacter() == 0 {
		return nil, errors.New("no support for UniqueGenerator if no split character is defined")
	}

	if !constraints.Numbers {
		return nil, errors.New("no support if numbers are not allowed")
	}

	if constraints.MaxLength < 8 {
		return nil, errors.New("no support if maximum characters is less than 8")
	}

	translator, err := NewNameHintTranslator(constraints)
	if err != nil {
		return nil, err
	}

	return &uniqueGenerator{
		logger:         logger,
		prefix:         prefix,
		constraints:    constraints,
		translator:     translator,
		splitCharacter: constraints.splitCharacter(),
		existingNames:  make(map[string]uint),
	}, nil
}

func (g *uniqueGenerator) Generate(ap *sync_to_target.AccessProvider) (map[string]string, error) {
	maxLength := g.constraints.MaxLength - 6

	var nameHinting string
	if ap.NamingHint != "" {
		nameHinting = ap.NamingHint
	} else {
		nameHinting = ap.Name
	}

	name, err := g.translator.Translate(g.prefix + nameHinting)
	if err != nil {
		return nil, err
	}

	if uint(len(name)) > maxLength {
		name = name[:maxLength]
	}

	accessElements := make(map[string]*sync_to_target.Access)
	accessElementIds := make([]string, len(ap.Access))

	for i := range ap.Access {
		access := ap.Access[i]
		accessName := access.Id
		accessElements[accessName] = access

		accessElementIds[i] = accessName
	}

	sort.Strings(accessElementIds)

	result := make(map[string]string)

	for _, accessId := range accessElementIds {
		access := accessElements[accessId]

		if access.ActualName != nil {
			originalNameSplit := strings.Split(*access.ActualName, fmt.Sprintf("%[1]c%[1]c", g.splitCharacter))
			originalName := originalNameSplit[0]

			if originalName == name && len(originalNameSplit) > 1 {
				idNumber, err := strconv.ParseInt(originalNameSplit[1], 16, 16)

				if err == nil {
					number := uint(idNumber)

					if _, found := g.existingNames[name]; !found || g.existingNames[name] < number {
						g.existingNames[name] = number
					}
				} else {
					g.logger.Warn("Error while parsing id from actualName. Will ignore actualName")
				}
			}
		}

		if currentNumber, found := g.existingNames[name]; found {
			postfixId := fmt.Sprintf("%[1]c%[1]c%[2]x", g.splitCharacter, currentNumber)

			if !g.constraints.UpperCaseLetters {
				postfixId = strings.ToLower(postfixId)
			}

			g.existingNames[name] = currentNumber + 1

			result[accessId] = fmt.Sprintf("%s%s", name, postfixId)
		} else {
			g.existingNames[name] = 0
			result[accessId] = name
		}
	}

	return result, nil
}
