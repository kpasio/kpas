// This contains commands for generating the files needed to create
// a kpas compatible cluster on Hetzner Cloud.

package cmd

import (
	"fmt"
	"io/ioutil"
	"kpas/tpl"
	"path/filepath"

	vault "github.com/sosedoff/ansible-vault-go"
	"github.com/spf13/cobra"
)

var hCloudApiToken string
var hCloudSSHKeyName string
var vaultPassword string
var baseDomain string

// hcloudCmd represents the hcloud command
var hcloudCmd = &cobra.Command{
	Use:   "hcloud REPO_NAME CLUSTER_NAME",
	Short: "Commands for generating hcloud clusters",
	Long:  `Commands for generating hcloud clusters`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hcloud called")

		basePath := repoPath()
		repoPath := filepath.Join(basePath, args[0])
		// @TODO check if repo exists, exit if not
		// @TODO check if cluster exists, exit if so

		// Create the inventory directory
		inventoryPath := filepath.Join(repoPath, "inventory")
		createDiretoryIfNotExist(inventoryPath)

		// Collate cluster information
		hCloudClusterInfo := tpl.HCloudClusterInfo{
			ClusterName:    args[1],
			HCloudApiToken: hCloudApiToken,
			SshKeyName:     hCloudSSHKeyName,
		}

		genericClusterInfo := tpl.GenericClusterInfo{
			ClusterName:          args[1],
			CloudFlareDnsEnabled: false,
			BaseDomain:           baseDomain,
		}

		// Write out the inventory
		inventoryContent := tpl.HCloudInventory(hCloudClusterInfo)
		err := ioutil.WriteFile(filepath.Join(inventoryPath, "inventory.yml"), inventoryContent, 0644)
		CheckIfError(err)

		// Generate the hcloud specific secrets encrypted with ansible vault
		providerSecrets := tpl.HCloudProviderVault(hCloudClusterInfo)
		err = vault.EncryptFile(filepath.Join(repoPath, "provider_vault.yml"), string(providerSecrets), vaultPassword)
		CheckIfError(err)

		// Generate the generic secrets encrypted with ansible vault
		sharedSecrets := tpl.SharedVault(genericClusterInfo)
		err = vault.EncryptFile(filepath.Join(repoPath, "shared_vault.yml"), string(sharedSecrets), vaultPassword)
		CheckIfError(err)

		// Generate the provider values
		providerValues := tpl.HCloudProviderValues(hCloudClusterInfo)
		err = ioutil.WriteFile(filepath.Join(repoPath, "provider_values.yml"), providerValues, 0644)
		CheckIfError(err)

		// Generate the the shared values
		sharedValues := tpl.SharedValues(genericClusterInfo)
		err = ioutil.WriteFile(filepath.Join(repoPath, "shared_values.yml"), sharedValues, 0644)
		CheckIfError(err)
	},
}

func init() {
	generateCmd.AddCommand(hcloudCmd)
	hcloudCmd.Flags().StringVarP(&hCloudApiToken, "hcloud-api-token", "", "", "HCloud API Token (required)")
	hcloudCmd.MarkFlagRequired("hcloud-api-token")
	hcloudCmd.Flags().StringVarP(&hCloudSSHKeyName, "ssh-key-name", "", "", "HCloud SSH Key Name (required)")
	hcloudCmd.MarkFlagRequired("hcloud-ssh-key-name")
	hcloudCmd.Flags().StringVarP(&vaultPassword, "vault-password", "", "", "Password that will be used to encrypt secrets")
	hcloudCmd.MarkFlagRequired("vault-password")
	hcloudCmd.Flags().StringVarP(&baseDomain, "base-domain", "", "", "Base domain which has a wildcard DNS entry pointing at the clusters ingress")
	hcloudCmd.MarkFlagRequired("base-domain")
}
