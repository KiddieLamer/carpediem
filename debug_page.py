#!/usr/bin/env python3
from camoufox import Camoufox

with Camoufox(headless=False) as b:
    p = b.new_page()
    p.goto("https://chatgpt.com/auth/signup", wait_until="networkidle")
    import time; time.sleep(5)
    p.screenshot(path="/tmp/chatgpt_debug.png")
    html = p.content()
    with open("/tmp/chatgpt_debug.html", "w") as f:
        f.write(html)
    print("Screenshot: /tmp/chatgpt_debug.png")
    print("HTML: /tmp/chatgpt_debug.html")
    input("Tekan Enter...")
