// File: core/v3.go

// Package core contains the stable, versioned logic for the BTXZ archive format.
// This file implements the v3 specification (The "Pro" Version).
// Improvements:
// - Security: Switched from AES-256-GCM (12-byte nonce) to XChaCha20-Poly1305 (24-byte nonce).
//   This eliminates the risk of nonce collision with random nonces.
// - Compression: Switched back to XZ (LZMA2) for maximum compression ratio, but with optimized presets.
// Core Version: v3
package core

import (
	"archive/tar"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

// --- v3 Core Constants & Header Definition ---

const (
	// coreVersionV3 is the integer identifier for this version of the format.
	coreVersionV3 = 3

	// XChaCha20-Poly1305 Constants
	xNonceSize = 24 // XChaCha20 uses a 24-byte nonce (192 bits)
	xKeyLength = 32 // 32 bytes = 256 bits

	// Helper sizes
	v3HeaderSize = 4 + 2 + 1 + saltSize + 4 + 4 + 1 + xNonceSize
)

// BtxzHeaderV3 defines the binary structure of the v3 archive header.
// It uses XChaCha20-Poly1305 for superior security.
type BtxzHeaderV3 struct {
	Signature        [4]byte // "BTXZ"
	Version          uint16  // 3
	CompressionLevel uint8   // 1=Fast, 2=Default, 3=Best
	Salt             [saltSize]byte
	Argon2Time       uint32
	Argon2Memory     uint32
	Argon2Threads    uint8
	Nonce            [xNonceSize]byte // 24 bytes for XChaCha20
}

// CreateArchiveV3 creates a new archive using the v3 format (Tar -> XZ -> XChaCha20-Poly1305).
// It now supports adaptive profiles for hardware optimization.
func CreateArchiveV3(archivePath string, inputPaths []string, password string, level string) error {
	if len(inputPaths) == 0 {
		return errors.New("no input files or folders specified")
	}
	if password == "" {
		return errors.New("a password is required for v3 archives")
	}

	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("could not create archive file: %w", err)
	}
	defer archiveFile.Close()

	// 1. Configure Header and Crypto Params based on Profile
	header := BtxzHeaderV3{
		Signature:     [4]byte{'B', 'T', 'X', 'Z'},
		Version:       coreVersionV3,
		Argon2Threads: argon2Threads,
	}

	// Adaptive Profiles Configuration
	var xzDictCap int
	
	switch level {
	case "fast", "low": // Low-End Hardware Mode
		header.CompressionLevel = levelFast
		header.Argon2Memory = 64 * 1024       // 64 MB (Good for Pi/Mobile)
		header.Argon2Time = 1                 // 1 Pass
		xzDictCap = 1 * 1024 * 1024           // 1 MiB Dictionary (Very low memory usage)
	case "best", "max": // Max Security & Compression Mode
		header.CompressionLevel = levelBest
		header.Argon2Memory = 512 * 1024      // 512 MB (High Security)
		header.Argon2Time = 4                 // 4 Passes
		xzDictCap = 64 * 1024 * 1024          // 64 MiB Dictionary (Better compression, higher memory)
	default: // Default / Balanced Mode
		header.CompressionLevel = levelDefault
		header.Argon2Memory = 128 * 1024      // 128 MB Standard
		header.Argon2Time = 1                 // 1 Pass
		xzDictCap = 8 * 1024 * 1024           // 8 MiB Dictionary
	}

	// Generate Salt and Nonce
	if _, err := rand.Read(header.Salt[:]); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}
	if _, err := rand.Read(header.Nonce[:]); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Derive Key
	key := argon2.IDKey([]byte(password), header.Salt[:], header.Argon2Time, header.Argon2Memory, header.Argon2Threads, xKeyLength)

	// 2. Prepare Tar and XZ Writers
	compressedBuffer := new(bytes.Buffer)
	
	// Configure XZ Writer with Profile Settings
	// Using a larger dictionary improves compression but requires more memory for both compression and decompression.
	xzConfig := xz.WriterConfig{
		DictCap: xzDictCap,
	}
	xzWriter, err := xzConfig.NewWriter(compressedBuffer)
	if err != nil {
		return fmt.Errorf("failed to create xz writer: %w", err)
	}
	
	tarWriter := tar.NewWriter(xzWriter)

	// 3. Add files to Tar
	for _, path := range inputPaths {
		basePath := filepath.Dir(path)
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("could not stat input path %s: %w", path, err)
		}
		if info.IsDir() {
			basePath = path
		}
		
		walkErr := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			return addFileToTar(tarWriter, filePath, basePath)
		})
		if walkErr != nil {
			tarWriter.Close()
			xzWriter.Close()
			return fmt.Errorf("failed while walking path %s: %w", path, walkErr)
		}
	}
	
	if err := tarWriter.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := xzWriter.Close(); err != nil {
		return fmt.Errorf("failed to close xz writer: %w", err)
	}

	// 4. Write Header
	if err := binary.Write(archiveFile, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("failed to write archive header: %w", err)
	}

	// 5. Encrypt with XChaCha20-Poly1305
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return fmt.Errorf("failed to create XChaCha20-Poly1305 AEAD: %w", err)
	}

	// Seal appends to the first argument (dst). We pass nil to allocate new slice.
	encryptedPayload := aead.Seal(nil, header.Nonce[:], compressedBuffer.Bytes(), nil)

	if _, err := archiveFile.Write(encryptedPayload); err != nil {
		return fmt.Errorf("failed to write encrypted payload: %w", err)
	}

	return nil
}

