package naming_hint

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/raito-io/cli/base/access_provider/sync_to_target"
)

type UniqueGenerator interface {

	// Generate create unique and consistent name that can be used in the target data source to create access elements.
	// The argument ap is an access provider pointer. Each next ap that is provided to the class should have an ascending ID to guarantee a deterministic ID generation
	// The output is a map that will point AccessIDs to the unique generate name
	Generate(ap *sync_to_target.AccessProvider) (string, error)
}

// uniqueNameGenerator implements a Generate method which generates unique names that can be used to create access elements.
// If defined, a prefix is used when generating the names.
// If the naming hint of an accessProvider is specified, the naming hint will be reformed to a name only consisting valid characters.
// After that, a validation is executed to check if the valid name should be post-fixed with a unique ID.
// The post-fixes start with two splitCharacters and end with a 4 character hexadecimal number. Note that this is the only place that 2 splitCharacters can be used after each other.
// The algorithm has the prerequisite that accessProviders are provided in ascending order by ID.
// The algorithm will check if an already existing postfix exists in the actual name. If that is the case, the current post-fix will be used.
// Later created access providers will end up with a higher post-fix if they are reusing the same valid name.
// Example:
//
//		constraints: Uppercase, Numbers, '_', maxLength: 16
//		- AP{namingHint: "lowerCaseNamingHint"} => "LOWER_CASE"
//	 	- AP{namingHint: "lowerCaseNaming2"} => "LOWER_CASE__0"
//	 	- AP{namingHint: "lowerCaseNaming3"} => "LOWER_CASE__1"
//	 	- AP{namingHint: "UPPER_CASE_HINT", actualName: "UPPER_CASE__3"} => "UPPER_CASE__3
//	 	- AP{name: "UPPER_CASE_HINT2"} => "UPPER_CASE__4
type uniqueNameGenerator struct {
	logger         hclog.Logger
	prefix         string
	constraints    *NamingConstraints
	splitCharacter rune
	translator     Translator
	existingNames  map[string]uint
}

// NewUniqueNameGenerator will create an implementation of the UniqueGenerator interface. The UniqueGenerator will ensure the constraints provided in the first argument
func NewUniqueNameGenerator(logger hclog.Logger, prefix string, constraints *NamingConstraints) (UniqueGenerator, error) {
	if constraints.SplitCharacter() == 0 {
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

	return &uniqueNameGenerator{
		logger:         logger,
		prefix:         prefix,
		constraints:    constraints,
		translator:     translator,
		splitCharacter: constraints.SplitCharacter(),
		existingNames:  make(map[string]uint),
	}, nil
}

func (g *uniqueNameGenerator) Generate(ap *sync_to_target.AccessProvider) (string, error) {
	// Reserve 6 character for post fix ID
	maxLength := g.constraints.MaxLength - 6

	var nameHinting string
	if ap.NamingHint != "" {
		nameHinting = ap.NamingHint
	} else {
		nameHinting = ap.Name
	}

	name, err := g.translator.Translate(g.prefix + nameHinting)
	if err != nil {
		return "", err
	}

	if uint(len(name)) > maxLength {
		name = name[:maxLength]
	}

	if ap.ActualName != nil && *ap.ActualName == ap.NamingHint {
		// This case is when external access provider is imported. We try to avoid renaming if possible.
		if _, found := g.existingNames[*ap.ActualName]; !found {
			// Actual name was not used before so we can use it directly
			g.existingNames[*ap.ActualName] = 0

			return *ap.ActualName, nil
		}
		// Else we will generate a new name
	} //nolint:wsl

	if ap.ActualName != nil {
		// Search for post fix ID
		originalNameSplit := strings.Split(*ap.ActualName, fmt.Sprintf("%[1]c%[1]c", g.splitCharacter))
		originalName := originalNameSplit[0]

		if len(originalNameSplit) > 2 {
			g.logger.Warn(fmt.Sprintf("Current actual name %q does not fit expected name pattern. Rename is required", *ap.ActualName))
		} else if originalName == name && len(originalNameSplit) > 1 {
			idNumber, err := strconv.ParseInt(originalNameSplit[1], 16, 16)

			if err == nil {
				number := uint(idNumber) //nolint:gosec

				if _, found := g.existingNames[name]; !found || g.existingNames[name] < number {
					g.existingNames[name] = number
				}
			} else {
				g.logger.Warn("Error while parsing id from actualName. Will ignore actualName")
			}
		}
	}

	if currentNumber, found := g.existingNames[name]; found {
		postfixId := fmt.Sprintf("%[1]c%[1]c%[2]X", g.splitCharacter, currentNumber)

		if !g.constraints.UpperCaseLetters {
			postfixId = strings.ToLower(postfixId)
		}

		g.existingNames[name] = currentNumber + 1

		name = fmt.Sprintf("%s%s", name, postfixId)
	} else {
		g.existingNames[name] = 0
	}

	g.logger.Info(fmt.Sprintf("Generate unique name for ap %q: %+v", ap.Name, name))

	return name, nil
}
