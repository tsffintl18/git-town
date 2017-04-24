package cmd

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/steps"
	"github.com/spf13/cobra"
)

type newPullRequestConfig struct {
	InitialBranch  string
	BranchesToSync []string
}

var newPullRequestCommand = &cobra.Command{
	Use:   "new-pull-request",
	Short: "Creates a new pull request",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "new-pull-request",
			IsAbort:              abortFlag,
			IsContinue:           continueFlag,
			IsSkip:               false,
			IsUndo:               false,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				config := checkNewPullRequestPreconditions()
				return getNewPullRequestStepList(config)
			},
		})
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func checkNewPullRequestPreconditions() (result newPullRequestConfig) {
	if git.HasRemote("origin") {
		script.Fetch()
	}
	result.InitialBranch = git.GetCurrentBranchName()
	prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
	result.BranchesToSync = append(git.GetAncestorBranches(result.InitialBranch), result.InitialBranch)
	return
}

func getNewPullRequestStepList(config newPullRequestConfig) (result steps.StepList) {
	for _, branchName := range config.BranchesToSync {
		result.AppendList(steps.GetSyncBranchSteps(branchName))
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	result.Append(steps.CreatePullRequestStep{BranchName: config.InitialBranch})
	return
}

func init() {
	newPullRequestCommand.Flags().BoolVar(&abortFlag, "abort", false, abortFlagDescription)
	newPullRequestCommand.Flags().BoolVar(&continueFlag, "continue", false, continueFlagDescription)
	RootCmd.AddCommand(newPullRequestCommand)
}
