package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var createCmd = &cobra.Command{
	Use: "create <name>",
	Short: "Create a new encrypted arc",
	Long: `Create a new encrypted arc with a password and security question.
	The arc will be encrypted using AES-256-GCM with Argon2id key derivation.`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	fmt.Printf("Creating arc: %s\n\n", name)

	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	password := string(passwordBytes)
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 9 characters")
	}

	fmt.Print("Confirm password: ")
	confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	if password != string(confirmBytes) {
		return fmt.Errorf("passwords do not match")
	}

	var securityQuestion string
	fmt.Print("Security question: ")
	fmt.Scanln(&securityQuestion)

	if securityQuestion == "" {
		return fmt.Errorf("security question cannot be empty")
	}
	fmt.Print("Answer: ")
	answerBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read answer: %w", err)
	}
	fmt.Println()

	answer := string(answerBytes)
	if answer == "" {
		return fmt.Errorf("security answer cannot be empty")
	}

	arc, err := arcManager.Create(name, password, securityQuestion, answer)
	if err != nil {
		return fmt.Errorf("failed to create arc: %w", err)
	}

	fmt.Printf("\nArc created successfully!\n")
	fmt.Printf("	ID: %s\n", arc.ID)
	fmt.Printf("	Name: %s\n", arc.Name)
	fmt.Printf("	Created: %s\n", arc.CreatedAt.Format("2006-01-02 15:04:05"))

	return nil
}
