package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func Init() {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".carpediem")
	path := filepath.Join(dir, "accounts.txt")

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Already exists: %s\n", path)
		return
	}

	os.MkdirAll(dir, 0755)
	content := `# Format: email:password
# Contoh:
# email1@gmail.com:password123
# email2@outlook.com:pass456
`
	os.WriteFile(path, []byte(content), 0600)
	fmt.Printf("Created: %s\n", path)
	fmt.Println("Edit the file, then run: carpediem run")
}
