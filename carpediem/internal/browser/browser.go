package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

var firstNames = []string{"Alice","Bobby","Charlie","Diana","Ethan","Fiona",
	"George","Hannah","Ivan","Julia","Kevin","Luna","Marcus","Nina",
	"Oscar","Penny","Quinn","Riley","Sofia","Tyler","Uma","Victor",
	"Wendy","Xander","Yara","Zane"}

var lastNames = []string{"Smith","Jones","Brown","Taylor","Wilson","Moore",
	"Anderson","Thomas","Jackson","White","Harris","Martin","Garcia",
	"Martinez","Robinson","Clark","Lewis","Lee","Walker"}

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%")

func randName() string {
	return firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
}

func randPass() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func randYear() int { return 1975 + rand.Intn(31) }

func Run(ctx context.Context, email, password string, otpDelay int, progress chan<- string) (*Session, error) {
	defer close(progress)

	name := randName()
	year := randYear()
	progress <- fmt.Sprintf("  🆔 %s (%d)", name, year)

	progress <- "  🚀 Launching browser..."
	u := launcher.New().
		Headless(false).
		Logger(io.Discard).
		Set("disable-blink-features", "AutomationControlled").
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36").
		MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.Close()

	// Stealth: bypass automation detection
	page := browser.MustPage("https://chatgpt.com/").Context(ctx)
	page.MustEval(`() => {
		Object.defineProperty(navigator, "webdriver", { get: () => false });
		window.chrome = { runtime: {} };
		Object.defineProperty(navigator, "plugins", { get: () => [1,2,3] });
		Object.defineProperty(navigator, "languages", { get: () => ["en-US","en"] });
	}`)
	page.MustWaitLoad()
	time.Sleep(3)

	progress <- "  🔘 Klik Log in..."
	page.MustElement(`button[data-testid="login-button"]`).MustClick()
	page.MustWaitStable()
	time.Sleep(3)

	url := page.MustInfo().URL
	progress <- fmt.Sprintf("  🌐 URL: %s", url)

	// Isi email (modal langsung muncul dengan #email)
	progress <- "  📧 Isi email..."
	page.Timeout(15 * time.Second).MustElement("#email").MustInput(email)
	time.Sleep(1)
	page.Timeout(15 * time.Second).MustElement(`button[type="submit"]`).MustClick()

	// Tunggu redirect ke auth.openai.com/email-verification
	progress <- "  🌐 Nunggu redirect ke email-verification..."
	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		u := page.MustInfo().URL
		if strings.Contains(u, "email-verification") {
			break
		}
		time.Sleep(2)
	}
	progress <- fmt.Sprintf("  🌐 URL: %s", page.MustInfo().URL)

	// OTP — user input manual di browser
	progress <- "  📬 Cek email, input OTP di browser..."
	for i := otpDelay; i > 0; i-- {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			progress <- fmt.Sprintf("  ⏳ Nunggu OTP — %ds", i)
			time.Sleep(1 * time.Second)
		}
	}

	// Tunggu sampai page meninggalkan email-verification (OTP berhasil diverifikasi)
	progress <- "  ⏳ Nunggu verifikasi OTP..."
	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		u := page.MustInfo().URL
		if !strings.Contains(u, "email-verification") {
			break
		}
		time.Sleep(2)
	}
	progress <- fmt.Sprintf("  🌐 URL: %s", page.MustInfo().URL)

	// about-you (nama + umur) — kalo langsung ke chatgpt.com berarti akun existing
	progress <- "  👤 Isi Full name & Age..."
	time.Sleep(2)

	// Cari input name — label text "Full name"
	if el, err := page.ElementByJS(rod.Eval(`() => {
		const label = [...document.querySelectorAll('div')].find(d => d.textContent.trim() === 'Full name');
		if (!label) return null;
		const container = label.closest('div')?.parentElement;
		return container?.querySelector('input');
	}`)); err == nil {
		el.MustInput(name); progress <- "  ✅ Name diisi"
		time.Sleep(1)
	} else {
		progress <- "  ⚠️ Input name gak ketemu"
	}

	if btn, err := page.Element(`button[type="submit"]`); err == nil {
		btn.MustClick(); time.Sleep(2)
	}

	// Cari input age — label text "Age" atau "Birth date"
	if el, err := page.ElementByJS(rod.Eval(`() => {
		const label = [...document.querySelectorAll('div')].find(d => d.textContent.trim() === 'Age');
		if (!label) return null;
		const container = label.closest('div')?.parentElement;
		return container?.querySelector('input');
	}`)); err == nil {
		el.MustInput(fmt.Sprintf("%d", year)); progress <- "  ✅ Age diisi"
		time.Sleep(1)
	} else {
		progress <- "  ⚠️ Input age gak ketemu"
	}

	// Klik Finish / Continue untuk submit age
	if btn, err := page.Element(`button[type="submit"]`); err == nil {
		btn.MustClick(); time.Sleep(2)
	}

	// Skip onboarding
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
		for _, txt := range []string{"Skip","Next","Done","Get started","Start chatting","Go to ChatGPT"} {
			if btn, err := page.ElementByJS(rod.Eval(`() => [...document.querySelectorAll('button')].find(b => b.textContent.trim() === '`+txt+`')`)); err == nil {
				btn.MustClick(); time.Sleep(2); break
			}
		}
		if btn, err := page.Element(`button[type="submit"]`); err == nil {
			btn.MustClick(); time.Sleep(2)
		}
		time.Sleep(3)
	}

	progress <- "  ✅ Login. Mengambil session..."
	url = page.MustInfo().URL
	if !strings.Contains(url, "chatgpt.com") {
		page.MustNavigate("https://chatgpt.com/")
		page.MustWaitLoad()
		time.Sleep(3)
	}

	result := page.MustEval(`() => fetch('/api/auth/session', {credentials:'same-origin'}).then(r => r.json()).then(d => JSON.stringify(d))`)
	var session Session
	json.Unmarshal([]byte(result.Str()), &session)

	if session.AccessToken == "" {
		return nil, fmt.Errorf("accessToken kosong")
	}
	return &session, nil
}
