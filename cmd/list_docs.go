package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var listDocsCmd = &cobra.Command{
	Use:   "docs <arc-name-or-id>",
	Short: "List documents in an arc",
	Args:  cobra.ExactArgs(1),
	RunE:  runListDocs,
}

func init() {
	rootCmd.AddCommand(listDocsCmd)
}

func runListDocs(cmd *cobra.Command, args []string) error {
	arcNameOrID := args[0]

	entry, err := arcManager.FindArc(arcNameOrID)
	if err != nil {
		return err
	}

	fmt.Printf("Getting docs for arc: %s\n", entry.Name)
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	password := string(passwordBytes)

	arc, _, err := arcManager.Unlock(arcNameOrID, password)
	if err != nil {
		return err
	}

	docs := arcManager.ListDocuments(arc)

	if len(docs) == 0 {
		fmt.Println("No documents in this arc.")
		return nil
	}

	fmt.Printf("\nArc : %s\n", arc.Name)
	fmt.Printf("Documents: %d\n\n", len(docs))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tFILENAME\tSIZE\tADDED\tTAGS")
	fmt.Fprintln(w, "--\t--------\t----\t-----\t----")

	for _, doc := range docs {
		tags := arc.Tags[doc.ID]
		tagStr := ""
		if len(tags) > 0 {
			tagStr = fmt.Sprintf("[%s]", strings.Join(tags, ", "))
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			doc.ID,
			doc.Filename,
			formatSize(doc.Size),
			formatTime(doc.AddedAt),
			tagStr,
		)
	}

	w.Flush()
	return nil
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 24*time.Hour {
		return "today"
	} else if diff < 48*time.Hour {
		return "yesterday"
	} else if diff < 7*24*time.Hour {
		return fmt.Sprintf("%dd ago", int(diff.Hours()/24))
	}
	return t.Format("2006-01-02")
}
