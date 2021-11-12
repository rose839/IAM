package errors

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"basic", args{"hello"}, fmt.Errorf("hello")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := New(tt.args.message); err.Error() != tt.wantErr.Error() {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
