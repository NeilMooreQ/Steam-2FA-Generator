package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

type Response struct {
	Response struct {
		ServerTime string `json:"server_time"`
	} `json:"response"`
}

func getServerTime() (int64, error) {
	url := "https://api.steampowered.com:443/ITwoFactorService/QueryTime/v0001"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		return 0, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error response: %s", resp.Status)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("error decoding JSON: %w", err)
	}

	serverTime, err := strconv.ParseInt(response.Response.ServerTime, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing server_time: %w", err)
	}

	return serverTime, nil
}

func generate2fa(secretKey string) (string, error) {
	code := ""
	char := "23456789BCDFGHJKMNPQRTVWXY"

	timeValue, err := getServerTime()
	if err != nil {
		return "", err
	}

	hexTime := fmt.Sprintf("%016x", timeValue/30)

	byteTime, err := hex.DecodeString(hexTime)
	if err != nil {
		return "", fmt.Errorf("error decoding hex: %w", err)
	}

	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error decoding key: %w", err)
	}

	h := hmac.New(sha1.New, key)
	_, err = h.Write(byteTime)
	if err != nil {
		return "", fmt.Errorf("error writing to HMAC: %w", err)
	}
	digest := h.Sum(nil)

	begin := int(digest[19]) & 0x0F

	cInt := int32(digest[begin])<<24 | int32(digest[begin+1])<<16 | int32(digest[begin+2])<<8 | int32(digest[begin+3])
	cInt = cInt & 0x7fffffff

	for r := 0; r < 5; r++ {
		code += string(char[cInt%int32(len(char))])
		cInt /= int32(len(char))
	}

	return code, nil
}

func main() {
	secretKey := flag.String("secret_key", "", "The shared secret key for 2FA")
	flag.Parse()

	if *secretKey == "" {
		fmt.Println("Error: secret_key is required")
		return
	}

	code, err := generate2fa(*secretKey)
	if err != nil {
		fmt.Println("Error generating code:", err)
		return
	}

	fmt.Println(code)
}
