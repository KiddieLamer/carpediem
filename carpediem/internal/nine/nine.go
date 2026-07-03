package nine

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type importRequest struct {
	AccessToken string `json:"accessToken"`
	Name        string `json:"name"`
}

type importResponse struct {
	Success bool `json:"success"`
}

func Import(email, accessToken string) error {
	home, _ := os.UserHomeDir()
	datadir := filepath.Join(home, ".9router")

	machineID, err := os.ReadFile(filepath.Join(datadir, "machine-id"))
	if err != nil {
		return fmt.Errorf("baca machine-id: %w", err)
	}
	cliSecret, err := os.ReadFile(filepath.Join(datadir, "auth", "cli-secret"))
	if err != nil {
		return fmt.Errorf("baca cli-secret: %w", err)
	}

	mid := string(bytes.TrimSpace(machineID))
	sec := string(bytes.TrimSpace(cliSecret))

	h := sha256.Sum256([]byte(mid + "9r-cli-auth" + sec))
	token := hex.EncodeToString(h[:])[:16]

	body, _ := json.Marshal(importRequest{
		AccessToken: accessToken,
		Name:        email[:len(email)-10], // username part
	})

	req, _ := http.NewRequest("POST", "http://localhost:20128/api/oauth/codex/import-token",
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-9r-cli-token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request gagal: %w", err)
	}
	defer resp.Body.Close()

	var result importResponse
	json.NewDecoder(resp.Body).Decode(&result)
	if !result.Success {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(b))
	}

	return nil
}
