package librecursebuster

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pmezard/go-difflib/difflib"
)

//RandString will return a UUID
func RandString(printChan chan OutLine) string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

}

//returns a slice of strings containing urls
func getUrls(page []byte, printChan chan OutLine) ([]string, error) {

	ret := []string{}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		return nil, err
	}

	doc.Find("*").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		if len(link) > 0 {
			ret = append(ret, link)
		}
	})

	return ret, nil
}

func detectSoft404(a []byte, b []byte, ratio float64) bool {
	diff := difflib.SequenceMatcher{}
	diff.SetSeqs(strings.Split(string(a), " "), strings.Split(string(b), " "))

	if diff.Ratio() > ratio {
		return true
	}
	return false
}

func cleanURL(u *url.URL, actualURL string) string {
	var didHaveSlash bool
	if len(u.Path) > 0 {
		didHaveSlash = string(u.Path[len(u.Path)-1]) == "/"
		if string(u.Path[0]) != "/" {
			u.Path = "/" + u.Path
		}
	}

	cleaned := path.Clean(u.Path)

	if string(cleaned[0]) != "/" {
		cleaned = "/" + cleaned
	}
	if cleaned != "." {
		actualURL += cleaned
	}

	if didHaveSlash && cleaned != "/" {
		actualURL += "/"
	}
	return actualURL
}
