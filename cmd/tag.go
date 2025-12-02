package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var tagCmd = &cobra.Command{
	Use:   "tag <arc-name-or-id> <doc-id> <tag1> [tag2...]",
	Short: "Add tags to a document",
	Args:  cobra.MinimumNArgs(3),
	RunE:  runTag,
}

func init() {
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	arcNameOrID := args[0]
	docID := args[1]
	tags := args[2:]

	entry, err := arcManager.FindArc(arcNameOrID)
	if err != nil {
		return err
	}

	fmt.Printf("Adding tags to document in arc: %s\n", entry.Name)
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

	if err := arcManager.AddTags(entry.ID, arc, key, docID, tags); err != nil {
		return err
	}

	fmt.Printf("Tags added: %v\n", tags)
	return nil
}

