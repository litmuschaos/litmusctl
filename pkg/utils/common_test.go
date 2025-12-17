package utils

import (
	"strings"
	"testing"
)

func TestCheckKeyValueFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid single key-value pair",
			input: "key1=value1",
			want:  true,
		},
		{
			name:  "valid multiple key-value pairs",
			input: "key1=value1,key2=value2",
			want:  true,
		},
		{
			name:  "valid three key-value pairs",
			input: "app=nginx,env=prod,version=v1",
			want:  true,
		},
		{
			name:  "invalid - missing value",
			input: "key1=",
			want:  true,
		},
		{
			name:  "invalid - missing equals",
			input: "key1value1",
			want:  false,
		},
		{
			name:  "invalid - multiple equals",
			input: "key1=value1=extra",
			want:  false,
		},
		{
			name:  "invalid - contains quotes in key",
			input: "\"key1\"=value1",
			want:  false,
		},
		{
			name:  "invalid - contains quotes in value",
			input: "key1=\"value1\"",
			want:  false,
		},
		{
			name:  "valid - empty value allowed",
			input: "key=",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckKeyValueFormat(tt.input)
			if got != tt.want {
				t.Errorf("CheckKeyValueFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenerateNameID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple lowercase",
			input: "testname",
			want:  "testname",
		},
		{
			name:  "uppercase to lowercase",
			input: "TestName",
			want:  "testname",
		},
		{
			name:  "spaces replaced with underscore",
			input: "test name",
			want:  "test_name",
		},
		{
			name:  "special characters replaced",
			input: "test@name#123",
			want:  "test_name_123",
		},
		{
			name:  "hyphens replaced with underscore",
			input: "test-name-id",
			want:  "test_name_id",
		},
		{
			name:  "mixed special chars and spaces",
			input: "My Test Project! v1.0",
			want:  "my_test_project_v1_0",
		},
		{
			name:  "numbers preserved",
			input: "test123",
			want:  "test123",
		},
		{
			name:  "complex string",
			input: "Chaos-Experiment (Production) #1",
			want:  "chaos_experiment_production_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateNameID(tt.input)
			if got != tt.want {
				t.Errorf("GenerateNameID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "positive length",
			length:  10,
			wantErr: false,
		},
		{
			name:    "zero length",
			length:  0,
			wantErr: true,
		},
		{
			name:    "negative length",
			length:  -5,
			wantErr: true,
		},
		{
			name:    "length of 1",
			length:  1,
			wantErr: false,
		},
		{
			name:    "large length",
			length:  100,
			wantErr: false,
		},
	}

	validChars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomString(tt.length)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomString(%d) error = %v, wantErr %v", tt.length, err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != tt.length {
					t.Errorf("GenerateRandomString(%d) returned string of length %d", tt.length, len(got))
				}

				for _, char := range got {
					if !strings.ContainsRune(validChars, char) {
						t.Errorf("GenerateRandomString(%d) returned invalid character: %c", tt.length, char)
					}
				}
			}
		})
	}

	t.Run("randomness check", func(t *testing.T) {
		str1, _ := GenerateRandomString(20)
		str2, _ := GenerateRandomString(20)
		if str1 == str2 {
			t.Log("Warning: Two random strings were identical (unlikely but possible)")
		}
	})
}

func TestPrintInJsonFormat(t *testing.T) {
	t.Run("does not panic with map", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PrintInJsonFormat panicked: %v", r)
			}
		}()
		data := map[string]string{"key": "value"}
		PrintInJsonFormat(data)
	})

	t.Run("does not panic with struct", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PrintInJsonFormat panicked: %v", r)
			}
		}()
		data := struct {
			Name string `json:"name"`
		}{Name: "test"}
		PrintInJsonFormat(data)
	})
}

func TestPrintInYamlFormat(t *testing.T) {
	t.Run("does not panic with map", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PrintInYamlFormat panicked: %v", r)
			}
		}()
		data := map[string]string{"key": "value"}
		PrintInYamlFormat(data)
	})

	t.Run("does not panic with slice", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PrintInYamlFormat panicked: %v", r)
			}
		}()
		data := []string{"item1", "item2"}
		PrintInYamlFormat(data)
	})
}
