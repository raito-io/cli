package file

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUniqueFileNameForTarget(t *testing.T) {
	fileNames := make(map[string]struct{})
	for i := 0; i < 10; i++ {
		fileName := CreateUniqueFileNameForTarget("the // special ( name", "step-name", "yml")
		assert.True(t, strings.HasPrefix(fileName, "thespecialname-step-name-"), "Filename doesn't have the right prefix")
		assert.True(t, strings.HasSuffix(fileName, ".yml"), "Filename doesn't have the right suffix")
		fmt.Println(fileName)
		_, found := fileNames[fileName]
		assert.False(t, found, "Duplicate filename found ("+strconv.Itoa(i)+")")
		fileNames[fileName] = struct{}{}
		time.Sleep(1 * time.Millisecond)
	}

}
