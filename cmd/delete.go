package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	forceDelete bool
)

var deleteCmd = &cobra.Command{
	Use: "delete <arc-name-or-id>",
	Short: "Delete an arc permanently",
	Long: "Delete an arc and all its documents permanently. This cannot be undone!",
	Aliases: []string{"rm", "del"},
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Skip confirmation prompt")
}

func runDelete(cmd *cobra.Command, args []string) error {
	arcNameOrID := args[0]

	entry, err := arcManager.FindArc(arcNameOrID)
	if err != nil {
		return err
	}

	password, err := authManager.GetPassword(entry.ID, entry.Name, true)
	if err != nil {
		return err
	}

	_, _, err = arcManager.Unlock(entry.ID, password)
	if err != nil {
		return err
	}

	if !forceDelete {
		fmt.Printf("WARNING: This will permantently delete arc: '%s' and all its documents!\n", entry.Name)
		fmt.Print("Type arc name to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		confirmation, _ := reader.ReadString('\n')
		confirmation = strings.TrimSpace(confirmation)

		if confirmation != entry.Name {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	if err := arcManager.Delete(entry.ID); err != nil {
		return fmt.Errorf("failed to delete arc: %w", err)
	}

	fmt.Printf("Arc '%s' deleted successfully\n", entry.Name)
	return nil
}
