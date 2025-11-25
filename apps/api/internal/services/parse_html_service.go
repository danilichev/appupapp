package services

import (
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"

	"apps/api/internal/api"
)

type ParseHtmlService struct{}

func NewParseHtmlService() *ParseHtmlService {
	return &ParseHtmlService{}
}

func (s *ParseHtmlService) ParseHtml(
	url string,
) (*api.ParsedHtml, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body := io.LimitReader(resp.Body, 50*1024)
	z := html.NewTokenizer(body)

	inHead := false

	var title string

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		tok := z.Token()

		switch tt {
		case html.StartTagToken:
			if tok.Data == "head" {
				inHead = true
			}
			if tok.Data == "title" && title == "" {
				if z.Next() == html.TextToken {
					title = strings.TrimSpace(z.Token().Data)
				}
			}
		case html.EndTagToken:
			if tok.Data == "head" && inHead {
				break
			}
		}

		if title != "" {
			break
		}
	}

	return &api.ParsedHtml{
		Title: &title,
	}, nil
}
