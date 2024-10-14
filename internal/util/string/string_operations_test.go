package string

import "testing"

func TestTrimSpaceInCommaSeparatedList(t *testing.T) {
	type args struct {
		list string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "NoList",
			args: args{
				list: "SomeRandomString with spaces",
			},
			want: "SomeRandomString with spaces",
		},
		{
			name: "List without spaces",
			args: args{
				list: "SomeRandomString,without,spaces",
			},
			want: "SomeRandomString,without,spaces",
		},
		{
			name: "List with spaces",
			args: args{
				list: "SomeRandomString, without , spaces ",
			},
			want: "SomeRandomString,without,spaces",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimSpaceInCommaSeparatedList(tt.args.list); got != tt.want {
				t.Errorf("TrimSpaceInCommaSeparatedList() = %v, want %v", got, tt.want)
			}
		})
	}
}
