package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var removeCmd = &cobra.Command{
	Use:   "remove <arc-name-or-id> <doc-id>",
	Short: "Remove a document from an arc",
	Aliases: []string{"rm"},
	Args:  cobra.ExactArgs(2),
	RunE:  runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	arcNameOrID := args[0]
	docID := args[1]

	entry, err := arcManager.FindArc(arcNameOrID)
	if err != nil {
		return err
	}

	fmt.Printf("Removing from arc: %s\n", entry.Name)
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	password := string(passwordBytes)

	arc, key, err := arcManager.Unlock(entry.ID, password)
	if err != nil {
		return err
	}

	if err := arcManager.RemoveDocument(entry.ID, arc, key, docID); err != nil {
		return err
	}

	return nil
}
