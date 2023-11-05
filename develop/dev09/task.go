package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

/*
=== Утилита wget ===

# Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flags struct {
	d int
}

func main() {
	f := &flags{}
	flag.IntVar(&f.d, "d", 1, "Depth for recursive download")
	flag.Parse()
	nonFlagArgs := flag.Args()
	if len(nonFlagArgs) == 0 {
		log.Fatalln("No url specified")
	}
	urlInput := nonFlagArgs[0]
	wget(urlInput, f.d)
}

// скачать страницу
func getUrlPage(uri string) []byte {
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Fprintln(os.Stderr, "http err: ", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error status: %s\n", resp.Status)
		return nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading page: %v\n", err)
		return nil
	}
	return b
}

// запуск wget
func wget(url string, depth int) {
	fmt.Println("downloading " + url)
	if depth < 1 {
		fmt.Fprintln(os.Stderr, "wrong depth")
		return
	}
	trimmed := strings.TrimLeft(url, `https://`)
	splitted := strings.Split(trimmed, "/")
	baseUrl := splitted[0]
	wgetRec(depth, url, baseUrl, "", 0)
}

// рекурсивный wget
func wgetRec(depth int, url, baseUrl, oldPath string, recDepth int) {
	if recDepth == depth {
		return
	}
	page := getUrlPage(url)
	if page != nil {
		links := parseLinks(page, baseUrl)
		//находим ссылки, проходим по ним и меняем на соответствующие директории
		for _, val := range links {
			byteVal := []byte(val)
			newLink := "./" + oldPath + linkToFilePath(val)
			page = bytes.ReplaceAll(page, byteVal, []byte(newLink))
		}
		//записываем страницу
		err := writeToFile(page, oldPath+linkToFilePath(url)+"/")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		for _, v := range links {
			link := `https://` + linkToFilePath(v)
			splitted := strings.Split(v, "/")
			baseUrl = splitted[0]
			wgetRec(depth, link, baseUrl, linkToFilePath(url)+`/`, recDepth+1)
		}
	}
}

func writeToFile(data []byte, dirname string) error {
	if err := os.MkdirAll("./"+dirname, 0755); err != nil {
		return err
	}
	f, err := os.Create("./" + dirname + "index.html")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func parseLinks(data []byte, baseUrl string) []string {
	links := make([]string, 0)
	body := bytes.NewReader(data)
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						if strings.HasPrefix(attr.Val, "/") {
							attr.Val = baseUrl + attr.Val
						}
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
	return links
}

func linkToFilePath(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		fmt.Fprintln(os.Stderr, "problem while parsing url")
		return ""
	}
	p := u.Hostname() + "/" + u.EscapedPath()
	fullPath := strings.Split(p, "/")
	if len(fullPath[len(fullPath)-1]) == 0 {
		fullPath = fullPath[:len(fullPath)-1]
	}
	return path.Join(fullPath...)
}
