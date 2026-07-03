package accounts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Account struct {
	Email    string
	Password string
}

func Load(path string) ([]Account, error) {
	var p string
	if path != "" {
		p = path
	} else {
		home, _ := os.UserHomeDir()
		p = filepath.Join(home, ".carpediem", "accounts.txt")
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("baca file gagal: %w.\nBuat dulu: carpediem init", err)
	}

	var accounts []Account
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		accounts = append(accounts, Account{
			Email:    strings.TrimSpace(parts[0]),
			Password: strings.TrimSpace(parts[1]),
		})
		_ = i
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("accounts.txt kosong atau format salah")
	}
	return accounts, nil
}
