package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

const (
	blogName = "Jack Baker | Developer"
)

func main() {
	// read header
	headerBytes, err := os.ReadFile("res/header.html")
	if err != nil {
		log.Fatal(err.Error())
	}
	header := string(headerBytes)
	header += wrapInTagInline("h1", blogName)
	// read side column
	sideColumnBytes, err := os.ReadFile("res/sideColumn.html")
	if err != nil {
		log.Fatal(err.Error())
	}
	// some dummy footer
	footer := "\n</html>"

	// create list of all posts
	allPosts := make([]post, 0)
	_ = os.Mkdir("blog", os.ModePerm)

	// for every file in "posts" directory
	files, _ := os.ReadDir("posts")
	for _, postfile := range files {
		samplefile, err := os.Open(path.Join("posts", postfile.Name()))
		if err != nil {
			log.Fatal(err.Error())
		}
		filescanner := bufio.NewScanner(samplefile)

		newPost := post{}
		postBody := "\n"

		// first line post title
		filescanner.Scan()
		postBody += wrapInTagInline("h2", filescanner.Text())
		newPost.title = filescanner.Text()

		// second line date
		filescanner.Scan()
		date, err := time.Parse("1/2/06", filescanner.Text())
		if err != nil {
			log.Fatal(err.Error())
		}
		newPost.date = date
		postBody += wrapInTagInline("i", date.Format("January 2, 2006"))

		// third line blank for future extensbility (e.g. tags)
		filescanner.Scan()

		formatter := formatter{}

		// read rest of file and do formatting
		for filescanner.Scan() {
			line, err := formatter.formatLine(filescanner.Text())
			if err != nil {
				log.Fatal("error in ", postfile.Name(), err)
			}
			postBody += "\n" + line
			newPost.body += "\n" + line
		}

		// wrap it in a big div with these css classes
		postBody = wrapInDiv("column three-quarters", postBody)
		postBody += "\n" + string(sideColumnBytes)

		// put it all in the <body> tag
		postBody = wrapInTag("body", postBody)
		allPosts = append(allPosts, newPost)

		// assemble it all and write it out
		page := header + "\n" + postBody + "\n" + footer
		err = os.WriteFile(newPost.url(), []byte(page), os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println(newPost.title, "written to", newPost.url())
	}

	// sort all posts reverse chronological
	slices.SortFunc(allPosts, func(a post, b post) int {
		return -a.date.Compare(b.date)
	})

	// build the index
	indexBody := ""

	// starting with the post list
	for _, pst := range allPosts {
		// link the url with the title
		postLink := "<a href=\"/" + pst.url() + "\">" + pst.title + "</a>"
		// format the post date nicely
		postDate := wrapInTagInline("i", pst.date.Format("01-02-2006"))
		// combine and wrap as a list item
		indexBody += "\n" + wrapInTagInline("li", postLink+" - "+postDate)
	}

	// wrap the list
	indexBody = "\n" + wrapInTag("ul", indexBody)

	// add header
	indexBody = "\n" + wrapInTagInline("h2", "Archive of Posts") + indexBody
	// wrap in a div and add the side column
	indexBody = "\n" + wrapInDiv("column three-quarters", indexBody)
	indexBody += "\n" + string(sideColumnBytes)

	// add header and footer
	indexBody = wrapInTag("body", indexBody)
	index := header + "\n" + indexBody + "\n" + footer

	err = os.WriteFile("blog/index.html", []byte(index), os.ModePerm)

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(len(allPosts), "posts generated")
}

func wrapInTag(tag, str string) string {
	return fmt.Sprint("<", tag, ">", str, "\n</", tag, ">")
}

func wrapInTagInline(tag, str string) string {
	return fmt.Sprint("<", tag, ">", str, "</", tag, ">")
}

func wrapInDiv(class, str string) string {
	return fmt.Sprint("<div class=\"", class, "\">", str, "\n</div>")
}

type post struct {
	title string
	body  string
	date  time.Time
}

func (p post) url() string {
	return "blog/" + strings.ReplaceAll(strings.ToLower(p.title), " ", "-") + ".html"
}
