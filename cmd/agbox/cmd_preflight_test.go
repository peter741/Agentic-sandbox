package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreflightFindingsAndNoContentLeak(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{".env", "id_ed25519", "public.pem", "README.md"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("SECRET-CONTENT-MUST-NOT-LEAK"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(root, ".aws"), 0700); err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	code := run(context.Background(), []string{"preflight", "--workspace", root, "--json"}, &out, &stderr, os.LookupEnv)
	if code != 1 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%s", code, &stderr)
	}
	var report preflightReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if !report.Complete || !report.NeedsReview || len(report.Findings) != 4 {
		t.Fatalf("%+v", report)
	}
	if strings.Contains(out.String(), "SECRET-CONTENT") || strings.Contains(out.String(), root) {
		t.Fatal("report leaked contents or absolute root")
	}
	contents, err := os.ReadFile(filepath.Join(root, ".env"))
	if err != nil || string(contents) != "SECRET-CONTENT-MUST-NOT-LEAK" {
		t.Fatal("workspace changed")
	}
}
func TestPreflightSymlinkNotFollowed(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, ".env"), []byte("hidden"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	report := scanPreflight(context.Background(), root, 100)
	if !report.Complete || report.EntriesScanned != 1 || len(report.Findings) != 1 || report.Findings[0].Kind != "symlink" {
		t.Fatalf("%+v", report)
	}
}
func TestPreflightLimitAndCancellation(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a", "b"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	report := scanPreflight(context.Background(), root, 1)
	if report.Complete || report.EntriesScanned != 1 || report.Findings[0].Kind != "scan_limit" {
		t.Fatalf("%+v", report)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	report = scanPreflight(ctx, root, 100)
	if report.Complete || !report.NeedsReview || report.EntriesScanned != 0 {
		t.Fatalf("%+v", report)
	}
}
func TestPreflightCleanAndUsage(t *testing.T) {
	root := t.TempDir()
	for _, tc := range []struct {
		args []string
		code int
	}{
		{[]string{"preflight", "--workspace", root, "--json"}, 0},
		{[]string{"preflight", "--workspace", root, "--max-entries", "0"}, 2},
		{[]string{"preflight", "--workspace", "/"}, 2},
		{[]string{"preflight", "extra"}, 2},
	} {
		var out, stderr bytes.Buffer
		if code := run(context.Background(), tc.args, &out, &stderr, os.LookupEnv); code != tc.code {
			t.Fatalf("%v: %d, %s", tc.args, code, &stderr)
		}
	}
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, nil, 0600); err != nil {
		t.Fatal(err)
	}
	report := scanPreflight(context.Background(), file, 10)
	if report.Complete || !report.NeedsReview {
		t.Fatalf("%+v", report)
	}
}
