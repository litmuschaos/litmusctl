package version

import (
	"testing"
)

func TestVersionCmd_Properties(t *testing.T) {
	t.Run("Use field is correct", func(t *testing.T) {
		if VersionCmd.Use != "version" {
			t.Errorf("VersionCmd.Use = %q, want %q", VersionCmd.Use, "version")
		}
	})

	t.Run("Short description is set", func(t *testing.T) {
		if VersionCmd.Short == "" {
			t.Error("VersionCmd.Short should not be empty")
		}
		expected := "Displays the version of litmusctl"
		if VersionCmd.Short != expected {
			t.Errorf("VersionCmd.Short = %q, want %q", VersionCmd.Short, expected)
		}
	})

	t.Run("Command is runnable", func(t *testing.T) {
		if VersionCmd.Run == nil {
			t.Error("VersionCmd.Run should not be nil")
		}
	})
}

func TestUpdateCmd_Properties(t *testing.T) {
	t.Run("Use field is correct", func(t *testing.T) {
		if UpdateCmd.Use != "update" {
			t.Errorf("UpdateCmd.Use = %q, want %q", UpdateCmd.Use, "update")
		}
	})

	t.Run("Short description is set", func(t *testing.T) {
		if UpdateCmd.Short == "" {
			t.Error("UpdateCmd.Short should not be empty")
		}
		expected := "Changes the version of litmusctl"
		if UpdateCmd.Short != expected {
			t.Errorf("UpdateCmd.Short = %q, want %q", UpdateCmd.Short, expected)
		}
	})

	t.Run("Requires exactly one argument", func(t *testing.T) {
		if UpdateCmd.Args == nil {
			t.Error("UpdateCmd.Args should not be nil")
		}
	})

	t.Run("Command is runnable", func(t *testing.T) {
		if UpdateCmd.Run == nil {
			t.Error("UpdateCmd.Run should not be nil")
		}
	})
}

func TestVersionCmd_Subcommands(t *testing.T) {
	t.Run("UpdateCmd is subcommand of VersionCmd", func(t *testing.T) {
		found := false
		for _, cmd := range VersionCmd.Commands() {
			if cmd.Use == "update" {
				found = true
				break
			}
		}
		if !found {
			t.Error("UpdateCmd should be a subcommand of VersionCmd")
		}
	})

	t.Run("VersionCmd has exactly one subcommand", func(t *testing.T) {
		cmdCount := len(VersionCmd.Commands())
		if cmdCount != 1 {
			t.Errorf("VersionCmd has %d subcommands, want 1", cmdCount)
		}
	})
}

func TestVersionCmd_Help(t *testing.T) {
	t.Run("Command has help available", func(t *testing.T) {
		if VersionCmd.UsageString() == "" {
			t.Error("Usage string should not be empty")
		}
	})
}

func TestUpdateCmd_ArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "one argument (valid)",
			args:    []string{"0.23.0"},
			wantErr: false,
		},
		{
			name:    "two arguments",
			args:    []string{"0.23.0", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateCmd.Args(UpdateCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCmd.Args() with %v error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
		})
	}
}
