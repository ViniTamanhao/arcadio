// Package cmd defines all CLI commands for arcadio
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ViniTamanhao/arcadio/internal/arc"
	"github.com/ViniTamanhao/arcadio/internal/auth"
	"github.com/spf13/cobra"
)

var (
	arcManager  *arc.Manager
	authManager *auth.Manager
	baseDir     string
)

var rootCmd = &cobra.Command{
	Use:   "arc",
	Short: "arcadio - Secure encrypted document archives",
	Long: `arcadio (arc) is a CLI tool for creating and managing encrypted document archives.
Each arc is a secure arc that protects your documents with strong encryption.`,
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
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		baseDir = filepath.Join(home, ".arcadio", "arcs")
	}

	// Create base directory
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating base directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize managers
	var err error
	arcManager, err = arc.NewManager(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing arc manager: %v\n", err)
		os.Exit(1)
	}

	// Initialize auth manager
	authManager = auth.NewManager(filepath.Join(baseDir, ".."))
}
