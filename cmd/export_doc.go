package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var exportDocCmd = &cobra.Command{
	Use: "export <arc-name-or-id> <doc-id> <output-path>",
	Short: "Export a document from an arc",
	Args: cobra.ExactArgs(3),
	RunE: runExportDoc,
}

func init() {
	rootCmd.AddCommand(exportDocCmd)
}

func runExportDoc(cmd *cobra.Command, args []string) error {
	arcNameOrID := args[0]
	docID := args[1]
	outputPath := args[2]

	entry, err := arcManager.FindArc(arcNameOrID)	
	if err != nil {
		return err
	}

	fmt.Printf("Exporting arc: %s\n", entry.Name)

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

	if err := arcManager.ExportDocument(entry.ID, arc, key, docID, outputPath); err != nil {
		return err
	}

	return nil
}
