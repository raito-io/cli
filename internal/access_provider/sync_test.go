package access_provider

import (
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindLastUpdated(t *testing.T) {
	time, err := findLastCalculated("./testdata/input-ok.yaml", hclog.L())
	assert.NoError(t, err)
	assert.Equal(t, int64(2222), time)

	time, err = findLastCalculated("./testdata/input-okwithspaces.yaml", hclog.L())
	assert.NoError(t, err)
	assert.Equal(t, int64(2222), time)

	time, err = findLastCalculated("./testdata/input-invalid.yaml", hclog.L())
	assert.Error(t, err)
	assert.Equal(t, int64(0), time)

	time, err = findLastCalculated("./testdata/input-morethan10.yaml", hclog.L())
	assert.NoError(t, err)
	assert.Equal(t, int64(0), time)

	time, err = findLastCalculated("./testdata/input-notpresent.yaml", hclog.L())
	assert.NoError(t, err)
	assert.Equal(t, int64(0), time)
}
