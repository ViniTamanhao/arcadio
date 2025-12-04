package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <arc-name-or-id>",
	Short: "Show arc information",
	Long:  "Display detailed information about an arc (requires password).",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	arcNameOrID := args[0]

	entry, err := arcManager.FindArc(arcNameOrID)
	if err != nil {
		return fmt.Errorf("failed to find arc: %w", err)
	}



	password, err := authManager.GetPassword(entry.Name, entry.ID, true)
	if err != nil {
		return err
	}

	arc, _, err := arcManager.Unlock(entry.ID, password)
	if err != nil {
		return fmt.Errorf("failed to unlock arc: %w", err)
	}

	fmt.Printf("\nArc Information:\n")
	fmt.Printf("================\n")
	fmt.Printf("Name:         %s\n", arc.Name)
	fmt.Printf("ID:           %s\n", arc.ID)
	fmt.Printf("Created:      %s\n", arc.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Modified:     %s\n", arc.ModifiedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Documents:    %d\n", len(arc.Documents))
	fmt.Printf("Tags:         %d unique\n", countUniqueTags(arc.Tags))
	fmt.Printf("Encryption:   %s\n", arc.EncryptionVersion)

	var totalSize int64
	for _, doc := range arc.Documents {
		totalSize += doc.Size
	}
	fmt.Printf("Total Size:   %s\n", formatSize(totalSize))

	return nil
}

func countUniqueTags(tagMap map[string][]string) int {
	uniqueTags := make(map[string]bool)
	for _, tags := range tagMap {
		for _, tag := range tags {
			uniqueTags[tag] = true
		}
	}
	return len(uniqueTags)
}

