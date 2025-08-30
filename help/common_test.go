package help

import "testing"

func TestFormatFileSize(t *testing.T) {
	type args struct {
		fileSize int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "t1", args: args{fileSize: 10000}, want: "9.77 KB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFileSize(tt.args.fileSize); got != tt.want {
				t.Errorf("FormatFileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