// getDecryptedReaderV3 opens a v3 archive, handles XChaCha20 decryption.
func getDecryptedReaderV3(archivePath string, password string) (io.Reader, error) {
	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer archiveFile.Close()

	var header BtxzHeaderV3
	if err := binary.Read(archiveFile, binary.LittleEndian, &header); err != nil {
		return nil, fmt.Errorf("failed to read v3 archive header: %w", err)
	}

	key := argon2.IDKey([]byte(password), header.Salt[:], header.Argon2Time, header.Argon2Memory, header.Argon2Threads, xKeyLength)

	// Read Encrypted Payload
	encryptedPayload, err := io.ReadAll(archiveFile)
	if err != nil {
		return nil, fmt.Errorf("could not read encrypted payload: %w", err)
	}

	// Decrypt
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create XChaCha20-Poly1305 AEAD: %w", err)
	}

	decryptedPayload, err := aead.Open(nil, header.Nonce[:], encryptedPayload, nil)
	if err != nil {
		return nil, errors.New("decryption failed: incorrect password or tampered archive")
	}

	return bytes.NewReader(decryptedPayload), nil
}

// ExtractArchiveV3 extracts a v3 archive.
func ExtractArchiveV3(archivePath, outputDir, password string) ([]string, error) {
	var skippedFiles []string
	
	payloadReader, err := getDecryptedReaderV3(archivePath, password)
	if err != nil {
		return nil, err
	}

	// Decompress XZ
	xzReader, err := xz.NewReader(payloadReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create xz reader: %w", err)
	}
	
	tarReader := tar.NewReader(xzReader)
	cleanOutputDir, _ := filepath.Abs(filepath.Clean(outputDir))

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return skippedFiles, fmt.Errorf("error reading tar stream: %w", err)
		}

		targetPath := filepath.Join(cleanOutputDir, hdr.Name)
		cleanTargetPath := filepath.Clean(targetPath)

		if !strings.HasPrefix(cleanTargetPath, cleanOutputDir) {
			skippedFiles = append(skippedFiles, hdr.Name)
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(targetPath, os.FileMode(hdr.Mode))
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(targetPath), 0755)
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return skippedFiles, err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return skippedFiles, err
			}
			outFile.Close()
		}
	}
	return skippedFiles, nil
}

// TestArchiveV3 verifies the integrity of a v3 archive.
func TestArchiveV3(archivePath, password string) error {
	payloadReader, err := getDecryptedReaderV3(archivePath, password)
	if err != nil {
		return err
	}

	xzReader, err := xz.NewReader(payloadReader)
	if err != nil {
		return fmt.Errorf("integrity check failed: invalid compressed data: %w", err)
	}

	// Read and discard output to verify stream integrity
	if _, err := io.Copy(io.Discard, xzReader); err != nil {
		return fmt.Errorf("integrity check failed: data corruption detected: %w", err)
	}

	return nil
}

// ListArchiveContentsV3 lists contents of a v3 archive.
func ListArchiveContentsV3(archivePath, password string) ([]ArchiveEntry, error) {
	payloadReader, err := getDecryptedReaderV3(archivePath, password)
	if err != nil {
		return nil, err
	}

	xzReader, err := xz.NewReader(payloadReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create xz reader: %w", err)
	}
	
	tarReader := tar.NewReader(xzReader)
	var contents []ArchiveEntry

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		entry := ArchiveEntry{
			Mode: os.FileMode(hdr.Mode).String(),
			Size: hdr.Size,
			Name: hdr.Name,
		}
		contents = append(contents, entry)
	}
	return contents, nil
}
