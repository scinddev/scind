// internal/context/detector.go
// Context detection logic

package context

import (
	"fmt"
	"os"
	"path/filepath"
)

// Context holds the detected workspace and application context
type Context struct {
	WorkspacePath string
	WorkspaceName string
	AppPath       string
	AppName       string
}

// DetectContext walks up the directory tree to find workspace and application markers
func DetectContext(startDir string) (*Context, error) {
	ctx := &Context{}

	// Phase 1: Find workspace.yaml (walk up from startDir)
	workspaceRoot, err := findMarkerFile(startDir, "workspace.yaml")
	if err != nil {
		// Check if there's an orphaned application.yaml
		appPath, appErr := findMarkerFile(startDir, "application.yaml")
		if appErr == nil {
			return nil, fmt.Errorf(
				"no workspace found (workspace.yaml) in current directory or any parent directories, "+
					"but found an application (application.yaml) at: %s", appPath)
		}
		return nil, fmt.Errorf(
			"no workspace found (workspace.yaml) in current directory or any parent directories, " +
				"and no application (application.yaml) found either")
	}

	ctx.WorkspacePath = workspaceRoot
	ctx.WorkspaceName = filepath.Base(workspaceRoot) // Or parse from workspace.yaml

	// Phase 2: Find application.yaml (walk up from startDir, but only within workspace)
	appPath, err := findMarkerFileWithinBoundary(startDir, "application.yaml", workspaceRoot)
	if err == nil {
		ctx.AppPath = appPath
		ctx.AppName = filepath.Base(appPath) // Or parse from application.yaml
	}

	return ctx, nil
}

// findMarkerFile walks up from startDir looking for a file with the given name
func findMarkerFile(startDir, filename string) (string, error) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, filename)
		if _, err := os.Stat(candidate); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			return "", fmt.Errorf("%s not found", filename)
		}
		dir = parent
	}
}

// findMarkerFileWithinBoundary walks up from startDir looking for filename,
// but stops at boundaryDir (never traverses above it)
func findMarkerFileWithinBoundary(startDir, filename, boundaryDir string) (string, error) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, filename)
		if _, err := os.Stat(candidate); err == nil {
			return dir, nil
		}

		// Check if we've reached the boundary
		if dir == boundaryDir {
			return "", fmt.Errorf("%s not found within workspace", filename)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root (shouldn't happen if boundaryDir is above startDir)
			return "", fmt.Errorf("%s not found", filename)
		}
		dir = parent
	}
}
