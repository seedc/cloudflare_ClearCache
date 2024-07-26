package cache

import (
	"bytes"
	"cloudflare_ClearCache/settings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type PurgeCacheResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

// removeWWWPrefix 函数用于移除域名前的 'www.' 前缀
func removeWWWPrefix(domain string) string {
	// 检查域名是否以 'www.' 开头，如果是，则移除前缀
	if strings.HasPrefix(domain, "www.") {
		return strings.TrimPrefix(domain, "www.")
	}
	// 如果没有 'www.' 前缀，则返回原始域名
	return domain
}

// isMainDomain 函数检查 URL 是否仅为主域名（无路径或文件）
func isMainDomain(domain string) (bool, error) {
	// 解析传入的域名或 URL
	parsedURL, err := url.Parse(domain)
	if err != nil {
		return false, err
	}

	// 如果路径为空或为根路径（"/"），则认为是主域名
	if parsedURL.Path == "" || parsedURL.Path == "/" {
		return true, nil
	}

	// 否则认为是带有路径或文件的 URL
	return false, nil
}

// 此接口备用 用于判断是否是cf还是aws托管dns域名
func IfDomain(domain string, zonesid string) {
	apiKey := settings.Conf.Token
	apiEmail := settings.Conf.Email

	//判断是url还是域名
	isMain, err := isMainDomain(domain)
	if err != nil {
		fmt.Printf("解析URL时出错 %s: %v\n", domain, err)
		return
	}

	if isMain {
		fmt.Printf("URL: %s is a main domain.\n", domain)
		// 移除可能存在的 'www.' 前缀
		domainWithoutWWW := removeWWWPrefix(domain)
		// 检查去除 'www.' 前缀后的域名是否为主域名
		_, err := isMainDomain(domainWithoutWWW)
		if err != nil {
			fmt.Printf("解析 URL %s 时出错: %v\n", domainWithoutWWW, err)
			return
		} else {
			//直接清理整个域名
			DelCfCache(domainWithoutWWW, apiKey, apiEmail, zonesid)
		}
	} else {
		//执行文件单独清理
		fmt.Printf("URL: %s 具有路径或文件\n", domain)
		CfSubDomains(domain, apiKey, apiEmail, zonesid)
	}

}

func DelCfCache(domain string, apiKey string, apiEmail string, zonesid string) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", zonesid)
	method := "POST"

	// JSON 数据
	jsonData := `{"purge_everything":true}`

	// 创建一个新的请求
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 设置请求头
	req.Header.Add("X-Auth-Email", apiEmail)
	req.Header.Add("X-Auth-Key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s delete Cache: %s", domain, string(body))
}

// 子域名单独清理
func CfSubDomains(domain string, apiKey string, apiEmail string, zonesid string) {

	// 设置请求 URL
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", zonesid)

	// 构建请求体
	requestBody, err := json.Marshal(map[string][]string{
		"files": {domain},
	})
	if err != nil {
		fmt.Println("Error creating request body:", err)
		os.Exit(1)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", apiEmail)
	req.Header.Set("X-Auth-Key", apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时出错:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 打印响应
	//fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("错误: 接收到的响应状态码不是200: %d\n", resp.StatusCode)
	}

	// 解析响应体
	var purgeResponse PurgeCacheResponse
	err = json.NewDecoder(resp.Body).Decode(&purgeResponse)
	if err != nil {
		fmt.Println("解析响应时出错:", err)
	}

	// 检查清除缓存请求是否成功
	if purgeResponse.Success {
		fmt.Println("缓存清除成功。")
	} else {
		fmt.Println("缓存清除失败。")
		for _, e := range purgeResponse.Errors {
			fmt.Printf("错误: %d - %s\n", e.Code, e.Message)
		}
	}
}
