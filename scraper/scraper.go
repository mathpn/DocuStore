package scraper

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"golang.org/x/net/html"
)

var URLRegex = regexp.MustCompile(`^htt(p|ps)://(.*)(\s|$)`)

type ScrapeData struct {
	Title   string
	Content string
}

func ScrapeText(url string, log logger.Logger) (*ScrapeData, error) {
	buffer := bytes.NewBufferString("")
	response, err := http.Get(strings.TrimSpace(url))
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resBody))
	if err != nil {
		return nil, err
	}

	title := doc.Find("title").Text()
	buffer.WriteString(title + "\n")
	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); strings.ToLower(name) == "description" {
			description, _ := s.Attr("content")
			buffer.WriteString(description + "\n")
		}
	})

	textTags := []string{
		"a",
		"p",
		"strong",
		"code",
		"span",
		// "em",
		// "string",
		// "blockquote",
		// "q",
		// "cite",
		"h1",
		"h2",
		"h3",
		"h4",
		"h5",
		"h6",
	}

	tag := ""
	enter := false

	tokenizer := html.NewTokenizer(bytes.NewReader(resBody))
	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		tokenString := token.String()
		if strings.HasPrefix(tokenString, "<footer") {
			break
		}

		switch tt {
		case html.ErrorToken:
			log.Error(fmt.Sprintf("bad HTML token: %s", token))
		case html.StartTagToken, html.SelfClosingTagToken:
			enter = false

			tag = token.Data
			for _, ttt := range textTags {
				if tag == ttt {
					enter = true
					buffer.WriteString("\n")
					break
				}
			}
		case html.TextToken:
			if enter {
				data := strings.TrimSpace(token.Data)

				if len(data) > 0 {
					data = URLRegex.ReplaceAllString(data, "")
					buffer.WriteString(data + " ")
				}
			}
		}
	}
	result := &ScrapeData{Title: title, Content: buffer.String()}
	return result, nil
}

// func main() {
// 	out := ScrapeText("https://spark.apache.org/docs/latest/")
// 	fmt.Println("--------------")
// 	fmt.Print(out)
// }
