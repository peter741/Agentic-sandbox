// Agentic Sandbox addition: metadata-only workspace exposure review.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type preflightFinding struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	Advice string `json:"advice"`
}
type preflightReport struct {
	Complete       bool               `json:"complete"`
	NeedsReview    bool               `json:"needs_review"`
	EntriesScanned int                `json:"entries_scanned"`
	Findings       []preflightFinding `json:"findings"`
}

var errPreflightLimit = errors.New("preflight entry limit")

func newPreflightCommand() *cobra.Command {
	var workspace string
	var asJSON bool
	var maxEntries int
	cmd := &cobra.Command{
		Use: "preflight", Short: "Review workspace filenames for potential credential exposure before copying",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if maxEntries < 1 || maxEntries > 1000000 {
				return usageErrorf("--max-entries must be between 1 and 1000000")
			}
			root, err := validateWorkspacePath(workspace)
			if err != nil {
				return err
			}
			report := scanPreflight(cmd.Context(), root, maxEntries)
			if asJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				if err := enc.Encode(report); err != nil {
					return runtimeErrorf("write preflight report: %v", err)
				}
			} else {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Workspace exposure review (filenames only; no file contents read)"); err != nil {
					return runtimeErrorf("write preflight report: %v", err)
				}
				for _, finding := range report.Findings {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "[%s] %q: %s\n", finding.Kind, finding.Path, finding.Advice); err != nil {
						return runtimeErrorf("write preflight report: %v", err)
					}
				}
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Scanned %d entries. Complete: %t. Review required: %t.\nNo changes made. This is not a secret-content scan or a security guarantee.\n", report.EntriesScanned, report.Complete, report.NeedsReview); err != nil {
					return runtimeErrorf("write preflight report: %v", err)
				}
			}
			if report.NeedsReview || !report.Complete {
				return exitCodeError(exitCodeRuntimeError)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&workspace, "workspace", ".", "Existing project directory to review")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Print a machine-readable review")
	cmd.Flags().IntVar(&maxEntries, "max-entries", 10000, "Maximum entries inspected; an incomplete scan requires review")
	return cmd
}

func scanPreflight(ctx context.Context, root string, maxEntries int) preflightReport {
	report := preflightReport{Complete: true, Findings: []preflightFinding{}}
	add := func(path, kind, advice string) {
		report.Findings = append(report.Findings, preflightFinding{filepath.ToSlash(path), kind, advice})
		report.NeedsReview = true
	}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if walkErr != nil {
			report.Complete = false
			add(rel, "unreadable", "This entry could not be inspected; check access before relying on the review.")
			return nil
		}
		if rel == "." {
			if !entry.IsDir() {
				report.Complete = false
				add(rel, "invalid_workspace", "Select a project directory.")
				return fs.SkipAll
			}
			return nil
		}
		if report.EntriesScanned >= maxEntries {
			return errPreflightLimit
		}
		report.EntriesScanned++
		if entry.Type()&fs.ModeSymlink != 0 {
			add(rel, "symlink", "Link target was not followed or inspected; review its behavior before copying.")
			return nil
		}
		name := strings.ToLower(entry.Name())
		if entry.IsDir() && (name == ".ssh" || name == ".aws" || name == ".azure" || name == ".kube" || name == ".gnupg") {
			add(rel, "credential_directory", "May contain authentication material; use a clean workspace copy without credentials.")
		}
		if !entry.IsDir() && sensitivePreflightName(name) {
			add(rel, "sensitive_filename", "May contain credentials or private keys; review and remove sensitive material from the copy you give the agent.")
		}
		return nil
	})
	if err != nil {
		report.Complete = false
		if errors.Is(err, errPreflightLimit) {
			add(".", "scan_limit", "Entry limit reached; increase --max-entries or review a smaller workspace.")
		} else {
			add(".", "scan_interrupted", "Scan interrupted; rerun and resolve filesystem access or cancellation.")
		}
	}
	return report
}
func sensitivePreflightName(name string) bool {
	return name == ".env" || strings.HasPrefix(name, ".env.") || name == "id_rsa" || name == "id_ed25519" || name == "id_ecdsa" || name == "id_dsa" || name == ".netrc" || name == ".npmrc" || name == ".pypirc" || name == "credentials" || name == "credentials.json" || name == "application_default_credentials.json" || strings.HasSuffix(name, ".pem") || strings.HasSuffix(name, ".key") || strings.HasSuffix(name, ".p12") || strings.HasSuffix(name, ".pfx")
}
