/*
Package internal
Copyright © 2022 Pavel Sidlo <github.com/dyamon-cz>
*/
package internal

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/term"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Module struct {
	Name        string
	Path        string
	Description string
	Imports     string
	Version     string
	Published   string
	Licence     string
}

const uri = "https://pkg.go.dev/search?limit=100&m=package&q="

func SearchModules(pkg string) []Module {
	res, err := http.Get(uri + url.QueryEscape(pkg))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var modules = make([]Module, 0)

	doc.Find(".SearchSnippet").Each(func(i int, s *goquery.Selection) {

		module := Module{}

		name, path := parseNameAndPath(s)
		module.Name = name
		module.Path = path

		// there is a bug in prompt UI, if description overflow single line
		// check what is the width of terminal and shorten the description
		description := parseDescription(s)
		width, _, _ := term.GetSize(0)
		module.Description = ellipsis(description, width-4) // 4 is the length of ellipsis

		module.Imports = parseImports(s)
		module.Version = parseVersion(s)
		module.Published = parsePublished(s)
		module.Licence = parseLicence(s)

		modules = append(modules, module)
	})

	return modules
}

func parseNameAndPath(s *goquery.Selection) (string, string) {
	innerText := s.Find("a[data-gtmc='search result']").First().Text()

	re := regexp.MustCompile(`\r?\n\s+`)
	namePath := re.ReplaceAllString(innerText, " ")

	namePath = strings.TrimSpace(namePath)
	nameAndPath := strings.Split(namePath, " ")

	// remove brackets
	nameAndPath[1] = strings.Replace(nameAndPath[1], "(", "", -1)
	nameAndPath[1] = strings.Replace(nameAndPath[1], ")", "", -1)

	return nameAndPath[0], nameAndPath[1]
}

func parseDescription(s *goquery.Selection) string {
	innerText := s.Find(".SearchSnippet-synopsis").First().Text()

	description := strings.TrimSpace(innerText)

	return description
}

func parseImports(s *goquery.Selection) string {
	innerText := s.Find(".SearchSnippet-infoLabel a strong").First().Text()

	imports := strings.TrimSpace(innerText)

	return imports
}

func parseVersion(s *goquery.Selection) string {
	innerText := s.Find(".SearchSnippet-infoLabel .go-textSubtle strong").First().Text()

	version := strings.TrimSpace(innerText)

	return version
}

func parsePublished(s *goquery.Selection) string {
	innerText := s.Find(".SearchSnippet-infoLabel .go-textSubtle span[data-test-id='snippet-published'] strong").First().Text()

	published := strings.TrimSpace(innerText)

	return published
}

func parseLicence(s *goquery.Selection) string {
	innerText := s.Find(".SearchSnippet-infoLabel span[data-test-id='snippet-license'] a").First().Text()

	licence := strings.TrimSpace(innerText)

	return licence
}
