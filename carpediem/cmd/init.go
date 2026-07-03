package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const accountsPath = "accounts.txt"

func accountsFile() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".carpediem", accountsPath)
}

func Init() {
	path := accountsFile()
	dir := filepath.Dir(path)

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Already exists: %s\n", path)
		return
	}

	os.MkdirAll(dir, 0755)
	content := "# Format: email per baris (password opsional: email:pass)\n"
	os.WriteFile(path, []byte(content), 0600)
	fmt.Printf("Created: %s\n", path)
}

func InputAccounts() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n📋 Paste email, tiap baris satu akun.")
	fmt.Println("Kalo udah, ketik DONE.\n")

	var lines []string
	for {
		fmt.Print("  > ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if strings.ToUpper(line) == "DONE" || line == "" {
			break
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		fmt.Println("  ⏭️  Gak ada akun.")
		return
	}

	content := strings.Join(lines, "\n") + "\n"
	os.WriteFile(accountsFile(), []byte(content), 0600)
	fmt.Printf("\n✅ %d akun tersimpan\n", len(lines))
}
