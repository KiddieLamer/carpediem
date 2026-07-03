#!/usr/bin/env python3
"""Auto signup ChatGPT → otomatis isi nama/tgl lahir → capture auth.json → import 9router."""

import json, os, sys, time, base64, hashlib, random
from urllib.request import Request, urlopen
from camoufox import Camoufox

OUTPUT_DIR = os.path.dirname(os.path.abspath(__file__))
ACCOUNTS = os.path.join(OUTPUT_DIR, "accounts.txt")
N9 = os.path.expanduser("~/.9router")

FIRST_NAMES = ["Alice","Bobby","Charlie","Diana","Ethan","Fiona","George","Hannah","Ivan","Julia",
               "Kevin","Luna","Marcus","Nina","Oscar","Penny","Quinn","Riley","Sofia","Tyler",
               "Uma","Victor","Wendy","Xander","Yara","Zane"]

LAST_NAMES = ["Smith","Jones","Brown","Taylor","Wilson","Moore","Anderson","Thomas","Jackson",
              "White","Harris","Martin","Garcia","Martinez","Robinson","Clark","Lewis","Lee","Walker"]

def rand_name(): return f"{random.choice(FIRST_NAMES)} {random.choice(LAST_NAMES)}"
def rand_year(): return random.randint(1975, 2005)

def next_file():
    for i in range(999):
        n = f"auth{'' if i==0 else f'({i})'}.json"
        p = os.path.join(OUTPUT_DIR, n)
        if not os.path.exists(p): return p, n
    return None, None

def to_9router(token, name):
    try:
        m = open(os.path.join(N9, "machine-id")).read().strip()
        s = open(os.path.join(N9, "auth", "cli-secret")).read().strip()
        t = hashlib.sha256(f"{m}9r-cli-auth{s}".encode()).hexdigest()[:16]
        r = Request("http://localhost:20128/api/oauth/codex/import-token",
            json.dumps({"accessToken": token, "name": name}).encode(),
            headers={"Content-Type":"application/json","x-9r-cli-token":t})
        with urlopen(r) as x: return json.loads(x.read()).get("success", False)
    except Exception as e:
        print(f"  ⚠️  9router: {e}"); return False

def wait_visible(p, selector, timeout=15):
    try: p.wait_for_selector(selector, timeout=timeout*1000); return True
    except: return False

def run(email, pw):
    print(f"\n{'='*50}\n📧 {email}\n{'='*50}")
    path, fname = next_file()
    if not path: print("  SKIP"); return False

    name, year = rand_name(), rand_year()
    print(f"  🆔 {name} ({year})")

    with Camoufox(headless=False, humanize=True) as b:
        p = b.new_page()
        # Langsung ke halaman signup
        p.goto("https://chatgpt.com/auth/signup", wait_until="networkidle")
        time.sleep(3)
        p.fill("#email", email); time.sleep(1)
        p.locator("button[type='submit']:has-text('Continue')").click(timeout=5000)
        time.sleep(2)

        # Cek: akun sudah ada? Kalau iya, klik Log in instead
        try:
            if p.locator("text=already exists").is_visible(timeout=3000) or \
               p.locator("text=already have an account").is_visible(timeout=2000):
                print("  🔁 Akun sudah ada, login instead...")
                p.locator("a:has-text('Log in')").first.click(timeout=5000)
                time.sleep(2)
        except: pass

        # Password
        if wait_visible(p, "#password", 5):
            p.fill("#password", pw); time.sleep(1)
            p.locator("button[type='submit']:has-text('Continue')").click(timeout=5000)

        # OTP — user input manual
        print("  📬 OTP dikirim. Input di browser...")
        for i in range(60, 0, -1):
            print(f"\r  ⏳ Nunggu OTP — {i}s ", end="", flush=True)
            time.sleep(1)
        print()

        # Tunggu redirect ke auth.openai.com/about-you
        print("  👤 Ngisi nama & umur...")
        p.wait_for_url("**/about-you**", timeout=60000)
        time.sleep(2)

        # Isi nama
        n_input = p.locator("#name").first
        if n_input.is_visible(timeout=5000):
            n_input.fill(name); time.sleep(1)
            p.locator("button[type='submit']:has-text('Continue')").first.click(timeout=5000)
        else:
            p.locator("input[name='name']").first.fill(name); time.sleep(1)
            p.locator("button[type='submit']").first.click(timeout=5000)

        time.sleep(2)

        # Isi tahun lahir
        try:
            p.wait_for_url("**/about-you**", timeout=10000)
            p.locator("#birthDate").first.fill(str(year)); time.sleep(1)
            p.locator("button[type='submit']").first.click(timeout=5000)
        except:
            try:
                p.locator("input[placeholder*='YYYY']").first.fill(str(year)); time.sleep(1)
                p.locator("button[type='submit']").first.click(timeout=5000)
            except: pass

        # Skip onboarding/redirect sampai balik ke chatgpt.com
        print("  ⏳ Nunggu redirect ke chatgpt.com...")
        for _ in range(30):
            time.sleep(3)
            url = p.url
            if "chatgpt.com" in url and "/auth" not in url and "/onboarding" not in url:
                break
            # Klik tombol apapun yang muncul (Skip/Next/Continue/Done/etc)
            try:
                for btn in ["Skip", "Next", "Continue", "Done", "Get started", "Start chatting", "Go to ChatGPT"]:
                    b = p.locator(f"button:has-text('{btn}')")
                    if b.is_visible(timeout=1000):
                        b.click(timeout=2000); time.sleep(2); break
            except: pass

        print("  ✅ Login berhasil. Mengambil session...")
        time.sleep(2)

        # Redirect ke chatgpt.com kalo masih di domain lain
        if "chatgpt.com" not in p.url:
            p.goto("https://chatgpt.com/", wait_until="networkidle")
            time.sleep(3)

        # Fetch session
        result = p.evaluate("""
            () => fetch('/api/auth/session', {credentials:'same-origin'})
                .then(r => r.json())
                .then(d => JSON.stringify(d, null, 2))
        """)
        data = json.loads(result)

        if not data.get("accessToken"):
            print("  ⚠️  accessToken kosong, refresh...")
            p.goto("https://chatgpt.com/", wait_until="networkidle")
            time.sleep(3)
            result = p.evaluate("""
                () => fetch('/api/auth/session', {credentials:'same-origin'})
                    .then(r => r.json())
                    .then(d => JSON.stringify(d, null, 2))
            """)
            data = json.loads(result)

        if not data.get("accessToken"):
            print("  ❌ Gagal. Skip.")
            return False

        with open(path, "w") as f:
            json.dump(data, f, indent=2)

        em = data.get("user", {}).get("email", "?")
        print(f"\n  ✅ {fname} → {em}")
        if to_9router(data["accessToken"], em.split("@")[0]):
            print("  🚀 9router: OK")

    return True

def main():
    if not os.path.exists(ACCOUNTS):
        print(f"Buat {ACCOUNTS} dulu (email:password per baris)"); sys.exit(1)
    with open(ACCOUNTS) as f:
        lines = [l.strip() for l in f if l.strip() and ":" in l]
    if not lines: print("accounts.txt kosong"); sys.exit(1)
    print(f"📋 {len(lines)} akun\n")
    for i, l in enumerate(lines, 1):
        e, pw = l.split(":", 1)
        print(f"\n[{i}/{len(lines)}]"); run(e, pw)
    print(f"\n{'='*50}\n✅ Selesai!")

if __name__ == "__main__":
    main()
