package html2text

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHTML2Text(t *testing.T) {
	Convey("HTML2Text should work", t, func() {

		Convey("Links", func() {
			So(HTML2Text(`<div></div>`), ShouldEqual, "")
			So(HTML2Text(`<div>simple text</div>`), ShouldEqual, "simple text")

			// the original behavior
			So(HTML2Text(`click <a href="test">here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <A hRef="test">here</A>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href='test'>here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href=test>here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href =test>here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href = test>here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href      =     test>here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href      =     test target="_blank">here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a class="x" href="test">here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href="javascript:void(0)">here</a>`), ShouldEqual, "click ")
			So(HTML2Text(`click <a href="JaVaScRiPt:alert(1)">here</a>`), ShouldEqual, "click ")
			So(HTML2Text(`click <a href="&#106;avascript:alert(1)">here</a>`), ShouldEqual, "click ")
			So(HTML2Text(`click <a href="test"><span>here</span> or here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href="http://bit.ly/2n4wXRs">news</a>`), ShouldEqual, "click http://bit.ly/2n4wXRs")
			So(HTML2Text(`<a rel="mw:WikiLink" href="/wiki/yet#English" title="yet">yet</a>, <a rel="mw:WikiLink" href="/wiki/not_yet#English" title="not yet">not yet</a>`), ShouldEqual, "/wiki/yet#English, /wiki/not_yet#English")

			// with inner text
			So(HTML2TextWithOptions(`click <a href="test">here</a>`, WithLinksInnerText()), ShouldEqual, "click here <test>")
			So(HTML2TextWithOptions(`click <A hRef="test">here</A>`, WithLinksInnerText()), ShouldEqual, "click here <test>")
			So(HTML2TextWithOptions(`click <a href='test'>here</a>`, WithLinksInnerText()), ShouldEqual, "click here <test>")
			So(HTML2TextWithOptions(`click <a href=test>here</a>`, WithLinksInnerText()), ShouldEqual, "click here <test>")
			So(HTML2TextWithOptions(`click <a href =test>here</a>`, WithLinksInnerText()), ShouldEqual, "click here <test>")
			So(HTML2TextWithOptions(`click <a href = test>here</a>`, WithLinksInnerText()), ShouldEqual, "click here <test>")
			So(HTML2TextWithOptions(`click <a href      =     test>here</a>`, WithLinksInnerText()), ShouldEqual, "click here <test>")
			So(HTML2TextWithOptions(`click <a href      =     test target="_blank">here</a>`, WithLinksInnerText()), ShouldEqual, "click here <test>")
			So(HTML2TextWithOptions(`click <a class="x" href="test">here</a>`, WithLinksInnerText()), ShouldEqual, "click here <test>")
			So(HTML2TextWithOptions(`click <a href="ents/&apos;x&apos;">here</a>`, WithLinksInnerText()), ShouldEqual, "click here <ents/'x'>")
			So(HTML2TextWithOptions(`click <a href="javascript:void(0)">here</a>`, WithLinksInnerText()), ShouldEqual, "click here")
			So(HTML2TextWithOptions(`click <a href="JaVaScRiPt:alert(1)">here</a>`, WithLinksInnerText()), ShouldEqual, "click here")
			So(HTML2TextWithOptions(`click <a href="&#106;avascript:alert(1)">here</a>`, WithLinksInnerText()), ShouldEqual, "click here")
			So(HTML2TextWithOptions(`click <a href="test"><span>here</span> or here</a>`, WithLinksInnerText()), ShouldEqual, "click here or here <test>")
			So(HTML2TextWithOptions(`click <a href="http://bit.ly/2n4wXRs">news</a>`, WithLinksInnerText()), ShouldEqual, "click news <http://bit.ly/2n4wXRs>")
			So(HTML2TextWithOptions(`<a rel="mw:WikiLink" href="/wiki/yet#English" title="yet">yet</a>, <a rel="mw:WikiLink" href="/wiki/not_yet#English" title="not yet">not yet</a>`, WithLinksInnerText()), ShouldEqual, "yet </wiki/yet#English>, not yet </wiki/not_yet#English>")
			So(HTML2TextWithOptions(`click <a href="one">here<a href="two"> or</a><span> here</span></a>`, WithLinksInnerText()), ShouldEqual, "click here or <two> here <one>")
		})

		Convey("Inlines", func() {
			So(HTML2Text(`strong <strong>text</strong>`), ShouldEqual, "strong text")
			So(HTML2Text(`some <div id="a" class="b">div</div>`), ShouldEqual, "some div")
		})

		Convey("Line breaks and spaces", func() {
			So(HTML2Text("should    ignore more spaces"), ShouldEqual, "should ignore more spaces")
			So(HTML2Text("should \nignore \r\nnew lines"), ShouldEqual, "should ignore new lines")
			So(HTML2Text("a\nb\nc"), ShouldEqual, "a b c")
			So(HTML2Text(`two<br>line<br/>breaks`), ShouldEqual, "two\r\nline\r\nbreaks")
			So(HTML2Text(`<p>two</p><p>paragraphs</p>`), ShouldEqual, "two\r\n\r\nparagraphs")
		})

		Convey("Headings", func() {
			So(HTML2Text("<h1>First</h1>main text"), ShouldEqual, "First\r\n\r\nmain text")
			So(HTML2Text("First<h2>Second</h2>next section"), ShouldEqual, "First\r\n\r\nSecond\r\n\r\nnext section")
			So(HTML2Text("<h2>Second</h2>next section"), ShouldEqual, "Second\r\n\r\nnext section")
			So(HTML2Text("Second<h3>Third</h3>next section"), ShouldEqual, "Second\r\n\r\nThird\r\n\r\nnext section")
			So(HTML2Text("<h3>Third</h3>next section"), ShouldEqual, "Third\r\n\r\nnext section")
			So(HTML2Text("Third<h4>Fourth</h4>next section"), ShouldEqual, "Third\r\n\r\nFourth\r\n\r\nnext section")
			So(HTML2Text("<h4>Fourth</h4>next section"), ShouldEqual, "Fourth\r\n\r\nnext section")
			So(HTML2Text("Fourth<h5>Fifth</h5>next section"), ShouldEqual, "Fourth\r\n\r\nFifth\r\n\r\nnext section")
			So(HTML2Text("<h5>Fifth</h5>next section"), ShouldEqual, "Fifth\r\n\r\nnext section")
			So(HTML2Text("Fifth<h6>Sixth</h6>next section"), ShouldEqual, "Fifth\r\n\r\nSixth\r\n\r\nnext section")
			So(HTML2Text("<h6>Sixth</h6>next section"), ShouldEqual, "Sixth\r\n\r\nnext section")
			So(HTML2Text("<h7>Not Header</h7>next section"), ShouldEqual, "Not Headernext section")
		})

		Convey("HTML entities", func() {
			So(HTMLEntitiesToText("collapsible&#x20;  and explicit&nbsp;&nbsp;spaces"), ShouldEqual, "collapsible and explicit  spaces")
			So(HTML2Text(`two&nbsp;&nbsp;spaces`), ShouldEqual, "two  spaces")
			So(HTML2Text(`&copy; 2017 K3A`), ShouldEqual, "© 2017 K3A")
			So(HTML2Text("&lt;printtag&gt;"), ShouldEqual, "<printtag>")
			So(HTML2Text(`would you pay in &cent;, &pound;, &yen; or &euro;?`),
				ShouldEqual, "would you pay in ¢, £, ¥ or €?")
			So(HTML2Text(`Tom & Jerry is not an entity`), ShouldEqual, "Tom & Jerry is not an entity")
			So(HTML2Text(`this &neither; as you see`), ShouldEqual, "this &neither; as you see")
			So(HTML2Text(`list of items<ul><li>One</li><li>Two</li><li>Three</li></ul>`), ShouldEqual, "list of items\r\nOne\r\nTwo\r\nThree\r\n")
			So(HTML2Text(`fish &amp; chips`), ShouldEqual, "fish & chips")
			So(HTML2Text(`&quot;I'm sorry, Dave. I'm afraid I can't do that.&quot; – HAL, 2001: A Space Odyssey`), ShouldEqual, "\"I'm sorry, Dave. I'm afraid I can't do that.\" – HAL, 2001: A Space Odyssey")
			So(HTML2Text(`Google &reg;`), ShouldEqual, "Google ®")
			So(HTML2Text(`&#8268; decimal and hex entities supported &#x204D;`), ShouldEqual, "⁌ decimal and hex entities supported ⁍")
		})

		Convey("Large Entity", func() {
			So(HTMLEntitiesToText("&abcdefghij;"), ShouldEqual, "&abcdefghij;")
		})

		Convey("Numeric HTML Entities", func() {
			So(HTMLEntitiesToText("&#39;single quotes&#39; and &#52765;"), ShouldEqual, "'single quotes' and 츝")
			So(HTML2Text(`a null&#0; character is ignored`), ShouldEqual, "a null character is ignored")
			So(HTML2Text(`a control&#x1b; character is ignored`), ShouldEqual, "a control character is ignored")
			So(HTML2Text(`space entities&#x20;&#x20;collapsed correctly`), ShouldEqual, "space entities collapsed correctly")
		})

		Convey("HTML Entities in hrefs", func() {
			So(HTML2TextWithOptions(`click <a href="./filenamewith  spaces">here</a>`, WithLinksInnerText()), ShouldEqual, "click here <./filenamewith  spaces>")
			So(HTML2TextWithOptions(`click <a href="ents/control&#9;char">here</a>`, WithLinksInnerText()), ShouldEqual, "click here <ents/controlchar>")
			So(HTML2Text(`<a href="http://example.com/`+"\r\nBcc: victim@example.com"+`">link</a>`), ShouldEqual, "http://example.com/Bcc: victim@example.com")
		})

		Convey("Full HTML structure", func() {
			So(HTML2Text(``), ShouldEqual, "")
			So(HTML2Text("Bad chars\x1b ignored"), ShouldEqual, "Bad chars ignored")
			So(HTML2Text(`<html><head><title>Good</title></head><body>x</body>`), ShouldEqual, "x")
			So(HTML2Text(`<html><head href="foo"><title>Good</title></head><body>x</body>`), ShouldEqual, "x")
			So(HTML2Text(`<htMl><hEad><titLe>Good</Title></head><boDy>x</Body>`), ShouldEqual, "x")
			So(HTML2Text(`we are not <script type="javascript"></script>interested in scripts`),
				ShouldEqual, "we are not interested in scripts")
		})

		Convey("Switching Unix and Windows line breaks (original behavior)", func() {
			SetUnixLbr(true)
			So(HTML2Text(`two<br>line<br/>breaks`), ShouldEqual, "two\nline\nbreaks")
			So(HTML2Text(`<p>two</p><p>paragraphs</p>`), ShouldEqual, "two\n\nparagraphs")
			SetUnixLbr(false)
			So(HTML2Text(`two<br>line<br/>breaks`), ShouldEqual, "two\r\nline\r\nbreaks")
			So(HTML2Text(`<p>two</p><p>paragraphs</p>`), ShouldEqual, "two\r\n\r\nparagraphs")
		})

		Convey("Switching Unix and Windows line breaks (new options)", func() {
			So(HTML2TextWithOptions(`two<br>line<br/>breaks`, WithUnixLineBreaks()), ShouldEqual, "two\nline\nbreaks")
			So(HTML2TextWithOptions(`<p>two</p><p>paragraphs</p>`, WithUnixLineBreaks()), ShouldEqual, "two\n\nparagraphs")
			So(HTML2TextWithOptions(`two<br>line<br/>breaks`), ShouldEqual, "two\r\nline\r\nbreaks")
			So(HTML2TextWithOptions(`<p>two</p><p>paragraphs</p>`), ShouldEqual, "two\r\n\r\nparagraphs")
		})

		Convey("No list support by default (original behavior)", func() {
			So(HTML2Text(`list of items<ul><li>One</li><li>Two</li><li>Three</li></ul>`), ShouldEqual, "list of items\r\nOne\r\nTwo\r\nThree\r\n")
		})

		Convey("Tags with attributes", func() {
			So(HTML2Text(`list of items<ul><li class="menu-item">One</li><li class="menu-item">Two</li><li class="menu-item">Three</li></ul>`), ShouldEqual, "list of items\r\nOne\r\nTwo\r\nThree\r\n")
			So(HTML2Text(`list of items<ol><li class="menu-item">One</li><li class="menu-item">Two</li><li class="menu-item">Three</li></ol>`), ShouldEqual, "list of items\r\nOne\r\nTwo\r\nThree\r\n")
			So(HTML2Text(`<p class="content">content</p><div id="status">is ok</div>`), ShouldEqual, "content\r\n\r\nis ok")
		})

		Convey("Optional list support", func() {
			So(HTML2TextWithOptions(`list of items<ul><li>One</li><li>Two</li><li>Three</li></ul>`, WithListSupport()), ShouldEqual, "list of items\r\n - One\r\n - Two\r\n - Three\r\n")
			So(HTML2TextWithOptions(`list of items<ol><li>One</li><li>Two</li><li>Three</li></ol>`, WithListSupport()), ShouldEqual, "list of items\r\n - One\r\n - Two\r\n - Three\r\n")
		})

		Convey("Custom HTML Tags", func() {
			So(HTML2Text(`<aa>hello</aa>`), ShouldEqual, "hello")
			So(HTML2Text(`<aa >hello</aa>`), ShouldEqual, "hello")
			So(HTML2Text(`<aa x="1">hello</aa>`), ShouldEqual, "hello")
		})

		Convey("Keep spaces as they are", func() {
			So(HTML2TextWithOptions("should not    ignore spaces", WithKeepSpaces()), ShouldEqual, "should not    ignore spaces")
		})

		Convey("URL scheme validation and dangerous schemes blocked by default", func() {
			Convey("Garbage and unsafe URLs are ignored", func() {
				So(HTML2Text(`<a href=":">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#x09;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="http://example.com/&#x1b;ok">click</a>`), ShouldEqual, "http://example.com/ok")
				So(HTML2Text(`<a href="http://example.com/\x1b[2J\x1b[H">click</a>`), ShouldEqual, "http://example.com//x1b[2J/x1b[H")
			})

			Convey("Dangerous schemes blocked by default", func() {
				So(HTML2Text(`<a href="data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="DATA:text/html,test">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="vbscript:msgbox(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="VBSCRIPT:msgbox(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="file:///etc/passwd">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="FILE:///etc/passwd">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="blob:http://example.com/uuid">click</a>`), ShouldEqual, "")
			})

			Convey("WithLinksInnerText suppresses dangerous schemes but keeps inner text", func() {
				So(HTML2TextWithOptions(`<a href="data:text/html,test">click here</a>`, WithLinksInnerText()), ShouldEqual, "click here")
				So(HTML2TextWithOptions(`<a href="vbscript:msgbox(1)">click here</a>`, WithLinksInnerText()), ShouldEqual, "click here")
				So(HTML2TextWithOptions(`<a href="file:///etc/passwd">click here</a>`, WithLinksInnerText()), ShouldEqual, "click here")
			})

			Convey("Safe schemes allowed by default", func() {
				So(HTML2Text(`<a href="http://example.com">site</a>`), ShouldEqual, "http://example.com")
				So(HTML2Text(`<a href="https://example.com">site</a>`), ShouldEqual, "https://example.com")
				So(HTML2Text(`<a href="mailto:test@example.com">email</a>`), ShouldEqual, "mailto:test@example.com")
				So(HTML2Text(`<a href="tel:+33639981234">phone</a>`), ShouldEqual, "tel:+33639981234")
				So(HTML2TextWithOptions(`Send <a href="sms:+33639981234">SMS</a> for more info`, WithLinksInnerText()), ShouldEqual, "Send SMS <sms:+33639981234> for more info")
			})

			Convey("Relative URLs allowed by default", func() {
				So(HTML2Text(`<a href="/path/to/page">page</a>`), ShouldEqual, "/path/to/page")
				So(HTML2Text(`<a href="relative/path">page</a>`), ShouldEqual, "relative/path")
				So(HTML2Text(`<a href="#section">anchor</a>`), ShouldEqual, "#section")
			})

			Convey("Configured custom schemes allowed", func() {
				So(HTML2TextWithOptions(`<a href="ourapp://link">open app</a>`, WithAllowedURLSchemes([]string{"ourapp"})), ShouldEqual, "ourapp://link")
			})
		})

		Convey("URLs with whitespaces and zero-width characters are ignored", func() {
			Convey("Non-breaking spaces ignored: \\u00A0 (&nbsp;), \\u202F, \\u2007", func() {
				So(HTML2Text(`<a href="java&nbsp;script:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="&nbsp;javascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="java&#160;script:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="java&#xa0;script:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="java&#x202f;script:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="java&#x2007;script:alert(1)">click</a>`), ShouldEqual, "")
			})

			Convey("Zero-width spaces ignored: \\u200B (&#8203;), \\u200C, \\u200D, \\u2060, \\uFEFF", func() {
				So(HTML2Text(`<a href="jav&#8203;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#x200b;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#8204;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#x200c;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#8205;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#x200d;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#8288;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#x2060;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#65279;ascript:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="jav&#xfeff;ascript:alert(1)">click</a>`), ShouldEqual, "")
			})

			Convey("Other Unicode spaces ignored: \\u1680, \\u2000–\\u200A, \\u3000", func() {
				So(HTML2Text(`<a href="java&#x1680;script:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="java&#x2000;script:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="java&#x2001;script:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="java&#x2008;script:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="java&#x3000;script:alert(1)">click</a>`), ShouldEqual, "")
			})

			Convey("Whitespace or zero-width before colon is invalid", func() {
				So(HTML2Text(`<a href="javascript :alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="javascript&nbsp;:alert(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="javascript&#8203;:alert(1)">click</a>`), ShouldEqual, "")
			})

			Convey("Dangerous schemes obfuscated with unicode spaces / zero-width characters are ignored", func() {
				So(HTML2Text(`<a href="d&nbsp;ata:text/html,abc">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="v&#8203;bscript:msgbox(1)">click</a>`), ShouldEqual, "")
				So(HTML2Text(`<a href="f&#xfeff;ile:///etc/passwd">click</a>`), ShouldEqual, "")
			})

			Convey("WithLinksInnerText retains inner text when bad href is suppressed", func() {
				So(HTML2TextWithOptions(`<a href="java&nbsp;script:alert(1)">click here</a>`, WithLinksInnerText()), ShouldEqual, "click here")
				So(HTML2TextWithOptions(`<a href="jav&#8203;ascript:alert(1)">click here</a>`, WithLinksInnerText()), ShouldEqual, "click here")
				So(HTML2TextWithOptions(`<a href="jav&#x0a;ascript:alert(1)">click</a>`, WithLinksInnerText()), ShouldEqual, "click")
			})
		})

		Convey("Unmatched closing tags (negative depth flaw fix)", func() {
			Convey("Unmatched closing tags do not suppress subsequent text", func() {
				So(HTML2Text(`</head>Hello world`), ShouldEqual, "Hello world")
				So(HTML2Text(`</a>Hello world`), ShouldEqual, "Hello world")
			})

			Convey("Unmatched closing tags followed by bad tags do not leak contents", func() {
				So(HTML2Text(`</script><script>sensitive_token</script>Hello world`), ShouldEqual, "Hello world")
				So(HTML2Text(`</head><head><title>sensitive title</title></head>Hello world`), ShouldEqual, "Hello world")
			})
		})

		Convey("Anchor tag state desynchronization in WithLinksInnerText mode", func() {
			Convey("Anchors without href do not drop subsequent text", func() {
				So(HTML2TextWithOptions(`<a id="top"></a>Hello world`, WithLinksInnerText()), ShouldEqual, "Hello world")
				So(HTML2TextWithOptions(`<a></a>Hello world`, WithLinksInnerText()), ShouldEqual, "Hello world")
			})

			Convey("Inner text of anchors without href is retained", func() {
				So(HTML2TextWithOptions(`<a id="top">Anchor</a> Hello world`, WithLinksInnerText()), ShouldEqual, "Anchor Hello world")
				So(HTML2TextWithOptions(`<a>Anchor</a> Hello world`, WithLinksInnerText()), ShouldEqual, "Anchor Hello world")
			})

			Convey("Inner text is retained for mix of anchors with and without href", func() {
				So(HTML2TextWithOptions(`<a id="section-1">Section 1</a> <a href="http://example.com">Link</a> and <a name="section-2">Section 2</a>`, WithLinksInnerText()), ShouldEqual, "Section 1 Link <http://example.com> and Section 2")
			})

			Convey("Unclosed anchors without href do not drop subsequent text", func() {
				So(HTML2TextWithOptions(`<a id="top">Hello world`, WithLinksInnerText()), ShouldEqual, "Hello world")
			})

			Convey("Bad tags remain suppressed when anchors are present", func() {
				So(HTML2TextWithOptions(`<head><a id="foo"></a><title>sensitive title</title></head>Hello world`, WithLinksInnerText()), ShouldEqual, "Hello world")
			})

			Convey("Unmatched closing </a> in WithLinksInnerText mode won't suppress subsequent text", func() {
				So(HTML2TextWithOptions(`</a>Hello world`, WithLinksInnerText()), ShouldEqual, "Hello world")
			})
		})

	})
}
