package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Diniboy1123/transparnsee/config"
)

var httpClient = &http.Client{}

func CheckNSRecordsCloudflareDoH(domain string) bool {
	url := fmt.Sprintf("https://cloudflare-dns.com/dns-query?name=%s&type=NS", domain)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/dns-json")

	resp, err := httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Answer []struct {
			Data string `json:"data"`
		} `json:"Answer"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	count := 0
	for _, ans := range result.Answer {
		for _, ns := range config.AppConfig.CloudflareNS {
			if ans.Data == ns {
				count++
			}
		}
	}
	return count == len(config.AppConfig.CloudflareNS)
}
