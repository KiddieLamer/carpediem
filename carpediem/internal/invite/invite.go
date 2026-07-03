package invite

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const targetOrg = "5e4c9b31-1b4e-4887-839b-607597928d7c"

func Send(accessToken string) error {
	url := fmt.Sprintf("https://chatgpt.com/backend-api/accounts/%s/invites/request", targetOrg)
	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("buat request gagal: %w", err)
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("authorization", "Bearer "+accessToken)
	req.Header.Set("oai-language", "en-US")
	req.Header.Set("referer", "https://chatgpt.com/k12-verification")
	req.Header.Set("origin", "https://chatgpt.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request gagal: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}
