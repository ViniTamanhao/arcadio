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

func runDelete(cmd *cobra.Command, args []string) error {
	arcNameOrId := args[0]

	entry, err := arcManager.FindArc(arcNameOrId)
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
