package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Session struct {
	AccessToken string `json:"accessToken"`
	User        *struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}

var firstNames = []string{"Alice", "Bobby", "Charlie", "Diana", "Ethan", "Fiona",
	"George", "Hannah", "Ivan", "Julia", "Kevin", "Luna", "Marcus", "Nina",
	"Oscar", "Penny", "Quinn", "Riley", "Sofia", "Tyler", "Uma", "Victor",
	"Wendy", "Xander", "Yara", "Zane"}

var lastNames = []string{"Smith", "Jones", "Brown", "Taylor", "Wilson", "Moore",
	"Anderson", "Thomas", "Jackson", "White", "Harris", "Martin", "Garcia",
	"Martinez", "Robinson", "Clark", "Lewis", "Lee", "Walker"}

func randName() string {
	return firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
}

func randYear() int { return 1975 + rand.Intn(31) }

func Run(ctx context.Context, email, password string, otpDelay int, progress chan<- string) (*Session, error) {
	defer close(progress)

	name := randName()
	year := randYear()
	progress <- fmt.Sprintf("  🆔 %s (%d)", name, year)

	u := launcher.New().Headless(false).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.Close()

	page := browser.MustPage("https://chatgpt.com/auth/signup").Context(ctx)
	page.MustWaitLoad()
	time.Sleep(3 * time.Second)

	// Step 1: coba signup (email + password)
	progress <- "  📧 Isi email..."
	page.MustElement("#email").MustInput(email)
	time.Sleep(1)
	page.MustElementR("button", "Continue").MustClick()
	time.Sleep(2)

	// Cek "already exists" — kalau iya, retry dari awal pake login
	if has, _, _ := page.Has("text=already exists"); has {
		progress <- "  🔁 Akun sudah ada. Retry login..."
		page.Close()

		page = browser.MustPage("https://chatgpt.com/auth/login").Context(ctx)
		page.MustWaitLoad()
		time.Sleep(3)

		page.MustElement("#email").MustInput(email)
		time.Sleep(1)
		page.MustElementR("button", "Continue").MustClick()
		time.Sleep(2)

		page.MustElement("#password").MustInput(password)
		time.Sleep(1)
		page.MustElementR("button", "Continue").MustClick()
	}

	// Step 2: OTP
	progress <- "  📬 Input OTP di browser..."
	for i := otpDelay; i > 0; i-- {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			progress <- fmt.Sprintf("  ⏳ Nunggu OTP — %ds", i)
			time.Sleep(1 * time.Second)
		}
	}

	// Step 3: about-you (nama + umur)
	progress <- "  👤 Isi nama & umur..."
	time.Sleep(2)

	if el, err := page.Element("#name"); err == nil {
		el.MustInput(name)
		time.Sleep(1)
		if btn, err := page.ElementR("button", "Continue"); err == nil {
			btn.MustClick()
		}
		time.Sleep(2)
	}

	if el, err := page.Element("#birthDate"); err == nil {
		el.MustInput(fmt.Sprintf("%d", year))
		time.Sleep(1)
		if btn, err := page.ElementR("button", "Continue"); err == nil {
			btn.MustClick()
		}
		time.Sleep(2)
	}

	// Step 4: skip onboarding / tunggu redirect ke chatgpt.com
	progress <- "  ⏳ Nunggu redirect ke chatgpt.com..."
	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		url := page.MustInfo().URL
		if strings.Contains(url, "chatgpt.com") && !strings.Contains(url, "/auth") {
			break
		}

		for _, txt := range []string{"Skip", "Next", "Continue", "Done", "Get started", "Start chatting", "Go to ChatGPT"} {
			if btn, err := page.ElementR("button", txt); err == nil {
				btn.MustClick()
				time.Sleep(2)
				break
			}
		}
		time.Sleep(3)
	}

	progress <- "  ✅ Login berhasil. Mengambil session..."

	url := page.MustInfo().URL
	if !strings.Contains(url, "chatgpt.com") {
		page.MustNavigate("https://chatgpt.com/")
		page.MustWaitLoad()
		time.Sleep(3)
	}

	result := page.MustEval(`
		() => fetch('/api/auth/session', {credentials:'same-origin'})
			.then(r => r.json())
			.then(d => JSON.stringify(d))
	`)
	sessionJSON := result.Str()

	var session Session
	json.Unmarshal([]byte(sessionJSON), &session)

	if session.AccessToken == "" {
		return nil, fmt.Errorf("accessToken kosong")
	}

	return &session, nil
}
