package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var removeCmd = &cobra.Command{
	Use:   "remove <arc-id> <doc-id>",
	Short: "Remove a document from an arc",
	Aliases: []string{"rm"},
	Args:  cobra.ExactArgs(2),
	RunE:  runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	arcID := args[0]
	docID := args[1]

	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	password := string(passwordBytes)

	arc, key, err := arcManager.Unlock(arcID, password)
	if err != nil {
		return err
	}

	if err := arcManager.RemoveDocument(arcID, arc, key, docID); err != nil {
		return err
	}

	return nil
}
