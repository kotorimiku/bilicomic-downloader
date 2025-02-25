package bilicomicdownloader

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type ComicInfo struct {
	Series    string `xml:"Series"`
	Writer    string `xml:"Writer"`
	Publisher string `xml:"Publisher,omitempty"`
	Genre     string `xml:"Genre"`
	Summary   string `xml:"Summary"`
	Title     string `xml:"Title"`
	Number    string `xml:"Number,omitempty"`
	Volume    string `xml:"Volume,omitempty"`
	PageCount string `xml:"PageCount"`
}

// Build 方法将 ComicInfo 对象序列化并保存为 XML 文件
func (c *ComicInfo) Build(outputPath string) error {
	fileName := filepath.Join(outputPath, "ComicInfo.xml")

	// 创建文件
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// 创建 XML 编码器并序列化 ComicInfo
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ") // 设置缩进
	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("error serializing to XML: %v", err)
	}

	return nil
}
