package main

import (
	"errors"
	"regexp"
	"strings"
)

type formatter struct {
	isInCodeBlock bool
	indentLevel   int
	addLines      []string
}

func (f *formatter) formatLine(line string) (string, error) {
	output := ""
	var err error
	if len(line) == 0 {
		return output + "<br>", nil
	}

	// multiline code block handling
	if line == "```" {
		f.isInCodeBlock = !f.isInCodeBlock
		if f.isInCodeBlock {
			line = "<pre>"
		} else {
			line = "</pre>"
		}
	}

	output, err = f.checkListFormatting(line)
	if err != nil {
		return "", err
	}

	// no format if in code block
	if f.isInCodeBlock {
		return f.addExtraLines(line), nil
	}

	// order matters
	if output, err = applyMarkdown(output); err != nil {
		return "", err
	}
	output = replaceLinks(output)
	output = replaceImages(output)

	if f.indentLevel == 0 { // only top-level blocks have <p> tags
		output = wrapInTagInline("p", output)
	} else { // otherwise it's a list item
		output = wrapInTagInline("li", output)
	}

	return f.addExtraLines(output), nil
}

func (f *formatter) addExtraLines(inStr string) string {
	for _, addLine := range f.addLines {
		inStr = addLine + "\n" + inStr
	}
	f.addLines = make([]string, 0)

	return inStr
}

func (f *formatter) checkListFormatting(str string) (string, error) {
	desiredIndent := 0
	for _, char := range str {
		if char == '>' {
			desiredIndent++
		} else {
			break
		}
	}

	if desiredIndent > f.indentLevel {
		for range desiredIndent - f.indentLevel {
			f.addLines = append(f.addLines, "<ul>")
		}
	}

	if desiredIndent < f.indentLevel {
		for range f.indentLevel - desiredIndent {
			f.addLines = append(f.addLines, "</ul>")
		}
	}

	f.indentLevel = desiredIndent
	return str[desiredIndent:], nil
}

func applyMarkdown(str string) (string, error) {
	var (
		replaceMap = map[string]string{
			"**":  "b",
			"__":  "i",
			"--":  "s",
			"```": "code",
		}
	)
	for markdown, html := range replaceMap {
		if strings.Count(str, markdown)%2 == 0 {
			for strings.Count(str, markdown) != 0 {
				str = strings.Replace(str, markdown, "<"+html+">", 1)
				str = strings.Replace(str, markdown, "</"+html+">", 1)
			}
		} else {
			return "", errors.New("mismatched markdown formatting")
		}
	}
	return str, nil
}

func replaceImages(str string) string {
	regex := regexp.MustCompile(`\[(.*?)\]`)
	str = regex.ReplaceAllString(str, `<img src="$1"/>`)
	return str
}

func replaceLinks(str string) string {
	regex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	str = regex.ReplaceAllString(str, `<a href="$2">$1</a>`)
	return str
}
