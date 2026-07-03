# CarpeDiem — ChatGPT Auth Automator

> Single binary untuk registrasi/login ChatGPT massal + import session ke 9router.
> Bypass OTP SMS / region lock dengan cara ekstrak session dari browser.

## Stack

Go + Rod (Chrome DevTools Protocol) — **pure native binary, no Python, no Node**.

## Setup

```bash
# 1. Clone
git clone https://github.com/KiddieLamer/carpediem.git
cd carpediem

# 2. Build
go build -o carpediem .

# 3. Init accounts file
./carpediem init
# → Creates ~/.carpediem/accounts.txt

# 4. Edit accounts.txt
nano ~/.carpediem/accounts.txt
# Format:
# email:password
# h.omeofbad@gmail.com:pass123
# h.o.meofbad@gmail.com:pass456
```

## Usage

```bash
# Run semua akun di accounts.txt
./carpediem run

# Run dengan file custom
./carpediem run --accounts ./list.txt

# Dry run (skip import ke 9router)
./carpediem run --dry

# Custom OTP delay (default 60s)
./carpediem run --delay 90
```

## Flow per Akun

```
1. Buka Chromium → chatgpt.com/auth/signup
2. Isi email & password otomatis
3. ❓ Kalau "already exists" → switch ke login
4. ⏳ Countdown OTP — kamu input manual di browser
5. 👤 Isi nama & tahun lahir random (otomatis)
6. ⏭️ Skip onboarding otomatis
7. ✅ Fetch session → simpan auth.json
8. 🚀 Import ke 9router (POST /api/oauth/codex/import-token)
```

Kamu cukup input OTP dari email. Sisanya otomatis.

## Binary

File `carpediem` adalah **single binary native arm64** (14MB). Tinggal download + run, Chromium didownload otomatis pas pertama jalan.
