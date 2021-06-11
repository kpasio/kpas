package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"

	"kpas/tpl"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Commands for interacting with kpas repositories",
	Long: `In kpas a repository is a grouping of clusters stored as
a simple git repository. Sensitive files such as kubeconfigs
are encrypted in the repository and only decrypted locally.

When you add, push and pull repositories, behind the scenes
you are just pushing and pulling ordinary git repos in the
kpas repository directory, by default ~/.kpas/repos/

Encryption and decryption is handled using ansible-vault.`,
}

var addCmd = &cobra.Command{
	Use:   "add NAME REMOTE",
	Short: "Add a new remote (git) repository for storing details about clusters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Adding repo called %s from source %s\n", args[0], args[1])

		clonePath := filepath.Join(repoPath(), args[0])

		// Check if a folder with the repo name exists
		if _, err := os.Stat(clonePath); os.IsNotExist(err) {
			// if it does not then clone the repo into the repos subdirectory
			git.PlainClone(clonePath, false, &git.CloneOptions{
				URL:      args[1],
				Progress: os.Stdout,
			})

			// and decrypt the kubeconf files if any exist
		} else {
			// if it does then raise an error that the repo already exists
			fmt.Fprintln(os.Stderr, "A repository with this name already exists")
			os.Exit(1)
		}

	},
	Args: cobra.ExactArgs(2),
}

var initCmd = &cobra.Command{
	Use:   "init NAME REMOTE",
	Short: "Init a new remote (git) repository for storing details about clusters",
	Run: func(cmd *cobra.Command, args []string) {
		initPath := filepath.Join(repoPath(), args[0])

		// Check if a folder with the repo name exists
		if _, err := os.Stat(initPath); os.IsNotExist(err) {
			// if it does not exist then init a new repo
			repo, _ := git.PlainInit(initPath, false)

			// Add the remote path as origin
			_, err = repo.CreateRemote(&config.RemoteConfig{
				Name: "origin",
				URLs: []string{args[1]},
			})
			CheckIfError(err)

			// create a default gitignore file
			gitIgnoreContent := tpl.GitIgnoreTemplate()
			ioutil.WriteFile(filepath.Join(initPath, ".gitignore"), gitIgnoreContent, 0644)

			// commit default gitignore
			tree, err := repo.Worktree()
			CheckIfError(err)

			err = tree.AddWithOptions(&git.AddOptions{
				All: true,
			})
			CheckIfError(err)

			_, err = tree.Commit("Initial Commit", &git.CommitOptions{
				Author: &object.Signature{
					Name:  "kpas",
					Email: "kpas@kpas.io",
					When:  time.Now(),
				},
			})
			CheckIfError(err)

			// remind the user to push
			fmt.Fprintln(os.Stdout,
				"Repository",
				args[0],
				"created with remote",
				args[1],
				"don't forget to push with `kpas remote push",
				args[0],
				"`",
			)
		} else {
			// if it does then raise an error that the repo already exists
			fmt.Fprintln(os.Stderr, "A repository with this name already exists", initPath)
			os.Exit(1)
		}
	},
	Args: cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(addCmd)
	repoCmd.AddCommand(initCmd)
}
