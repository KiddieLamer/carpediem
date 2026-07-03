package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/KiddieLamer/carpediem/internal/accounts"
	"github.com/KiddieLamer/carpediem/internal/browser"
	"github.com/KiddieLamer/carpediem/internal/nine"
)

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
		cancel()
	}()

	for i, acc := range accs {
		select {
		case <-ctx.Done():
			fmt.Println("\n⚠️  Dihentikan.")
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

		if dry {
			fmt.Println("  ⏭️  Dry mode — skip 9router")
			continue
		}

		if err := nine.Import(acc.Email, session.AccessToken); err != nil {
			fmt.Printf("  ⚠️ 9router: %v\n", err)
		} else {
			fmt.Println("  🚀 9router: OK")
		}
	}

	fmt.Printf("\n%s\n✅ Selesai!\n", strings.Repeat("=", 50))
}
