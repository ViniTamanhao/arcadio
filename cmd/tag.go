package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var tagCmd = &cobra.Command{
	Use:   "tag <arc-id> <doc-id> <tag1> [tag2...]",
	Short: "Add tags to a document",
	Args:  cobra.MinimumNArgs(3),
	RunE:  runTag,
}

func init() {
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	arcID := args[0]
	docID := args[1]
	tags := args[2:]

	// Get password
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	password := string(passwordBytes)

	// Unlock arc
	arc, key, err := arcManager.Unlock(arcID, password)
	if err != nil {
		return err
	}

	// Add tags
	if err := arcManager.AddTags(arcID, arc, key, docID, tags); err != nil {
		return err
	}

	fmt.Printf("Tags added: %v\n", tags)
	return nil
}

