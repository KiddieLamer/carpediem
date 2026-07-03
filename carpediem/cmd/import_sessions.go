package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/KiddieLamer/carpediem/internal/browser"
	"github.com/KiddieLamer/carpediem/internal/nine"
)

func ImportSessions() {
	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, "Downloads")

	fmt.Printf("\n📁 Cari session JSON di direktori (enter buat %s): ", defaultDir)
	var dir string
	fmt.Scanln(&dir)
	if dir == "" {
		dir = defaultDir
	}

	var files []string
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, "session") && strings.HasSuffix(name, ".json") {
			files = append(files, path)
		}
		return nil
	})

	if len(files) == 0 {
		fmt.Println("❌ Gak ada file session*.json di direktori itu.")
		return
	}

	fmt.Printf("\n📄 Nemuin %d file session:\n", len(files))
	for i, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("  %d. %s — ⚠️ gagal baca: %v\n", i+1, path, err)
			continue
		}

		var session browser.Session
		if err := json.Unmarshal(data, &session); err != nil {
			fmt.Printf("  %d. %s — ⚠️ gagal parse: %v\n", i+1, path, err)
			continue
		}

		if session.AccessToken == "" {
			fmt.Printf("  %d. %s — ⚠️ accessToken kosong\n", i+1, path)
			continue
		}

		email := "?"
		if session.User != nil {
			email = session.User.Email
		}

		fmt.Printf("  %d. %s 📧 %s\n", i+1, path, email)

		if err := nine.Import(email, session.AccessToken); err != nil {
			fmt.Printf("     ⚠️ 9router: %v\n", err)
		} else {
			fmt.Println("     🚀 9router: OK")
		}
	}
}
