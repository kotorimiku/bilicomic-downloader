package bilicomicdownloader

import (
	"bytes"
	"io"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client = resty.New()

func ClientInit(cookie string) {
	headers := map[string]string{
		"Accept":                      "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Encoding":             "gzip, deflate",
		"Accept-Language":             "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7",
		"Cache-Control":               "max-age=0",
		"Content-Type":                "application/x-www-form-urlencoded",
		"Cookie":                      cookie,
		"Origin":                      "https://www.bilicomic.net",
		"Priority":                    "u=0, i",
		"Referer":                     "https://www.bilicomic.net/read/78/4616.html?__cf_chl_tk=xbQ0jTrU2wgptFcOGBeXnbEDCainfHlCIEIWb145Ys8-1740640867-1.0.1.1-BxummUHvbxhg9MNBrbr9QeBnZK1rrA3zJgYbUHN5Ex4",
		"Sec-Ch-Ua":                   `"Not A;Brand";v="99", "Chromium";v="133", "Google Chrome";v="133"`,
		"Sec-Ch-Ua-Arch":              `""`,
		"Sec-Ch-Ua-Bitness":           `"64"`,
		"Sec-Ch-Ua-Full-Version":      `"133.0.0.0"`,
		"Sec-Ch-Ua-Full-Version-List": `"Not(A:Brand";v="99.0.0.0", "Google Chrome";v="133.0.6943.128", "Chromium";v="133.0.6943.128"`,
		"Sec-Ch-Ua-Mobile":            `?1`,
		"Sec-Ch-Ua-Model":             `"Nexus 5"`,
		"Sec-Ch-Ua-Platform":          `"Android"`,
		"Sec-Ch-Ua-Platform-Version":  `"6.0"`,
		"Sec-Fetch-Dest":              "document",
		"Sec-Fetch-Mode":              "navigate",
		"Sec-Fetch-Site":              "same-origin",
		"Sec-Fetch-User":              `?1`,
		"Upgrade-Insecure-Requests":   "1",
		"User-Agent":                  "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Mobile Safari/537.36",
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

func GetImageRaw(url string) (io.ReadCloser, error) {
	resp, err := client.R().SetHeader("Accept", "image/webp,image/apng,image/*,*/*;q=0.8").Get(url)
	if err != nil {
		return nil, err
	}

	return resp.RawBody(), nil
}
