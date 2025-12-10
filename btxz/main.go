// File: main.go

// Package main implements the command-line interface for BTXZ.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"btxz/core"
	"btxz/update"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const version = "0.0.0‑dev" // <-- this will be auto‑replaced by CI

// main is the entry point for the application. It sets up the command structure
// and runs a background check for new updates.
func main() {
	// Run the update check in a separate goroutine so it doesn't block the UI.
	go update.CheckForUpdates(version)

	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

// NewRootCmd creates and configures the main 'btxz' command and its subcommands.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "btxz",
		Short:   "BTXZ: A professional, secure, and efficient file archiver.",
		Version: version,
		Long: `BTXZ is a professional command-line tool for creating and extracting
securely encrypted, highly compressed archives using a proprietary format.
Powered by XChaCha20-Poly1305 and LZMA2/XZ.`,
		// Suppress the default 'completion' command from cobra.
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Allow users to disable all styling for CI/CD or accessibility.
			if quiet, _ := cmd.Flags().GetBool("no-style"); quiet {
				pterm.DisableStyling()
				pterm.DisableColor()
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// After any command runs, display the update notification if one is available.
			update.DisplayUpdateNotification()
		},
	}

	rootCmd.SetVersionTemplate(`{{printf "btxz version %s\n" .Version}}`)
	rootCmd.Flags().Bool("no-style", false, "Disable all styling and colors")

	rootCmd.AddCommand(
		NewCreateCmd(),
		NewExtractCmd(),
		NewListCmd(),
		NewUpdateCmd(),
		NewTestCmd(),
	)

	return rootCmd
}

// NewCreateCmd configures the 'create' command.
func NewCreateCmd() *cobra.Command {
	var (
		outputFile string
		password   string
		level      string
	)
	createCmd := &cobra.Command{
		Use:   "create [file/folder...]",
		Short: "Create a new secure archive",
		Long: `Packages files into a secure .btxz archive using the V3 format (XChaCha20-Poly1305 + LZMA2).

ADAPTIVE PROFILES:
  --level low   : Low memory mode (64MB RAM, 1 pass). Good for Raspberry Pi/Mobile.
  --level default: Balanced mode (128MB RAM, 1 pass). Good for most laptops.
  --level max   : Paranoid mode (512MB RAM, 4 passes, Ultra Compression). High-end hardware only.`,
		Example: `  btxz create ./doc.pdf -o archive.btxz -p "pass" --level max`,
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			printCommandHeader("SECURE ARCHIVE CREATION")
			startTime := time.Now()

			if outputFile == "" {
				handleCmdError("Output file path must be specified with -o or --output.")
			}
			
			// Normalize level
			level = strings.ToLower(level)
			if level == "fast" { level = "low" }
			if level == "best" { level = "max" }

			if level != "low" && level != "default" && level != "max" {
				handleCmdError("Invalid level. Use: low, default, or max.")
			}
			
			promptForPassword(&password)

			pterm.DefaultSection.Println("Initialization")
			pterm.Info.Printf("Target: %s\n", outputFile)
			pterm.Info.Printf("Profile: %s\n", strings.ToUpper(level))
			pterm.Info.Println("Security: Enabled (XChaCha20-Poly1305)")

			pterm.DefaultSection.Println("Processing")
			spinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start(fmt.Sprintf("Compressing & Encrypting %d inputs...", len(args)))
			err := core.CreateArchive(outputFile, args, password, level)
			spinner.Stop()

			if err != nil {
				handleCmdError("Failed to create archive: %v", err)
			}
			
			duration := time.Since(startTime)

			// Show profile info
			var profileDesc string
			switch level {
			case "low":
				profileDesc = "Low-End / Fast"
			case "max":
				profileDesc = "Ultra / Hardened"
			default:
				profileDesc = "Balanced / Standard"
			}

			pterm.DefaultSection.Println("Mission Report")
			pterm.Success.Println("Operation Completed Successfully.")
			
			data := [][]string{
				{"Archive", outputFile},
				{"Security", "XChaCha20-Poly1305 (256-bit)"},
				{"Profile", profileDesc},
				{"Time Elapsed", duration.Round(time.Millisecond).String()},
				{"Status", "SECURED"},
			}
			
			pterm.DefaultTable.WithData(data).WithBoxed().Render()
		},
	}
	createCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Path for the new archive file (required)")
	createCmd.Flags().StringVarP(&password, "password", "p", "", "Password for encryption (prompts if empty, required)")
	createCmd.Flags().StringVarP(&level, "level", "l", "default", "Profile: low, default, max")

	return createCmd
}

