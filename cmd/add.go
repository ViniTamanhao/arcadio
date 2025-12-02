package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/ViniTamanhao/arcadio/pkg/models"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	addTags []string
	addRecursive bool
)

var addCmd = &cobra.Command{
	Use: "add <arc-id> <file-or-directory>",
	Short: "Add documents to an arc",
	Long: `Add one or more documents to an encrypted arc. Documents will be encrypted before storage.`,
	Args: cobra.ExactArgs(2),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringSliceVarP(&addTags, "tags", "t", []string{}, "Tags to add to the document(s)")
	addCmd.Flags().BoolVarP(&addRecursive, "recursive", "r", false, "Add directory recursively")
}

func runAdd(cmd *cobra.Command, args []string) error {
	arcID := args[0]
	path := args [1]

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

	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to access path: %w", err)
	}

	if fileInfo.IsDir() {
		if !addRecursive {
			return fmt.Errorf("path is a directory, use --recursive flag to add all files")
		}
		return addDirectory(arcID, arc, key, path)
	}

	doc, err := arcManager.AddDocument(arcID, arc, key, path, addTags)
	if err != nil {
		return err
	}

	fmt.Printf("\nDocument added: %s\n", doc.Filename)
	fmt.Printf("	ID: %s\n", doc.ID)
	fmt.Printf("	Size: %d bytes\n", doc.Size)
	if len(addTags) > 0 {
		fmt.Printf("	Tags: %v\n", addTags)
	}

	return nil
}

func addDirectory(arcID string, arc *models.Arc, key []byte, dirPath string) error {
	count := 0
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fmt.Printf("\nAdding: %s\n", path)
		_, err = arcManager.AddDocument(arcID, arc, key, path, addTags)
		if err != nil {
			fmt.Printf("Failed: %v\n", err)
			return nil
		}
		count++
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("\nAdded %d documents\n", count)
	return nil
}
