package html2text

import (
	"bytes"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Line break constants
// Deprecated: Please use HTML2TextWithOptions(text, WithUnixLineBreak())
const (
	WIN_LBR  = "\r\n"
	UNIX_LBR = "\n"
)

var (
	legacyLBR             = WIN_LBR
	badTagnamesRE         = regexp.MustCompile(`^(head|script|style)$`)
	hrefAttrRE            = regexp.MustCompile(`(?i)[ \t\n\r\f]href\s*=\s*('([^']*?)'|"([^"]*?)"|([^\s"'` + "`" + `=<>]+))`)
	headersRE             = regexp.MustCompile(`^(\/)?h[1-6]`)
	numericEntityRE       = regexp.MustCompile(`(?i)^#(x?[a-f0-9]+)$`)
	defaultAllowedSchemes = []string{"http", "https", "mailto", "tel", "sms"}
)

type options struct {
	lbr            string
	linksInnerText bool
	listPrefix     string
	keepSpaces     bool
	allowedSchemes []string
}

func newOptions() *options {
	// apply defaults
	return &options{
		lbr:        WIN_LBR,
		keepSpaces: false,
	}
}

// Option is a functional option
type Option func(*options)

// WithUnixLineBreaks instructs the converter to use unix line breaks ("\n" instead of "\r\n" default)
func WithUnixLineBreaks() Option {
	return func(o *options) {
		o.lbr = UNIX_LBR
	}
}

// WithLinksInnerText instructs the converter to retain link tag inner text and append href URLs in angle brackets after the text
// Example: click news <http://bit.ly/2n4wXRs>
func WithLinksInnerText() Option {
	return func(o *options) {
		o.linksInnerText = true
	}
}

// WithListSupportPrefix formats <ul> and <li> lists with the specified prefix
func WithListSupportPrefix(prefix string) Option {
	return func(o *options) {
		o.listPrefix = prefix
	}
}

// WithListSupport formats <ul> and <li> lists with " - " prefix
func WithListSupport() Option {
	return WithListSupportPrefix(" - ")
}

// WithKeepSpaces keep spaces as they are
func WithKeepSpaces() Option {
	return func(o *options) {
		o.keepSpaces = true
	}
}

// WithAllowedURLSchemes restricts valid URL schemes (example: []string{"http", "https", "mailto"}).
// URLs with invalid schemes are ignored.
// If this option is not used, default allowed schemes are http, https, and mailto.
func WithAllowedURLSchemes(s []string) Option {
	return func(o *options) {
		o.allowedSchemes = s
	}
}

func parseHTMLEntity(entName string) (string, bool) {
	if r, ok := entity[entName]; ok {
		return string(r), true
	}

	if match := numericEntityRE.FindStringSubmatch(entName); len(match) == 2 {
		var (
			err    error
			n      int64
			digits = match[1]
		)

		if digits != "" && (digits[0] == 'x' || digits[0] == 'X') {
			n, err = strconv.ParseInt(digits[1:], 16, 64)
		} else {
			n, err = strconv.ParseInt(digits, 10, 64)
		}

		if err == nil {
			// ignore null char, legacy override, out of range or lone surrogates
			if n == 0 || (n >= 0x80 && n <= 0x9F) || n > utf8.MaxRune || (n >= 0xD800 && n <= 0xDFFF) {
				return "", true
			}

			return string(rune(n)), true
		}
	}

	return "", false
}

func parseTagName(s string) string {
	for i, r := range s {
		switch r {
		case ' ', '\t', '\n', '\r', '\f':
			return s[:i]
		}
		if i > 0 && r == '/' { // tags like <br/>
			return s[:i]
		}
	}
	return s
}

// SetUnixLbr with argument true sets Unix-style line-breaks in output ("\n")
// with argument false sets Windows-style line-breaks in output ("\r\n", the default)
// Deprecated: Please use HTML2TextWithOptions(text, WithUnixLineBreak())
func SetUnixLbr(b bool) {
	if b {
		legacyLBR = UNIX_LBR
	} else {
		legacyLBR = WIN_LBR
	}
}

func htmlEntitiesToText(htmlEntsText string, hrefContext bool) string {
	outBuf := bytes.NewBufferString("")
	inEnt := false

	for i, r := range htmlEntsText {
		switch {
		case r == ';' && inEnt:
			inEnt = false
			continue

		case r == '&': //possible html entity
			start := i + 1
			end := start

			for end < len(htmlEntsText) && end-start < 10 && htmlEntsText[end] != ';' {
				end++
			}

			isEnt := end < len(htmlEntsText) && htmlEntsText[end] == ';'
			entName := htmlEntsText[start:end]

			if isEnt {
				if ent, isEnt := parseHTMLEntity(entName); isEnt {
					for _, er := range ent {
						if !hrefContext && isCollapsibleWhitespace(er) {
							outBuf.WriteString(" ")
						} else {
							outBuf.WriteRune(er)
						}
					}
					inEnt = true
					continue
				}
			}
		}

		if !inEnt {
			if !hrefContext && isCollapsibleWhitespace(r) {
				writeSpace(outBuf)
			} else {
				outBuf.WriteRune(r)
			}
		}
	}

	return outBuf.String()
}

// HTMLEntitiesToText decodes HTML entities inside a provided
// string and returns decoded text
func HTMLEntitiesToText(htmlEntsText string) string {
	return htmlEntitiesToText(htmlEntsText, false)
}

func writeSpace(outBuf *bytes.Buffer) {
	bts := outBuf.Bytes()
	if len(bts) > 0 && bts[len(bts)-1] != ' ' {
		outBuf.WriteString(" ")
	}
}

// isIgnoredChar returns true for control characters (except for HTML whitespace) and Unicode soft hyphen.
func isIgnoredChar(r rune) bool {
	if isHTMLWhitespace(r) {
		return false
	}
	return unicode.IsControl(r) || r == 0xad
}

// isSpace returns true for runes which should be treated as space characters.
func isSpace(r rune) bool {
	switch {
	case r == ' ', r >= 0x2008 && r <= 0x200B:
		return true
	}
	return false
}

// isHTMLWhitespace returns true for runes representing a collabsible HTML whitespace.
func isHTMLWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r', '\f':
		return true
	}
	return false
}

