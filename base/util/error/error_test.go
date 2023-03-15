package error

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToErrorResult(t *testing.T) {
	res := ToErrorResult(errors.New("An error"))

	require.Equal(t, status.New(codes.Internal, "An error"), errorToStatus(t, res))
}

func TestCreateSourceConnectionError(t *testing.T) {
	res := CreateSourceConnectionError("http://myurl", "failed")

	require.Equal(t, "error while connecting to \"http://myurl\": failed", res.Error())
}

func TestCreateBadInputParameterError(t *testing.T) {
	res := CreateBadInputParameterError("param1", "v1", "explained")

	require.Equal(t, "parameter \"param1\" has invalid value \"v1\". explained", res.Error())
}

func TestCreateBadInputParameterErrorNoExplanation(t *testing.T) {
	res := CreateBadInputParameterError("param2", "v2", "")

	require.Equal(t, "parameter \"param2\" has invalid value \"v2\"", res.Error())
}

func TestCreateMissingInputParameterError(t *testing.T) {
	res := CreateMissingInputParameterError("param666")

	require.Equal(t, "mandatory parameter \"param666\" is missing", res.Error())
}

func errorToStatus(t *testing.T, err error) *status.Status {
	t.Helper()

	s, ok := status.FromError(err)
	require.True(t, ok)

	return s
}
