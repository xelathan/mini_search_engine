package search

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type CrawlData struct {
	Url          string
	Success      bool
	ResponseCode int
	CrawlData    ParsedBody
}

type ParsedBody struct {
	CrawlTime       time.Duration
	PageTitle       string
	PageDescription string
	Headings        string
	Links           Links
}

type Links struct {
	Internal []string
	External []string
}

func runCrawl(inputUrl string) CrawlData {
	resp, err := http.Get(inputUrl)
	baseUrl, _ := url.Parse(inputUrl)

	if err != nil || resp == nil {
		fmt.Println("Error getting response")
		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: 0,
			CrawlData:    ParsedBody{},
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Non 200 response")
		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: resp.StatusCode,
			CrawlData:    ParsedBody{},
		}
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		data, err := parseBody(resp.Body, baseUrl)
		if err != nil {
			fmt.Println("Error parsing body")
			return CrawlData{
				Url:          inputUrl,
				Success:      false,
				ResponseCode: resp.StatusCode,
				CrawlData:    ParsedBody{},
			}
		}
		return CrawlData{
			Url:          inputUrl,
			Success:      true,
			ResponseCode: resp.StatusCode,
			CrawlData:    data,
		}
	} else {
		fmt.Println("Non HTML response")
		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: resp.StatusCode,
			CrawlData:    ParsedBody{},
		}
	}
}

func getLinks(node *html.Node, baseUrl *url.URL) Links {
	links := Links{}
	if node == nil {
		return links
	}

	var findLinks func(node *html.Node)
	findLinks = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					linkUrl, err := url.Parse(attr.Val)
					if err != nil || strings.HasPrefix(linkUrl.String(), "#") ||
						strings.HasPrefix(linkUrl.String(), "#mail") ||
						strings.HasPrefix(linkUrl.String(), "tel") ||
						strings.HasPrefix(linkUrl.String(), "javascript") ||
						strings.HasPrefix(linkUrl.String(), ".pdf") ||
						strings.HasPrefix(linkUrl.String(), ".md") {
						continue
					}

					if linkUrl.IsAbs() {
						if isSameHost(linkUrl.String(), baseUrl.String()) {
							links.Internal = append(links.Internal, linkUrl.String())
						} else {
							links.External = append(links.External, linkUrl.String())
						}
					} else {
						rel := baseUrl.ResolveReference(linkUrl)
						links.Internal = append(links.Internal, rel.String())
					}
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findLinks(child)
		}
	}

	findLinks(node)
	return links
}

func isSameHost(absoluteUrl, base string) bool {
	absUrl, err := url.Parse(absoluteUrl)
	if err != nil {
		return false
	}

	baseUrl, err2 := url.Parse(base)
	if err2 != nil {
		return false
	}

	return absUrl.Host == baseUrl.Host
}

func parseBody(body io.Reader, baseUrl *url.URL) (ParsedBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return ParsedBody{}, err
	}
	start := time.Now()

	links := getLinks(doc, baseUrl)

	title, description := getPageData(doc)

	headings := getHeadings(doc)

	end := time.Now()
	return ParsedBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: description,
		Headings:        headings,
		Links:           links,
	}, nil
}

func getHeadings(n *html.Node) string {
	if n == nil {
		return ""
	}

	var heading strings.Builder
	var findH1 func(n *html.Node)

	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			if n.FirstChild != nil {
				heading.WriteString(n.FirstChild.Data)
				heading.WriteString(", ")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findH1(c)
		}
	}

	return strings.TrimSuffix(heading.String(), ",")
}

func getPageData(node *html.Node) (string, string) {
	if node == nil {
		return "", ""
	}

	title, desc := "", ""
	var findMetaAndTitle func(n *html.Node)
	findMetaAndTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			} else {
				title = ""
			}
		} else if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					name = attr.Val
				} else if attr.Key == "content" {
					content = attr.Val
				}
			}
			if name == "description" {
				desc = content
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findMetaAndTitle(child)
	}

	return title, desc
}
