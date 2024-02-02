package config

import (
	"io/ioutil"
	"os"

	"testing"

	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestCreateNewLitmusCtlConfig(t *testing.T) {
	testFilename := "test-config.yaml"
	testConfig := types.LitmuCtlConfig{
		APIVersion: "v1",
		Kind:       "Config",
	}

	// unit test
	err := CreateNewLitmusCtlConfig(testFilename, testConfig)
	assert.NoError(t, err)

	assert.True(t, FileExists(testFilename))

	// Check if the file length is greater than 0
	length, err := GetFileLength(testFilename)
	assert.NoError(t, err)
	assert.Greater(t, length, 0)

	os.Remove(testFilename)
}

func TestFileExists(t *testing.T) {

	testFilename := "test-file-exists.txt"
	_, err := os.Create(testFilename)
	assert.NoError(t, err)

	exists := FileExists(testFilename)
	assert.True(t, exists)

	os.Remove(testFilename)
}

func TestGetFileLength(t *testing.T) {

	testFilename := "test-file-length.txt"
	file, err := os.Create(testFilename)
	assert.NoError(t, err)
	_, err = file.WriteString("This is a test.")
	assert.NoError(t, err)
	file.Close()

	length, err := GetFileLength(testFilename)
	assert.NoError(t, err)

	// Calculate the expected length based on the text content, accounting for line endings
	expectedLength := len("This is a test.")
	assert.Equal(t, expectedLength, length)

	err = os.Remove(testFilename)
	assert.NoError(t, err)
}

func TestYamltoObject(t *testing.T) {

	testConfig := types.LitmuCtlConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		Accounts:       []types.Account{},
		CurrentAccount: "test-account",
		CurrentUser:    "test-user",
	}

	// Serialize the test config to YAML and write it to a temporary test file
	data, err := yaml.Marshal(testConfig)
	assert.NoError(t, err)
	testFilename := "test-config.yaml"
	err = ioutil.WriteFile(testFilename, data, 0644)
	assert.NoError(t, err)

	//unit test
	obj, err := YamltoObject(testFilename)
	assert.NoError(t, err)
	assert.Equal(t, testConfig, obj)

	os.Remove(testFilename)
}

func TestConfigSyntaxCheck(t *testing.T) {

	testFilename := "test-config.yaml"
	testConfig := types.LitmuCtlConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		Accounts:       []types.Account{},
		CurrentAccount: "",
		CurrentUser:    "",
	}

	// Serialize the LitmusCtlConfig to a YAML file
	err := CreateNewLitmusCtlConfig(testFilename, testConfig)
	assert.NoError(t, err)

	//unit test
	err = ConfigSyntaxCheck(testFilename)
	assert.NoError(t, err)

	os.Remove(testFilename)

	invalidTestFilename := "invalid-test-config.yaml"
	invalidTestConfig := types.LitmuCtlConfig{
		APIVersion:     "v1",
		Kind:           "InvalidKind", // This is an invalid Kind value
		Accounts:       []types.Account{},
		CurrentAccount: "",
		CurrentUser:    "",
	}

	err = CreateNewLitmusCtlConfig(invalidTestFilename, invalidTestConfig)
	assert.NoError(t, err)

	//  unit test
	err = ConfigSyntaxCheck(invalidTestFilename)
	assert.Error(t, err)

	os.Remove(invalidTestFilename)
}

