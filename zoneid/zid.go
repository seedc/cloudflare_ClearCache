package zoneid

import (
	"cloudflare_ClearCache/settings"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type CloudflareResponse struct {
	Result  []Zone  `json:"result"`
	Success bool    `json:"success"`
	Errors  []Error `json:"errors"`
}

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// removeWWWPrefix 移除域名前的 'www.' 前缀
func removeWWWPrefix(domain string) string {
	if strings.HasPrefix(domain, "www.") {
		return strings.TrimPrefix(domain, "www.")
	}
	return domain
}

// getMainDomain 获取主域名
func getMainDomain(domain string) string {
	parsedURL, err := url.Parse(domain)
	if err != nil || parsedURL.Hostname() == "" {
		// 如果无法解析为 URL 或者没有主机名部分，直接返回输入的域名
		return domain
	}

	// 获取主机名部分
	hostname := parsedURL.Hostname()

	// 分割主机名部分
	parts := strings.Split(hostname, ".")
	if len(parts) > 2 {
		// 如果分割后的部分大于2，则取最后两个部分作为主域名
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return hostname
}

func Getzid(domain string) (zoneID string) {
	// 替换为你的API信息
	apiEmail := settings.Conf.Email
	apiKey := settings.Conf.Token

	// 清洗 URL，获取主域名
	domain = getMainDomain(removeWWWPrefix(domain))
	fmt.Println(domain)

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?name=%s", domain)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Add("X-Auth-Email", apiEmail)
	req.Header.Add("X-Auth-Key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code")
		os.Exit(1)
	}

	var cfResponse CloudflareResponse
	err = json.NewDecoder(resp.Body).Decode(&cfResponse)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		os.Exit(1)
	}

	if !cfResponse.Success {
		fmt.Println("Failed to fetch Zone ID")
		for _, e := range cfResponse.Errors {
			fmt.Printf("Error: %d - %s\n", e.Code, e.Message)
		}
		os.Exit(1)
	}

	if len(cfResponse.Result) > 0 {
		zoneID := cfResponse.Result[0].ID
		return fmt.Sprintf(zoneID)
	} else {
		fmt.Println("No zones found for the specified domain")
		return "nil"
	}
}
