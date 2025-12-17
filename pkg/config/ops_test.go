package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/litmuschaos/litmusctl/pkg/types"
)

func TestFileExists(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_file_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tmpDir, err := os.MkdirTemp("", "test_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "existing file",
			filename: tmpFile.Name(),
			want:     true,
		},
		{
			name:     "non-existing file",
			filename: "/path/to/non/existing/file.txt",
			want:     false,
		},
		{
			name:     "directory instead of file",
			filename: tmpDir,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FileExists(tt.filename)
			if got != tt.want {
				t.Errorf("FileExists(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestGetFileLength(t *testing.T) {
	content := "Hello, LitmusChaos!"
	tmpFile, err := os.CreateTemp("", "test_length_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	tests := []struct {
		name     string
		filename string
		want     int
		wantErr  bool
	}{
		{
			name:     "file with known content",
			filename: tmpFile.Name(),
			want:     len(content),
			wantErr:  false,
		},
		{
			name:     "non-existing file",
			filename: "/path/to/non/existing/file.txt",
			want:     -1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileLength(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileLength(%q) error = %v, wantErr %v", tt.filename, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFileLength(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestYamltoObject(t *testing.T) {
	testdataDir := filepath.Join("testdata")

	tests := []struct {
		name    string
		file    string
		wantErr bool
		check   func(t *testing.T, cfg types.LitmuCtlConfig)
	}{
		{
			name:    "valid config file",
			file:    filepath.Join(testdataDir, "valid_config.yaml"),
			wantErr: false,
			check: func(t *testing.T, cfg types.LitmuCtlConfig) {
				if cfg.APIVersion != "v1" {
					t.Errorf("Expected APIVersion 'v1', got '%s'", cfg.APIVersion)
				}
				if cfg.Kind != "Config" {
					t.Errorf("Expected Kind 'Config', got '%s'", cfg.Kind)
				}
				if cfg.CurrentUser != "testuser" {
					t.Errorf("Expected CurrentUser 'testuser', got '%s'", cfg.CurrentUser)
				}
				if cfg.CurrentAccount != "https://litmus.example.com" {
					t.Errorf("Expected CurrentAccount 'https://litmus.example.com', got '%s'", cfg.CurrentAccount)
				}
				if len(cfg.Accounts) != 2 {
					t.Errorf("Expected 2 accounts, got %d", len(cfg.Accounts))
				}
			},
		},
		{
			name:    "non-existing file",
			file:    filepath.Join(testdataDir, "nonexistent.yaml"),
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := YamltoObject(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("YamltoObject(%q) error = %v, wantErr %v", tt.file, err, tt.wantErr)
				return
			}
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func TestConfigSyntaxCheck(t *testing.T) {
	testdataDir := filepath.Join("testdata")

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid config syntax",
			file:    filepath.Join(testdataDir, "valid_config.yaml"),
			wantErr: false,
		},
		{
			name:    "invalid config syntax",
			file:    filepath.Join(testdataDir, "invalid_config.yaml"),
			wantErr: true,
		},
		{
			name:    "non-existing file",
			file:    filepath.Join(testdataDir, "nonexistent.yaml"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConfigSyntaxCheck(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigSyntaxCheck(%q) error = %v, wantErr %v", tt.file, err, tt.wantErr)
			}
		})
	}
}

func TestIsAccountExists(t *testing.T) {
	config := types.LitmuCtlConfig{
		Accounts: []types.Account{
			{
				Endpoint: "https://litmus.example.com",
				Users: []types.User{
					{Username: "testuser", Token: "token1"},
					{Username: "admin", Token: "token2"},
				},
			},
			{
				Endpoint: "https://staging.litmus.io",
				Users: []types.User{
					{Username: "devuser", Token: "token3"},
				},
			},
		},
	}

	tests := []struct {
		name     string
		username string
		endpoint string
		want     bool
	}{
		{
			name:     "existing user on correct endpoint",
			username: "testuser",
			endpoint: "https://litmus.example.com",
			want:     true,
		},
		{
			name:     "existing user on wrong endpoint",
			username: "testuser",
			endpoint: "https://staging.litmus.io",
			want:     false,
		},
		{
			name:     "non-existing user",
			username: "nonexistent",
			endpoint: "https://litmus.example.com",
			want:     false,
		},
		{
			name:     "non-existing endpoint",
			username: "testuser",
			endpoint: "https://unknown.endpoint.com",
			want:     false,
		},
		{
			name:     "admin user exists",
			username: "admin",
			endpoint: "https://litmus.example.com",
			want:     true,
		},
		{
			name:     "devuser on staging",
			username: "devuser",
			endpoint: "https://staging.litmus.io",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAccountExists(config, tt.username, tt.endpoint)
			if got != tt.want {
				t.Errorf("IsAccountExists(config, %q, %q) = %v, want %v",
					tt.username, tt.endpoint, got, tt.want)
			}
		})
	}
}

func TestCreateNewLitmusCtlConfig(t *testing.T) {
	// Create a temporary file path
	tmpDir, err := os.MkdirTemp("", "config_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "test_config.yaml")

	testConfig := types.LitmuCtlConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentAccount: "https://test.example.com",
		CurrentUser:    "testuser",
		Accounts: []types.Account{
			{
				Endpoint:       "https://test.example.com",
				ServerEndpoint: "https://test.example.com/api",
				Users: []types.User{
					{
						Username:  "testuser",
						Token:     "test-token",
						ExpiresIn: "1735689600",
					},
				},
			},
		},
	}

	// Test creating new config
	err = CreateNewLitmusCtlConfig(configPath, testConfig)
	if err != nil {
		t.Fatalf("CreateNewLitmusCtlConfig() error = %v", err)
	}

	if !FileExists(configPath) {
		t.Error("Config file was not created")
	}

	cfg, err := YamltoObject(configPath)
	if err != nil {
		t.Fatalf("Failed to read created config: %v", err)
	}

	if cfg.APIVersion != testConfig.APIVersion {
		t.Errorf("APIVersion mismatch: got %q, want %q", cfg.APIVersion, testConfig.APIVersion)
	}
	if cfg.CurrentUser != testConfig.CurrentUser {
		t.Errorf("CurrentUser mismatch: got %q, want %q", cfg.CurrentUser, testConfig.CurrentUser)
	}
}

func TestUpdateLitmusCtlConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "update_config_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "test_config.yaml")

	initialConfig := types.LitmuCtlConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentAccount: "https://initial.example.com",
		CurrentUser:    "initialuser",
		Accounts: []types.Account{
			{
				Endpoint:       "https://initial.example.com",
				ServerEndpoint: "https://initial.example.com/api",
				Users: []types.User{
					{Username: "initialuser", Token: "initial-token", ExpiresIn: "1735689600"},
				},
			},
		},
	}

	err = CreateNewLitmusCtlConfig(configPath, initialConfig)
	if err != nil {
		t.Fatalf("Failed to create initial config: %v", err)
	}

	t.Run("update existing user token", func(t *testing.T) {
		updateConfig := types.UpdateLitmusCtlConfig{
			CurrentAccount: "https://initial.example.com",
			CurrentUser:    "initialuser",
			ServerEndpoint: "https://initial.example.com/api",
			Account: types.Account{
				Endpoint: "https://initial.example.com",
				Users: []types.User{
					{Username: "initialuser", Token: "updated-token", ExpiresIn: "1800000000"},
				},
			},
		}

		err := UpdateLitmusCtlConfig(updateConfig, configPath)
		if err != nil {
			t.Fatalf("UpdateLitmusCtlConfig() error = %v", err)
		}

		cfg, err := YamltoObject(configPath)
		if err != nil {
			t.Fatalf("Failed to read updated config: %v", err)
		}

		if cfg.Accounts[0].Users[0].Token != "updated-token" {
			t.Errorf("Token not updated: got %q, want %q", cfg.Accounts[0].Users[0].Token, "updated-token")
		}
	})

	t.Run("add new user to existing endpoint", func(t *testing.T) {
		updateConfig := types.UpdateLitmusCtlConfig{
			CurrentAccount: "https://initial.example.com",
			CurrentUser:    "newuser",
			ServerEndpoint: "https://initial.example.com/api",
			Account: types.Account{
				Endpoint: "https://initial.example.com",
				Users: []types.User{
					{Username: "newuser", Token: "new-token", ExpiresIn: "1800000000"},
				},
			},
		}

		err := UpdateLitmusCtlConfig(updateConfig, configPath)
		if err != nil {
			t.Fatalf("UpdateLitmusCtlConfig() error = %v", err)
		}

		cfg, err := YamltoObject(configPath)
		if err != nil {
			t.Fatalf("Failed to read updated config: %v", err)
		}

		if len(cfg.Accounts[0].Users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(cfg.Accounts[0].Users))
		}
	})

	t.Run("add new endpoint", func(t *testing.T) {
		updateConfig := types.UpdateLitmusCtlConfig{
			CurrentAccount: "https://new.example.com",
			CurrentUser:    "newendpointuser",
			ServerEndpoint: "https://new.example.com/api",
			Account: types.Account{
				Endpoint:       "https://new.example.com",
				ServerEndpoint: "https://new.example.com/api",
				Users: []types.User{
					{Username: "newendpointuser", Token: "new-endpoint-token", ExpiresIn: "1800000000"},
				},
			},
		}

		err := UpdateLitmusCtlConfig(updateConfig, configPath)
		if err != nil {
			t.Fatalf("UpdateLitmusCtlConfig() error = %v", err)
		}

		cfg, err := YamltoObject(configPath)
		if err != nil {
			t.Fatalf("Failed to read updated config: %v", err)
		}

		if len(cfg.Accounts) != 2 {
			t.Errorf("Expected 2 accounts, got %d", len(cfg.Accounts))
		}
	})
}

func TestUpdateCurrent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "update_current_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "test_config.yaml")

	initialConfig := types.LitmuCtlConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentAccount: "https://old.example.com",
		CurrentUser:    "olduser",
		Accounts: []types.Account{
			{
				Endpoint: "https://old.example.com",
				Users:    []types.User{{Username: "olduser", Token: "token"}},
			},
			{
				Endpoint: "https://new.example.com",
				Users:    []types.User{{Username: "newuser", Token: "token2"}},
			},
		},
	}

	err = CreateNewLitmusCtlConfig(configPath, initialConfig)
	if err != nil {
		t.Fatalf("Failed to create initial config: %v", err)
	}

	t.Run("switch current account and user", func(t *testing.T) {
		current := types.Current{
			CurrentAccount: "https://new.example.com",
			CurrentUser:    "newuser",
		}

		err := UpdateCurrent(current, configPath)
		if err != nil {
			t.Fatalf("UpdateCurrent() error = %v", err)
		}

		cfg, err := YamltoObject(configPath)
		if err != nil {
			t.Fatalf("Failed to read updated config: %v", err)
		}

		if cfg.CurrentAccount != "https://new.example.com" {
			t.Errorf("CurrentAccount not updated: got %q", cfg.CurrentAccount)
		}
		if cfg.CurrentUser != "newuser" {
			t.Errorf("CurrentUser not updated: got %q", cfg.CurrentUser)
		}
	})

	t.Run("error on non-existing file", func(t *testing.T) {
		current := types.Current{
			CurrentAccount: "https://new.example.com",
			CurrentUser:    "newuser",
		}

		err := UpdateCurrent(current, "/nonexistent/path/config.yaml")
		if err == nil {
			t.Error("Expected error for non-existing file")
		}
	})
}