// isCollapsibleWhitespace returns true if the rune is a HTML whitespace or a Unicode line or paragraph separator.
func isCollapsibleWhitespace(r rune) bool {
	return isHTMLWhitespace(r) || r == 0x2028 || r == 0x2029
}

// isZeroWidth returns true for zero-width spaces and invisible format characters.
func isZeroWidth(r rune) bool {
	switch r {
	case '\u200B', '\u200C', '\u200D', '\u2060', '\uFEFF':
		return true
	default:
		return unicode.Is(unicode.Cf, r)
	}
}

func parseAndSanitizeHref(link string) string {
	decoded := strings.TrimSpace(htmlEntitiesToText(link, true))
	return strings.Map(func(r rune) rune {
		switch {
		case r == ' ':
			return r
		case r == '\\':
			return '/'
		case unicode.IsSpace(r) || isIgnoredChar(r) || isZeroWidth(r):
			return -1
		default:
			return r
		}

	}, decoded)
}

// isBadHref checks whether a link href is a valid and allowed URL.
func isBadHref(link string, allowedSchemes []string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return true
	}

	if u.Scheme == "" {
		return false
	}

	if allowedSchemes == nil {
		allowedSchemes = defaultAllowedSchemes
	}

	for _, sch := range allowedSchemes {
		if strings.EqualFold(u.Scheme, strings.TrimSuffix(sch, ":")) {
			return false
		}
	}

	return true
}

// HTML2Text converts html into a text form.
// The output is plain text, do not embed the result back into HTML without escaping.
func HTML2Text(html string) string {
	var opts []Option
	if legacyLBR == UNIX_LBR {
		opts = append(opts, WithUnixLineBreaks())
	}
	return HTML2TextWithOptions(html, opts...)
}

