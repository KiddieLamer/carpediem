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
	content := `# Format: email:password per baris
# Contoh:
email1@gmail.com:password123
email2@gmail.com:pass456
`
	os.WriteFile(path, []byte(content), 0600)
	fmt.Printf("Created: %s\n", path)
	fmt.Println("Edit, then run: carpediem run")
}

func InputAccounts() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n📋 Paste accounts (email:password), tiap baris satu akun.")
	fmt.Println("Kalo udah, ketik DONE di baris terakhir.\n")

	var lines []string
	for {
		fmt.Print("  > ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if strings.ToUpper(line) == "DONE" || strings.ToUpper(line) == "" {
			break
		}

		if strings.Contains(line, ":") {
			lines = append(lines, line)
		} else {
			fmt.Println("  ⚠️  Format salah. Harus email:password")
		}
	}

	if len(lines) == 0 {
		fmt.Println("  ⏭️  Gak ada akun.")
		return
	}

	content := strings.Join(lines, "\n") + "\n"
	path := accountsFile()
	os.WriteFile(path, []byte(content), 0600)

	fmt.Printf("\n✅ %d akun tersimpan ke %s\n", len(lines), path)
}
