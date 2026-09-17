// Package dict handles fetching and caching the EFF large wordlist.
package dict

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const effWordlistURL = "https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt"

const minEFFWordlistLines = 7700

// ----------------------------------------------------------------------------\\
func downloadEFFWordlist(destPath string) error {
	resp, err := http.Get(effWordlistURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download EFF wordlist: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// the EFF wordlist is tab-separated "diceroll\\tword" per line;
	// keep only the word so getRandomWord's punctuation stripping doesn't mangle it
	scanner := bufio.NewScanner(resp.Body)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		word := fields[len(fields)-1]
		if word == "" {
			continue
		}
		fmt.Fprintln(writer, word)
	}
	return scanner.Err()
}

// ----------------------------------------------------------------------------\\
func countUniqueLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	uniqueLines := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		uniqueLines[scanner.Text()] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return len(uniqueLines), nil
}

// ----------------------------------------------------------------------------\\
// GetEFFDictPath returns the path to a locally cached copy of the EFF large
// wordlist, downloading it first if it isn't already cached.
func GetEFFDictPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}

	xkcdCacheDir := filepath.Join(cacheDir, "xkcdpass")
	if err := os.MkdirAll(xkcdCacheDir, 0o755); err != nil {
		return "", err
	}

	destPath := filepath.Join(xkcdCacheDir, "eff_large_wordlist.txt")

	needsDownload := false
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		needsDownload = true
	} else if err != nil {
		return "", err
	} else {
		uniqueLines, err := countUniqueLines(destPath)
		if err != nil {
			return "", err
		}
		if uniqueLines < minEFFWordlistLines {
			if err := os.Remove(destPath); err != nil {
				return "", err
			}
			needsDownload = true
		}
	}

	if needsDownload {
		if err := downloadEFFWordlist(destPath); err != nil {
			return "", err
		}
	}

	return destPath, nil
}
