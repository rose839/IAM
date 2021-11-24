package idutil

import "testing"

func TestGetInstanceID(t *testing.T) {
	type args struct {
		uid    uint64
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				uid:    10,
				prefix: "secret-",
			},
			want: "secret-123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInstanceID(tt.args.uid, tt.args.prefix); got != tt.want {
				t.Errorf("GetInstanceID() = %v, want %v", got, tt.want)
			}
		})
	}
}
