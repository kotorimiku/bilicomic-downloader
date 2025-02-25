package bilicomicdownloader

type BookInfo struct {
	Author      []string
	Description string
	Genre       []string
	Title       string
	Cover       string
}

type Volume struct {
	Title    string
	Cover    string
	Chapters []*Chapter
}

type Chapter struct {
	Title    string
	Url      string
	progress float32
}
