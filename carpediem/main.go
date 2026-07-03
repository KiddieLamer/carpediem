package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/KiddieLamer/carpediem/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: carpediem <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  init         Create accounts.txt template")
		fmt.Println("  run          Run automation")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cmd.Init()
	case "run":
		fs := flag.NewFlagSet("run", flag.ExitOnError)
		accounts := fs.String("accounts", "", "Path to accounts.txt")
		delay := fs.Int("delay", 60, "OTP wait delay in seconds")
		dry := fs.Bool("dry", false, "Skip 9router import")
		fs.Parse(os.Args[2:])
		cmd.Run(*accounts, *delay, *dry)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
