package editor

import (
	"os"
	"os/exec"
)

// OpenInEditor opens content in the system's default editor
func OpenInEditor(content string) error {
    // Create a temporary file
    tmpFile, err := os.CreateTemp("", "goatmeal-*.txt")
    if err != nil {
        return err
    }
    tmpPath := tmpFile.Name()
    defer os.Remove(tmpPath)

    // Write content to file
    if _, err := tmpFile.WriteString(content); err != nil {
        tmpFile.Close()
        return err
    }
    tmpFile.Close()

    // Get the default editor
    editor := getDefaultEditor()

    // Run the editor
    cmd := exec.Command(editor, tmpPath)
    if editor != "cursor" {
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    }

    return cmd.Run()
}

func getDefaultEditor() string {
    if editor := os.Getenv("EDITOR"); editor != "" {
        return editor
    }
    if editor := os.Getenv("VISUAL"); editor != "" {
        return editor
    }
    
    // Try common editors
    for _, editor := range []string{"nvim", "nano", "vim"} {
        if _, err := exec.LookPath(editor); err == nil {
            return editor
        }
    }
    return "vim"
} 