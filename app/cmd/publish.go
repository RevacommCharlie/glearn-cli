package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gSchool/glearn-cli/api/learn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	branchCommand     = `git branch | grep \* | cut -d ' ' -f2`
	remoteNameCommand = `git remote -v | grep push | cut -f2- -d/ | sed 's/[.].*$//'`
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish master for your curriculum repository",
	Long: `
		The Learn system recognizes blocks of content held in GitHub respositories.
		This command pushes the latest commit for the remote origin master (which
		should be GitHub), then attemptes the release of a new Learn block version
		at the HEAD of master.
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if viper.Get("api_token") == "" || viper.Get("api_token") == nil {
			previewCmdError("Please set your API token first with `learn set --api_token=value`")
		}

		if len(args) != 0 {
			fmt.Println("Usage: `learn publish` takes no arguments, merely pushing latest master and releasing a version to Learn")
			os.Exit(1)
		}

		// Start benchmarking the total time spent in publish cmd
		startOfCmd := time.Now()

		remote, err := remoteName()
		if err != nil {
			log.Printf("Cannot run git remote detection with command: git remote -v | grep push | cut -f2- -d/ | sed 's/[.].*$//'\n%s", err)
			os.Exit(1)
		}
		if remote == "" {
			log.Println("no fetch remote detected")
			os.Exit(1)
		}

		block, err := learn.API.GetBlockByRepoName(remote)
		if err != nil {
			log.Printf("Error fetchng block from learn: %s", err)
			os.Exit(1)
		}
		if !block.Exists() {
			block, err = learn.API.CreateBlockByRepoName(remote)
			if err != nil {
				log.Printf("Error creating block from learn: %s", err)
				os.Exit(1)
			}
		}

		branch, err := currentBranch()
		if err != nil {
			log.Println("Cannot run git branch detection with bash:", err)
			os.Exit(1)
		}
		if branch != "master" {
			fmt.Println("You are currently not on branch 'master'- the `learn publish` command must be on master branch to push all currently committed work to your 'origin master' remote.")
			os.Exit(1)
		}

		fmt.Println("Pushing work to remote origin", branch)

		// TODO what happens when they do not have work in remote and push fails?
		err = pushToRemote(branch)
		if err != nil {
			fmt.Printf("Error pushing to origin remote on branch: %s", err)
			os.Exit(1)
		}

		// Start benchmark for creating master release & building on learn
		startOfMasterReleaseAndBuild := time.Now()

		// Create a release on learn, notify user
		releaseID, err := learn.API.CreateMasterRelease(block.ID)
		if err != nil || releaseID == 0 {
			fmt.Printf("error creating master release for releaseID: %d. Error: %s", releaseID, err)
			os.Exit(1)
		}

		var attempts uint8 = 20
		_, err = learn.API.PollForBuildResponse(releaseID, &attempts)
		if err != nil {
			fmt.Printf("Failed to fetch the state of your build: %s", err)
			os.Exit(1)
		}

		// Add benchmark in milliseconds for compressDirectory
		bench := &learn.CLIBenchmark{
			MasterReleaseAndBuild: time.Since(startOfMasterReleaseAndBuild).Milliseconds(),
			TotalCmdTime:          time.Since(startOfCmd).Milliseconds(),
			CmdName:               "publish",
		}

		fmt.Printf("Block %d released!\n", block.ID)

		err = learn.API.SendMetadataToLearn(&learn.CLIBenchmarkPayload{
			CLIBenchmark: bench,
		})
		if err != nil {
			learn.API.NotifySlack(err)
			os.Exit(1)
		}
	},
}

func currentBranch() (string, error) {
	return runBashCommand(branchCommand)
}

func remoteName() (string, error) {
	return runBashCommand(remoteNameCommand)
}

func runBashCommand(command string) (string, error) {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func pushToRemote(branch string) error {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("git push origin %s", branch)).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}