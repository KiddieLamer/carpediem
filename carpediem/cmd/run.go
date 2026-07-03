package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/KiddieLamer/carpediem/internal/accounts"
	"github.com/KiddieLamer/carpediem/internal/browser"
	"github.com/KiddieLamer/carpediem/internal/invite"
	"github.com/KiddieLamer/carpediem/internal/nine"
)

func RunInteractive() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\n📁 Path ke accounts.txt (enter buat default ~/.carpediem/accounts.txt): ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)

	fmt.Print("⏱️  OTP delay detik (enter buat default 60): ")
	delayStr, _ := reader.ReadString('\n')
	delayStr = strings.TrimSpace(delayStr)

	delay := 60
	if delayStr != "" {
		fmt.Sscanf(delayStr, "%d", &delay)
	}

	fmt.Print("🧪 Dry run? (skip 9router) [y/N]: ")
	dryStr, _ := reader.ReadString('\n')
	dry := strings.TrimSpace(strings.ToLower(dryStr)) == "y"

	Run(path, delay, dry)
}

func Run(accountsPath string, otpDelay int, dry bool) {
	accs, err := accounts.Load(accountsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\n⚠️  Dihentikan.")
		cancel()
	}()

	for i, acc := range accs {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Printf("\n%s\n", strings.Repeat("=", 50))
		fmt.Printf("[%d/%d] 📧 %s\n", i+1, len(accs), acc.Email)
		fmt.Println(strings.Repeat("=", 50))

		progress := make(chan string, 100)
		go func() {
			for msg := range progress {
				fmt.Println(msg)
			}
		}()

		session, err := browser.Run(ctx, acc.Email, acc.Password, otpDelay, progress)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ❌ %v\n", err)
			continue
		}

		email := "?"
		if session.User != nil {
			email = session.User.Email
		}
		fmt.Printf("  ✅ Login: %s\n", email)

		// Simpan session ke ~/.carpediem/<email>.json
		sessionPath := filepath.Join(os.Getenv("HOME"), ".carpediem", email+".json")
		sessionData, _ := json.MarshalIndent(session, "", "  ")
		os.WriteFile(sessionPath, sessionData, 0600)
		fmt.Printf("  💾 Session saved: %s\n", sessionPath)

		if dry {
			fmt.Println("  ⏭️  Dry mode")
			continue
		}

		fmt.Println("  📨 Kirim invite request...")
		if err := invite.Send(session.AccessToken); err != nil {
			fmt.Printf("  ⚠️ Invite: %v\n", err)
		} else {
			fmt.Println("  ✅ Invite OK")
		}

		if err := nine.Import(acc.Email, session.AccessToken); err != nil {
			fmt.Printf("  ⚠️ 9router: %v\n", err)
		} else {
			fmt.Println("  🚀 9router: OK")
		}
	}

	fmt.Printf("\n%s\n✅ Selesai!\n", strings.Repeat("=", 50))
}
