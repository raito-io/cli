package cmd

import (
	"testing"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_moreThanOneExecutionWithinAnHour(t *testing.T) {
	type args struct {
		cronExpression string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "At every minute",
			args: args{
				cronExpression: "* * * * *",
			},
			want: true,
		},
		{
			name: "At minute 0",
			args: args{
				cronExpression: "0 * * * *",
			},
			want: false,
		},
		{
			name: "At every minute past hour 1 on day-of-month 2 in April",
			args: args{
				cronExpression: "* 1 2 4 *",
			},
			want: true,
		},
		{
			name: "At every 15th minute from 0 through 59",
			args: args{
				cronExpression: "0/15 * * * *",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule, err := cron.ParseStandard(tt.args.cronExpression)
			require.NoError(t, err)

			specSchedule := schedule.(*cron.SpecSchedule)

			assert.Equalf(t, tt.want, moreThanOneExecutionWithinAnHour(specSchedule), "moreThanOneExecutionWithinAnHour(%v)", tt.args.cronExpression)
		})
	}
}