// NewExtractCmd configures the 'extract' command.
func NewExtractCmd() *cobra.Command {
	var (
		outputDir string
		password  string
	)
	extractCmd := &cobra.Command{
		Use:     "extract <archive.btxz>",
		Short:   "Extract files from an archive",
		Long:    `Decompresses and decrypts a .btxz archive into the specified directory. Automatically detects v1, v2, and v3 formats.`,
		Example: `  btxz extract data.btxz -o ./restored_data`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			printCommandHeader("ARCHIVE EXTRACTION")
			startTime := time.Now()
			archivePath := args[0]
			
			if password == "" {
				pass, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter decryption password")
				password = pass
			}

			pterm.DefaultSection.Println("Processing")
			spinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start(fmt.Sprintf("Decrypting '%s'...", filepath.Base(archivePath)))
			skippedFiles, err := core.ExtractArchive(archivePath, outputDir, password)
			spinner.Stop()

			if err != nil {
				if strings.Contains(err.Error(), "decryption failed") || strings.Contains(err.Error(), "authentication failed") {
					handleCmdError("Access Denied: Incorrect Password or Corrupted Archive.")
				}
				handleCmdError("Critical Error: %v", err)
			}

			duration := time.Since(startTime)
			pterm.DefaultSection.Println("Mission Report")

			if len(skippedFiles) > 0 {
				pterm.Warning.Println("Operation Completed with Warnings.")
				pterm.DefaultBox.WithTitle("Skipped Files (Safe Mode)").WithBoxStyle(pterm.NewStyle(pterm.FgYellow)).Println(
					strings.Join(skippedFiles, "\n"),
				)
			} else {
				pterm.Success.Println("All files extracted successfully.")
			}

			data := [][]string{
				{"Source", filepath.Base(archivePath)},
				{"Destination", outputDir},
				{"Time Elapsed", duration.Round(time.Millisecond).String()},
				{"Status", "RESTORED"},
			}
			pterm.DefaultTable.WithData(data).WithBoxed().Render()
		},
	}
	extractCmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "Directory to extract files to")
	extractCmd.Flags().StringVarP(&password, "password", "p", "", "Password for decryption (prompts if empty)")
	return extractCmd
}

// NewTestCmd configures the 'test' command.
func NewTestCmd() *cobra.Command {
	var password string
	testCmd := &cobra.Command{
		Use:     "test <archive.btxz>",
		Short:   "Test integrity of an archive",
		Long:    `Verifies the integrity of a .btxz archive (V3+) by decrypting and decompressing the stream without writing to disk.`,
		Example: `  btxz test backup.btxz -p "s3cr3t!"`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			printCommandHeader("INTEGRITY VERIFICATION")
			startTime := time.Now()
			archivePath := args[0]

			if password == "" {
				pass, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter decryption password")
				password = pass
			}

			pterm.DefaultSection.Println("Analysis")
			spinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Verifying structure and checksums...")
			err := core.TestArchive(archivePath, password)
			spinner.Stop()

			if err != nil {
				pterm.Error.Println("INTEGRITY CHECK FAILED")
				pterm.Error.Println(err.Error())
				os.Exit(1)
			}

			duration := time.Since(startTime)
			pterm.DefaultSection.Println("Mission Report")
			pterm.Success.Println("Verification Passed.")
			
			data := [][]string{
				{"Target", filepath.Base(archivePath)},
				{"Integrity", "VALID"},
				{"Time Elapsed", duration.Round(time.Millisecond).String()},
				{"Status", "VERIFIED"},
			}
			pterm.DefaultTable.WithData(data).WithBoxed().Render()
		},
	}
	testCmd.Flags().StringVarP(&password, "password", "p", "", "Password for decryption (prompts if empty)")
	return testCmd
}

// NewListCmd configures the 'list' command.
func NewListCmd() *cobra.Command {
	var password string
	listCmd := &cobra.Command{
		Use:     "list <archive.btxz>",
		Short:   "List the contents of an archive",
		Long:    `Shows a list of files and folders inside a .btxz archive without extracting them. Automatically handles all versions.`,
		Example: `  btxz list my_archive.btxz -p "s3cr3t!"`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			printCommandHeader("ARCHIVE CONTENTS")
			archivePath := args[0]
			
			if password == "" {
				pass, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter decryption password")
				password = pass
			}

			spinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Decrypting metadata...")
			contents, err := core.ListArchiveContents(archivePath, password)
			spinner.Stop()

			if err != nil {
				if strings.Contains(err.Error(), "decryption failed") || strings.Contains(err.Error(), "authentication failed") {
					handleCmdError("Access Denied: Incorrect Password.")
				}
				handleCmdError("Failed to list archive contents: %v", err)
			}

			pterm.Success.Printf("Index retrieved for %s.\n", filepath.Base(archivePath))
			tableData := pterm.TableData{{"Mode", "Size (bytes)", "Name"}}
			for _, item := range contents {
				tableData = append(tableData, []string{item.Mode, fmt.Sprintf("%d", item.Size), item.Name})
			}
			pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).Render()
		},
	}
	listCmd.Flags().StringVarP(&password, "password", "p", "", "Password for decryption (prompts if empty)")
	return listCmd
}

// NewUpdateCmd configures the 'update' command.
func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update btxz to the latest version",
		Long:  `Checks for the latest version on GitHub and performs an in-place update if available.`,
		Run: func(cmd *cobra.Command, args []string) {
			printCommandHeader("SYSTEM UPDATE")
			if err := update.PerformUpdate(version); err != nil {
				handleCmdError("Update failed: %v", err)
			}
			pterm.Success.Println("BTXZ has been updated successfully!")
		},
	}
}

// --- Helper Functions ---

// handleCmdError prints a formatted error message and exits the application.
func handleCmdError(format string, a ...interface{}) {
	pterm.Error.Printf(format+"\n", a...)
	os.Exit(1)
}

// promptForPassword checks if a password string is empty and, if so, prompts
// the user for it.
func promptForPassword(password *string) {
	if *password == "" {
		pterm.Info.Println("No password provided via flags.")
		pass, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Set encryption password")
		*password = pass
	}
	if *password == "" {
		handleCmdError("Aborted: A password is required to encrypt the archive.")
	}
}

// printCommandHeader displays the standard logo and title for a command.
func printCommandHeader(title string) {
	// Clear screen for a fresh look
	print("\033[H\033[2J")
	// Cyberpunk/Matrix style gradient logo
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("BT", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("XZ", pterm.NewStyle(pterm.FgLightMagenta)),
	).Render()
	
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgBlack)).WithTextStyle(pterm.NewStyle(pterm.FgCyan, pterm.Bold)).Println(title)
	fmt.Println()
}
