// File: update.go
// This package handles the application's self-updating mechanism.

package update

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"bytes"

	"github.com/inconshreveable/go-update"
	"github.com/pterm/pterm"
)

const versionURL = "https://raw.githubusercontent.com/BlackTechX011/BTXZ/main/version.json"

// updateArt is the visual warning for an available update.
const updateArt = `
 ╔══════════════════════════════════╗
 ║      SYSTEM UPDATE DETECTED      ║
 ╚══════════════════════════════════╝
      ↓↓ INSTALLING PATCHES ↓↓`

// ReleaseInfo defines the structure of the version.json file on GitHub.
type ReleaseInfo struct {
	Version   string                     `json:"version"`
	Notes     string                     `json:"notes"`
	Platforms map[string]PlatformDetails `json:"platforms"`
}

// PlatformDetails contains the download URL and security checksum for a specific architecture.
type PlatformDetails struct {
	URL      string `json:"url"`
	Checksum string `json:"sha256"` // SHA256 hash of the binary
}

// Cache for the latest release info.
var (
	latestRelease *ReleaseInfo
	checkOnce     sync.Once
	mu            sync.RWMutex
)

// CheckForUpdates fetches release information from GitHub.
// It is designed to be run in a goroutine and will not block.
// It handles network errors gracefully by simply doing nothing.
func CheckForUpdates(currentVersion string) {
	checkOnce.Do(func() {
		resp, err := http.Get(versionURL)
		if err != nil {
			return // Fail silently on network errors
		}
		defer resp.Body.Close()

		var release ReleaseInfo
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return // Fail silently on JSON parsing errors
		}

		// Simple string comparison works for "v1.0" > "v1.0"
		if release.Version > currentVersion {
			mu.Lock()
			latestRelease = &release
			mu.Unlock()
		}
	})
}

// DisplayUpdateNotification prints a prominent warning if a new version is available.
func DisplayUpdateNotification() {
	mu.RLock()
	release := latestRelease
	mu.RUnlock()

	if release != nil {
		pterm.Println() // Add some space
		message := fmt.Sprintf("A new version (%s) is available!\n\nNotes: %s\n\n%s",
			pterm.LightGreen(release.Version),
			release.Notes,
			pterm.LightYellow("Run 'btxz update' to get the latest features and security fixes."),
		)
		pterm.DefaultBox.WithTitle(pterm.LightYellow("UPDATE AVAILABLE")).WithTitleTopCenter().WithBoxStyle(pterm.NewStyle(pterm.FgYellow)).Println(pterm.FgYellow.Sprint(updateArt) + "\n\n" + message)
	}
}

// PerformUpdate executes the self-update process.
func PerformUpdate(currentVersion string) error {
	// Force a fresh check
	checkOnce = sync.Once{} 
	CheckForUpdates(currentVersion)

	mu.RLock()
	release := latestRelease
	mu.RUnlock()

	if release == nil {
		pterm.Info.Println("Your system is up to date.")
		return nil
	}

	platformKey := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	platformInfo, ok := release.Platforms[platformKey]
	if !ok {
		return fmt.Errorf("no update available for your platform: %s", platformKey)
	}

	pterm.DefaultSection.Println("Update Found")
	pterm.Info.Printf("Current: %s\n", currentVersion)
	pterm.Info.Printf("Latest:  %s\n", pterm.Green(release.Version))
	pterm.Info.Printf("Notes:   %s\n", release.Notes)

	// --- DOWNLOAD PHASE ---
	pterm.DefaultSection.Println("Downloading")
	
	req, err := http.NewRequest("GET", platformInfo.URL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.ContentLength > 0 {
		bar, _ := pterm.DefaultProgressbar.WithTotal(int(resp.ContentLength)).WithTitle("Downloading update...").Start()
		// Wrap body to update the progress bar
		proxyReader := &progressReader{
			Reader: resp.Body,
			Bar:    bar,
		}
		// Read into memory
		data, err := io.ReadAll(proxyReader)
		if err != nil {
			return fmt.Errorf("download interrupted: %w", err)
		}
		bar.Stop() // Ensure bar finishes
		
		// --- VERIFICATION PHASE ---
		pterm.DefaultSection.Println("Security Checks")
		
		if platformInfo.Checksum != "" {
			spinner, _ := pterm.DefaultSpinner.Start("Verifying SHA256 checksum...")
			hash := sha256.Sum256(data)
			calculatedHash := hex.EncodeToString(hash[:])
			if calculatedHash != platformInfo.Checksum {
				spinner.Fail("Checksum Mismatch!")
				return fmt.Errorf("security check failed: expected %s, got %s", platformInfo.Checksum, calculatedHash)
			}
			spinner.Success("Checksum Verified")
		} else {
			pterm.Warning.Println("Skipping checksum checks (not provided in manifest).")
		}

		// --- INSTALLATION PHASE ---
		pterm.DefaultSection.Println("Installation")
		pterm.Info.Println("Replacing binary...")
		
		reader := bytes.NewReader(data)
		err = update.Apply(reader, update.Options{})
		if err != nil {
			if rerr := update.RollbackError(err); rerr != nil {
				return fmt.Errorf("failed to apply update and rollback failed: %v", rerr)
			}
			return fmt.Errorf("failed to apply update: %w", err)
		}
		
		// --- SUMMARY ---
		pterm.DefaultSection.Println("Mission Report")
		reportData := [][]string{
			{"Previous Version", currentVersion},
			{"New Version", pterm.Green(release.Version)},
			{"Platform", platformKey},
			{"Status", "UPDATED"},
		}
		pterm.DefaultTable.WithData(reportData).WithBoxed().Render()
		pterm.Success.Println("BTXZ has been updated successfully. Please restart your terminal.")

		return nil
	}
	
	return fmt.Errorf("server returned unknown content length")
}

// progressReader wraps an io.Reader to update a pterm.Progressbar
type progressReader struct {
	io.Reader
	Bar *pterm.ProgressbarPrinter
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	if n > 0 && pr.Bar != nil {
		pr.Bar.Add(n)
	}
	return
}
