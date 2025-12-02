// Package cmd defines all CLI commands for arcadio
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ViniTamanhao/arcadio/internal/arc"
	"github.com/spf13/cobra"
)

var (
	arcManager 	*arc.Manager
	baseDir 		string
)

var rootCmd = &cobra.Command{
	Use: "arc",
	Short: "arcadio - Secure encrypted document archives",
	Long: `arcadio (arc) is a CLI tool for creating and managing encrypted document archives.
	Each arc is a secure vault that protects your documents with string encryption.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&baseDir, "base-dir", "", "Base directory for arcs (default: ~/.arcadio)")
}

func initConfig() {
	if baseDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting home directory: %v\n", err)
			os.Exit(1)
		}
		baseDir = filepath.Join(home, ".arcadio", "arcs")
	}

	if err := os.MkdirAll(baseDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "error creating base directory: %v\n", err)
		os.Exit(1)
	}

	var err error
	arcManager, err = arc.NewManager(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initiating arc manager: %v\n", err)
		os.Exit(1)
	}
}