// HTML2TextWithOptions converts html into a text form with additional options.
// The output is plain text, do not embed the result back into HTML without escaping.
func HTML2TextWithOptions(html string, reqOpts ...Option) string {
	opts := newOptions()
	for _, opt := range reqOpts {
		opt(opts)
	}

	inLen := len(html)
	tagStart := 0
	inEnt := false
	var badTagStack []string // stack of open tag names for <head>, <script>, <style>, <a>
	shouldOutput := true
	hrefs := []string{}     // maintain a stack of sanitized <a> tag href links with html entities decoded
	tagQuoteChar := rune(0) // tracks quote context inside tags (0 = no quote, '"' or '\'' when inside quoted attr value)
	// new line cannot be printed at the beginning or
	// for <p> after a new line created by previous <p></p>
	canPrintNewline := false

	outBuf := bytes.NewBufferString("")

	for i, r := range html {
		if inLen > 0 && i == inLen-1 {
			// prevent new line at the end of the document
			canPrintNewline = false
		}

		switch {
		case isIgnoredChar(r):
			continue
		// new lines and spaces adding a single space if not there yet
		case isCollapsibleWhitespace(r):
			if shouldOutput && len(badTagStack) == 0 && !inEnt {
				if !isSpace(r) || !opts.keepSpaces {
					writeSpace(outBuf)
					continue
				}
			}

		case r == ';' && inEnt: // end of html entity
			inEnt = false
			continue

		case r == '&' && shouldOutput: // possible html entity
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
					for _, er := range ent {
						if isIgnoredChar(er) {
						} else if isCollapsibleWhitespace(er) {
							writeSpace(outBuf)
						} else {
							outBuf.WriteRune(er)
						}
					}
					inEnt = true
					continue
				}
			}

		case r == '"': // double-quote in attribute
			if !shouldOutput {
				if tagQuoteChar == '"' {
					tagQuoteChar = 0
				} else if tagQuoteChar == 0 {
					tagQuoteChar = '"'
				}
			}

		case r == '\'': // single-quote in attribute
			if !shouldOutput {
				if tagQuoteChar == '\'' {
					tagQuoteChar = 0
				} else if tagQuoteChar == 0 {
					tagQuoteChar = '\''
				}
			}

		case r == '<': // start of a tag
			if tagQuoteChar != 0 {
				continue // inside attribute quotes, not a real tag start
			}
			tagStart = i + 1
			shouldOutput = false
			continue

		case r == '>': // end of a tag
			if tagQuoteChar != 0 {
				continue // inside attribute quotes, not a real tag end
			}
			shouldOutput = true
			tag := html[tagStart:i]
			tagContentLowercase := strings.ToLower(tag)
			tagNameLowercase := parseTagName(tagContentLowercase)

			if tagNameLowercase == "/ul" || tagNameLowercase == "/ol" {
				outBuf.WriteString(opts.lbr)
			} else if tagNameLowercase == "li" || tagNameLowercase == "li/" {
				if opts.listPrefix != "" {
					outBuf.WriteString(opts.lbr)
					outBuf.WriteString(opts.listPrefix)
				} else {
					outBuf.WriteString(opts.lbr)
				}
			} else if headersRE.MatchString(tagNameLowercase) {
				if canPrintNewline {
					outBuf.WriteString(opts.lbr)
					outBuf.WriteString(opts.lbr)
				}
				canPrintNewline = false
			} else if tagNameLowercase == "br" || tagNameLowercase == "br/" {
				// new line
				outBuf.WriteString(opts.lbr)
			} else if tagNameLowercase == "p" || tagNameLowercase == "/p" {
				if canPrintNewline {
					outBuf.WriteString(opts.lbr)
					outBuf.WriteString(opts.lbr)
				}
				canPrintNewline = false
			} else if tagNameLowercase == "/a" {
				// end of link
				// links can be empty can happen if isBadHref matches the link
				if !opts.linksInnerText {
					// end of unwanted block - only pop if <a> is on top
					if len(badTagStack) > 0 && badTagStack[len(badTagStack)-1] == "a" {
						badTagStack = badTagStack[:len(badTagStack)-1]
					}
				}

				if len(hrefs) > 0 {
					last := len(hrefs) - 1
					if hrefs[last] != "" {
						if opts.linksInnerText {
							outBuf.WriteString(" <")
							outBuf.WriteString(hrefs[last])
							outBuf.WriteString(">")
						} else {
							outBuf.WriteString(hrefs[last])
						}
					}
					hrefs = hrefs[:last]
				}
			} else if tagNameLowercase == "a" {
				// start of link
				if !opts.linksInnerText {
					// unwanted block
					badTagStack = append(badTagStack, "a")
				}

				// parse link href
				m := hrefAttrRE.FindStringSubmatch(tag)
				if len(m) == 5 {
					link := m[2]
					if len(link) == 0 {
						link = m[3]
						if len(link) == 0 {
							link = m[4]
						}
					}

					link = parseAndSanitizeHref(link)

					if !isBadHref(link, opts.allowedSchemes) {
						hrefs = append(hrefs, link)
					} else {
						hrefs = append(hrefs, "") // maintain stack alignment
					}
				} else {
					hrefs = append(hrefs, "") // maintain stack alignment
				}
			} else if badTagnamesRE.MatchString(tagNameLowercase) {
				// unwanted block
				badTagStack = append(badTagStack, tagNameLowercase)
			} else if len(tagNameLowercase) > 0 && tagNameLowercase[0] == '/' &&
				badTagnamesRE.MatchString(tagNameLowercase[1:]) {
				// end of unwanted block - only pop if matching tag is on top
				openingTag := tagNameLowercase[1:]
				if len(badTagStack) > 0 && badTagStack[len(badTagStack)-1] == openingTag {
					badTagStack = badTagStack[:len(badTagStack)-1]
				}
			}
			continue

		} // switch end

		if shouldOutput && len(badTagStack) == 0 && !inEnt {
			canPrintNewline = true
			outBuf.WriteRune(r)
		}
	}

	return outBuf.String()
}
