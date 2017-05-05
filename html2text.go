package html2text

import (
	"bytes"
	"regexp"
	"strings"
)

var badTagnamesRE = regexp.MustCompile(`^(head|script|style|a)($|\s*)`)
var linkTagRE = regexp.MustCompile(`a.*href=('([^']*?)'|"([^"]*?)")`)
var badLinkHrefRE = regexp.MustCompile(`#|javascript:`)

func parseHTMLEntity(entName string) (string, bool) {
	entName = strings.ToLower(entName)

	// possible entities
	switch entName {
	case "nbsp":
		return " ", true
	case "gt":
		return ">", true
	case "lt":
		return "<", true
	case "amp":
		return "&", true
	case "quot":
		return "\"", true
	case "apos":
		return "'", true
	case "cent":
		return "¢", true
	case "pound":
		return "£", true
	case "yen":
		return "¥", true
	case "euro":
		return "€", true
	case "copy":
		return "©", true
	case "reg":
		return "®", true
	default:
		return "", false
	}

}

// HTMLEntitiesToText decodes HTML entities inside a provided
// string and returns decoded text
func HTMLEntitiesToText(htmlEntsText string) string {
	outBuf := bytes.NewBufferString("")
	inEnt := false

	for i, r := range htmlEntsText {
		switch {
		case r == ';' && inEnt:
			inEnt = false
			continue

		case r == '&': //possible html entity
			entName := ""
			isEnt := false

			// parse the entity name - max 10 chars
			chars := 0
			for _, er := range htmlEntsText[i+1:] {
				if er == ';' {
					isEnt = true
					break
				} else {
					entName += string(er)
				}

				chars++
				if chars == 10 {
					break
				}
			}

			if isEnt {
				if ent, isEnt := parseHTMLEntity(entName); isEnt {
					outBuf.WriteString(ent)
					inEnt = true
					continue
				}
			}
		}

		if !inEnt {
			outBuf.WriteRune(r)
		}
	}

	return outBuf.String()
}

// HTML2Text converts html into a text form
func HTML2Text(html string) string {
	tagStart := 0
	inEnt := false
	badTagStackDepth := 0 // if == 1 it means we are inside <head>...</head>
	shouldOutput := true

	outBuf := bytes.NewBufferString("")

	for i, r := range html {
		switch {
		case r == ';' && inEnt:
			inEnt = false
			shouldOutput = true
			continue

		case r == '&' && shouldOutput: //possible html entity
			entName := ""
			isEnt := false

			// parse the entity name - max 10 chars
			chars := 0
			for _, er := range html[i+1:] {
				if er == ';' {
					isEnt = true
					break
				} else {
					entName += string(er)
				}

				chars++
				if chars == 10 {
					break
				}
			}

			if isEnt {
				if ent, isEnt := parseHTMLEntity(entName); isEnt {
					outBuf.WriteString(ent)
					inEnt = true
					shouldOutput = false
					continue
				}
			}

		case r == '<':
			tagStart = i + 1
			shouldOutput = false
			continue
		case r == '>':
			shouldOutput = true
			tagName := strings.ToLower(html[tagStart:i])

			if badTagnamesRE.MatchString(tagName) {
				badTagStackDepth++

				// parse link href
				m := linkTagRE.FindStringSubmatch(tagName)
				if len(m) == 4 {
					link := m[2]
					if len(link) == 0 {
						link = m[3]
					}

					if !badLinkHrefRE.MatchString(link) {
						outBuf.WriteString(HTMLEntitiesToText(link))
					}
				}

			} else if len(tagName) > 0 && tagName[0] == '/' &&
				badTagnamesRE.MatchString(tagName[1:]) {
				badTagStackDepth--
			}

			continue
		}

		if shouldOutput && badTagStackDepth == 0 && !inEnt {
			outBuf.WriteRune(r)
		}
	}

	return outBuf.String()
}
