package bilicomicdownloader

import (
	"bytes"
	"io"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client = resty.New()

func ClientInit(cookie string) {
	var headers = map[string]string{
		"User-Agent":                "Mozilla/5.0 (Linux; Android 11; M2102J20SG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Mobile Safari/537.36 EdgA/97.0.1072.78",
		"Accept-Language":           "en,zh-HK;q=0.9,zh-TW;q=0.8,zh-CN;q=0.7,zh;q=0.6,en-GB;q=0.5,en-US;q=0.4",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Cookie":                    cookie,
		"Referer":                   "www.bilicomic.net",
		"Accept-Encoding":           "gzip, deflate",
		"Priority":                  "u=0, i",
		"Sec-CH-UA":                 `"Microsoft Edge";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
		"Sec-CH-UA-Mobile":          "?1",
		"Sec-CH-UA-Platform":        `"Android"`,
		"Sec-Fetch-Dest":            "document",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-User":            "?1",
		"Sec-Fetch-Site":            "same-origin",
		"Sec-Fetch-Mode":            "navigate",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
	}

	client.SetHeaders(headers)
}

func GetText(url string) (string, error) {
	resp, err := client.R().Get(url)
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}

func GetBytes(url string) ([]byte, error) {
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func GetRaw(url string) (io.Reader, error) {
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(resp.Body()), nil
}

func GetImage(url string) ([]byte, error) {
	resp, err := client.R().SetHeader("Accept", "image/webp,image/apng,image/*,*/*;q=0.8").Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
