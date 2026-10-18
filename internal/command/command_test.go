package command_test

import (
	"bytes"
	"testing"

	"github.com/laurilyytikainen/go-bfc/internal/command"
)

func Test_parseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *command.Command
		wantErr bool
	}{
		{
			name: "Empty args return None command",
			args: []string{},
			want: &command.Command{
				SourceFilePath: "",
				Command:        command.None,
				OutputFilePath: "",
				Verbose:        false,
			},
			wantErr: false,
		},
		{
			name: "Run command with a file returns Run command",
			args: []string{"run", "test.b"},
			want: &command.Command{
				SourceFilePath: "test.b",
				Command:        command.Run,
				OutputFilePath: "",
				Verbose:        false,
			},
			wantErr: false,
		},
		{
			name: "Run command with verbose flag returns Run command",
			args: []string{"run", "--verbose", "test.b"},
			want: &command.Command{
				SourceFilePath: "test.b",
				Command:        command.Run,
				OutputFilePath: "",
				Verbose:        true,
			},
			wantErr: false,
		},
		{
			name: "Build command with a file returns Build command",
			args: []string{"build", "test.b"},
			want: &command.Command{
				SourceFilePath: "test.b",
				Command:        command.Build,
				OutputFilePath: "",
				Verbose:        false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			got, gotErr := command.ParseArgs(tt.args, &buf)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parseArgs() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("parseArgs() succeeded unexpectedly")
			}
			if !tt.want.Equals(got) || ((gotErr == nil) == tt.wantErr) {
				t.Errorf("parseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
