package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func runList(cmd *cobra.Command, args []string) error {
	arcs := arcManager.ListArcs()

	if len(arcs) == 0 {
		fmt.Println("No arcs found. Create one with: arc create <name>")
		return nil
	}

	sort.Slice(arcs, func(i, j int) bool {
		return arcs[i].CreatedAt.After(arcs[j].CreatedAt)
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tID\tCREATED")
	fmt.Fprintln(w, "----\t--\t-------")

	for _, arc := range arcs {
		fmt.Fprintf(w, "%s\t%s\t%s\n", 
			arc.Name,
			arc.ID[:8]+"...",
			arc.CreatedAt.Format("2006-01-02 15:04"),
		)
	}

	w.Flush()
	return nil
}
