package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var keyringCmd = &cobra.Command{
	Use:   "keyring",
	Short: "Manage stored passwords",
	Long:  "Manage passwords stored in the system keyring",
}

var keyringListCmd = &cobra.Command{
	Use:   "list",
	Short: "List arcs with stored passwords",
	RunE:  runKeyringList,
}

var keyringDeleteCmd = &cobra.Command{
	Use:   "delete <arc-name-or-id>",
	Short: "Delete a stored password",
	Args:  cobra.ExactArgs(1),
	RunE:  runKeyringDelete,
}

var keyringSaveCmd = &cobra.Command{
	Use:   "save <arc-name-or-id>",
	Short: "Save password to keyring",
	Args:  cobra.ExactArgs(1),
	RunE:  runKeyringSave,
}

var keyringClearCmd = &cobra.Command{
	Use:   "clear-session",
	Short: "Clear session password cache",
	RunE:  runKeyringClear,
}

func init() {
	rootCmd.AddCommand(keyringCmd)
	keyringCmd.AddCommand(keyringListCmd)
	keyringCmd.AddCommand(keyringDeleteCmd)
	keyringCmd.AddCommand(keyringSaveCmd)
	keyringCmd.AddCommand(keyringClearCmd)
}

func runKeyringList(cmd *cobra.Command, args []string) error {
	arcs := arcManager.ListArcs()
	
	fmt.Println("Arcs with stored passwords:")
	fmt.Println()
	
	count := 0
	for _, arc := range arcs {
		if authManager.HasStoredPassword(arc.ID) {
			fmt.Printf("  ✓ %s (%s)\n", arc.Name, arc.ID[:8]+"...")
			count++
		}
	}
	
	if count == 0 {
		fmt.Println("  No passwords stored in keyring")
	}
	
	return nil
}

func runKeyringDelete(cmd *cobra.Command, args []string) error {
	arcNameOrID := args[0]
	
	entry, err := arcManager.FindArc(arcNameOrID)
	if err != nil {
		return err
	}
	
	if err := authManager.DeletePassword(entry.ID); err != nil {
		return err
	}
	
	fmt.Printf("✓ Password deleted from keyring: %s\n", entry.Name)
	return nil
}

func runKeyringSave(cmd *cobra.Command, args []string) error {
	arcNameOrID := args[0]
	
	entry, err := arcManager.FindArc(arcNameOrID)
	if err != nil {
		return err
	}
	
	// Prompt for password
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()
	
	password := string(passwordBytes)
	
	// Verify password by trying to unlock
	_, _, err = arcManager.Unlock(entry.ID, password)
	if err != nil {
		return fmt.Errorf("invalid password")
	}
	
	// Save to keyring
	if err := authManager.SavePassword(entry.ID, password); err != nil {
		return err
	}
	
	fmt.Printf("✓ Password saved to keyring: %s\n", entry.Name)
	return nil
}

func runKeyringClear(cmd *cobra.Command, args []string) error {
	authManager.ClearSession()
	fmt.Println("✓ Session cache cleared")
	return nil
}
