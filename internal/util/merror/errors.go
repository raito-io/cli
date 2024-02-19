package merror

import (
	"strings"

	"github.com/raito-io/cli/internal/util/array"
)

type Errors []error

func (e Errors) Error() string {
	return strings.Join(array.Map(e, func(err *error) string { return (*err).Error() }), ",")
}

func Append(err error, errs ...error) error {
	if len(errs) == 0 {
		return err
	} else if err == nil && len(errs) == 1 {
		return errs[0]
	} else if err == nil {
		merr := &Errors{}
		*merr = append(*merr, errs...)

		return merr
	}

	merr, ok := err.(*Errors)
	if ok {
		*merr = append(*merr, errs...)
	} else {
		merr = &Errors{}
		*merr = append(*merr, err)
		*merr = append(*merr, errs...)
	}

	return merr
}