func TestUpdateLitmusCtlConfig(t *testing.T) {

	testFilename := "test-update-config.yaml"
	initialConfig := types.LitmuCtlConfig{
		APIVersion: "v1",
		Kind:       "Config",
		Accounts: []types.Account{
			{
				Endpoint: "https://cluster1",
				Users: []types.User{
					{
						Username:  "user1",
						Token:     "token1",
						ExpiresIn: "3600s",
					},
				},
			},
			{
				Endpoint: "https://cluster2",
				Users: []types.User{
					{
						Username:  "user2",
						Token:     "token2",
						ExpiresIn: "7200s",
					},
				},
			},
		},
		CurrentAccount: "cluster1",
		CurrentUser:    "user1",
	}

	// Serialize the initial LitmusCtlConfig to a YAML file
	err := CreateNewLitmusCtlConfig(testFilename, initialConfig)
	assert.NoError(t, err)

	// Define an update configuration
	updateConfig := types.UpdateLitmusCtlConfig{
		Account: types.Account{
			Endpoint: "https://cluster2",
			Users: []types.User{
				{
					Username:  "user2",    // Existing username
					Token:     "newtoken", // Updated token
					ExpiresIn: "5400s",    // Updated ExpiresIn
				},
			},
		},
		CurrentAccount: "cluster2", // Updated current account
		CurrentUser:    "user2",    // Updated current user
	}

	// unit test to update the configuration
	err = UpdateLitmusCtlConfig(updateConfig, testFilename)
	assert.NoError(t, err)

	// Read the updated configuration from the file
	updatedConfig, err := YamltoObject(testFilename)
	assert.NoError(t, err)

	// Check if the user's token and ExpiresIn have been updated
	assert.Equal(t, "newtoken", updatedConfig.Accounts[1].Users[0].Token)
	assert.Equal(t, "5400s", updatedConfig.Accounts[1].Users[0].ExpiresIn)
	// Check if the current account and user have been updated
	assert.Equal(t, "cluster2", updatedConfig.CurrentAccount)
	assert.Equal(t, "user2", updatedConfig.CurrentUser)

	os.Remove(testFilename)
}

func TestUpdateCurrent(t *testing.T) {

	testFilename := "test-update-current.yaml"
	defer os.Remove(testFilename)

	initialConfig := types.LitmuCtlConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentUser:    "user1",
		CurrentAccount: "account1",
		Accounts: []types.Account{
			{
				Endpoint: "example.com",
				Users: []types.User{
					{
						Username:  "user1",
						Token:     "token1",
						ExpiresIn: "3600s",
					},
				},
			},
		},
	}

	err := CreateNewLitmusCtlConfig(testFilename, initialConfig)
	assert.NoError(t, err)

	// updated current user and account
	updatedCurrent := types.Current{
		CurrentUser:    "user2",
		CurrentAccount: "account2",
	}

	// Update the current user and account
	err = UpdateCurrent(updatedCurrent, testFilename)
	assert.NoError(t, err)

	updatedConfig, err := YamltoObject(testFilename)
	assert.NoError(t, err)

	// Check that the current user and account have been updated as expected
	assert.Equal(t, "user2", updatedConfig.CurrentUser)
	assert.Equal(t, "account2", updatedConfig.CurrentAccount)
}

func TestWriteObjToFile(t *testing.T) {

	testFilename := "test-write-obj.yaml"
	testConfig := types.LitmuCtlConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentUser:    "user1",
		CurrentAccount: "account1",
		Accounts: []types.Account{
			{
				Endpoint: "https://cluster1",
				Users: []types.User{
					{
						Username:  "user1",
						Token:     "token1",
						ExpiresIn: "3600s",
					},
				},
			},
		},
	}

	err := WriteObjToFile(testConfig, testFilename)
	assert.NoError(t, err)

	loadedConfig, err := YamltoObject(testFilename)
	assert.NoError(t, err)

	// Verify that the loaded configuration matches the original test configuration
	assert.Equal(t, testConfig, loadedConfig)

	os.Remove(testFilename)
}

func TestIsAccountExists(t *testing.T) {

	testConfig := types.LitmuCtlConfig{
		Accounts: []types.Account{
			{
				Endpoint: "example.com",
				Users: []types.User{
					{
						Username:  "user1",
						Token:     "token1",
						ExpiresIn: "3600s",
					},
					{
						Username:  "user2",
						Token:     "token2",
						ExpiresIn: "7200s",
					},
				},
			},
			{
				Endpoint: "another.com",
				Users: []types.User{
					{
						Username:  "user3",
						Token:     "token3",
						ExpiresIn: "5400s",
					},
				},
			},
		},
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentAccount: "example.com",
		CurrentUser:    "user1",
	}

	// Check if a known username and endpoint exist
	exists := IsAccountExists(testConfig, "user1", "example.com")
	assert.True(t, exists)

	// Check if a username that doesn't exist returns false
	exists = IsAccountExists(testConfig, "nonexistentuser", "example.com")
	assert.False(t, exists)

	// Check if an endpoint that doesn't exist returns false
	exists = IsAccountExists(testConfig, "user1", "nonexistent.com")
	assert.False(t, exists)
}
