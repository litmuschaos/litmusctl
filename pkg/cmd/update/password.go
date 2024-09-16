package update

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var PasswordCmd = &cobra.Command{
	Use: "password",
	Short: `Updates an account's password.
		Examples(s)
		#update a user's password
		litmusctl update password
		`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			updatePasswordRequest types.UpdatePasswordInput
		)

		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		promptUsername := promptui.Prompt{
			Label: "Username",
		}
		updatePasswordRequest.Username, err = promptUsername.Run()
		utils.PrintError(err)

		promptOldPassword := promptui.Prompt{
			Label: "Old Password",
			Mask:  '*',
		}
		updatePasswordRequest.OldPassword, err = promptOldPassword.Run()
		utils.PrintError(err)
	
	NEW_PASSWORD:

		promptNewPassword := promptui.Prompt{
			Label: "New Password",
			Mask:  '*',
		}
		updatePasswordRequest.NewPassword, err = promptNewPassword.Run()
		utils.PrintError(err)

		promptConfirmPassword := promptui.Prompt{
			Label: "Confirm New Password",
			Mask:  '*',
		}
		confirmPassword, err := promptConfirmPassword.Run()
		utils.PrintError(err)

		if updatePasswordRequest.NewPassword != confirmPassword {
			utils.Red.Println("\nPasswords do not match. Please try again.")
			goto NEW_PASSWORD
		}
		payloadBytes, _ := json.Marshal(updatePasswordRequest)

		resp, err := apis.SendRequest(
			apis.SendRequestParams{
				Endpoint: credentials.Endpoint + utils.AuthAPIPath + "/update/password/",
				Token:    "Bearer " + credentials.Token,
			},
			payloadBytes,
			string(types.Post),
		)
		utils.PrintError(err)
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		utils.PrintError(err)

		if resp.StatusCode == http.StatusOK {
			utils.White_B.Println("\nPassword updated successfully!")
		} else {
			err := errors.New("Unmatched status code: " + string(bodyBytes))
			utils.Red.Println("\nError updating password: ", err)
		}
	},
}

func init() {
	UpdateCmd.AddCommand(PasswordCmd)
}