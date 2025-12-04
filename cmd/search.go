package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <arc-name-or-id> <query>",
	Short: "Search for documents in an arc",
	Args:  cobra.ExactArgs(2),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	arcNameOrID := args [0]
	query := args[1]

	entry, err := arcManager.FindArc(arcNameOrID)
	if err != nil {
		return err
	}

	fmt.Printf("Searching arc: %s\n", entry.Name)

	password, err := authManager.GetPassword(entry.ID, entry.Name, true)
	if err != nil {
		return err
	}

	arc, _, err := arcManager.Unlock(entry.ID, password)
	if err != nil {
		return err
	}

	results := arcManager.SearchDocuments(arc, query)
	if len(results) == 0 {
		fmt.Printf("No documents found matching: %s\n", query)
		return nil
	}

	fmt.Printf("\nFound %d document(s) matching: %s\n\n", len(results), query)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tFILENAME\tSIZE\tTAGS")
	fmt.Fprintln(w, "--\t--------\t----\t----")

	for _, doc := range results {
		tags := arc.Tags[doc.ID]
		tagStr := ""
		if len(tags) > 0 {
			tagStr = fmt.Sprintf("[%s]", strings.Join(tags, ", "))
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n",
			doc.ID[:8]+"...",
			doc.Filename,
			formatSize(doc.Size),
			tagStr,
		)
	}

	w.Flush()
	return nil
}
