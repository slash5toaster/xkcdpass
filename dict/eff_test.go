package dict

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
    return f(req)
}

func TestCountUniqueLines(t *testing.T) {
    path := filepath.Join(t.TempDir(), "words.txt")
    content := "apple\nbanana\napple\ncherry\n"
    if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
        t.Fatalf("write temp file: %v", err)
    }

    count, err := countUniqueLines(path)
    if err != nil {
        t.Fatalf("countUniqueLines returned error: %v", err)
    }
    if count != 3 {
        t.Fatalf("expected 3 unique lines, got %d", count)
    }
}

func TestGetEFFDictPathUsesValidCachedFile(t *testing.T) {
    cacheRoot := t.TempDir()
    t.Setenv("XDG_CACHE_HOME", cacheRoot)

    cacheDir := filepath.Join(cacheRoot, "xkcdpass")
    if err := os.MkdirAll(cacheDir, 0o755); err != nil {
        t.Fatalf("create cache dir: %v", err)
    }

    cachedPath := filepath.Join(cacheDir, "eff_large_wordlist.txt")
    lines := make([]string, minEFFWordlistLines+10)
    for i := range lines {
        lines[i] = fmt.Sprintf("%d\tword-%d", i, i)
    }
    if err := os.WriteFile(cachedPath, []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
        t.Fatalf("write cached dict: %v", err)
    }

    got, err := GetEFFDictPath()
    if err != nil {
        t.Fatalf("GetEFFDictPath returned error: %v", err)
    }
    if got != cachedPath {
        t.Fatalf("expected %q, got %q", cachedPath, got)
    }
}

func TestGetEFFDictPathRefreshesSmallCache(t *testing.T) {
    oldTransport := http.DefaultTransport
    http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
        if req.URL.String() != effWordlistURL {
            t.Fatalf("unexpected URL requested: %s", req.URL.String())
        }

        lines := make([]string, minEFFWordlistLines+5)
        for i := range lines {
            lines[i] = fmt.Sprintf("%d\tword-%d", i, i)
        }

        return &http.Response{
            StatusCode: http.StatusOK,
            Body:       io.NopCloser(strings.NewReader(strings.Join(lines, "\n") + "\n")),
            Header:     make(http.Header),
            Request:    req,
        }, nil
    })
    defer func() { http.DefaultTransport = oldTransport }()

    cacheRoot := t.TempDir()
    t.Setenv("XDG_CACHE_HOME", cacheRoot)

    cacheDir := filepath.Join(cacheRoot, "xkcdpass")
    if err := os.MkdirAll(cacheDir, 0o755); err != nil {
        t.Fatalf("create cache dir: %v", err)
    }

    cachedPath := filepath.Join(cacheDir, "eff_large_wordlist.txt")
    if err := os.WriteFile(cachedPath, []byte("too\nshort\n"), 0o600); err != nil {
        t.Fatalf("write small cache file: %v", err)
    }

    got, err := GetEFFDictPath()
    if err != nil {
        t.Fatalf("GetEFFDictPath returned error: %v", err)
    }
    if got != cachedPath {
        t.Fatalf("expected %q, got %q", cachedPath, got)
    }

    count, err := countUniqueLines(cachedPath)
    if err != nil {
        t.Fatalf("countUniqueLines returned error: %v", err)
    }
    if count < minEFFWordlistLines {
        t.Fatalf("expected refreshed cache to have at least %d unique lines, got %d", minEFFWordlistLines, count)
    }
}
