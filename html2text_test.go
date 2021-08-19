// MIT License

// Copyright (c) 2017 Mario K3A Hros (www.k3a.me)

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
			So(HTML2Text(`click <a href="test">here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a class="x" href="test">here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href="ents/&apos;x&apos;">here</a>`), ShouldEqual, "click ents/'x'")
			So(HTML2Text(`click <a href="javascript:void(0)">here</a>`), ShouldEqual, "click ")
			So(HTML2Text(`click <a href="test"><span>here</span> or here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href="http://bit.ly/2n4wXRs">news</a>`), ShouldEqual, "click http://bit.ly/2n4wXRs")
			So(HTML2Text(`<a rel="mw:WikiLink" href="/wiki/yet#English" title="yet">yet</a>, <a rel="mw:WikiLink" href="/wiki/not_yet#English" title="not yet">not yet</a>`), ShouldEqual, "/wiki/yet#English, /wiki/not_yet#English")
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
			So(HTML2Text(`two&nbsp;&nbsp;spaces`), ShouldEqual, "two  spaces")
			So(HTML2Text(`&copy; 2017 K3A`), ShouldEqual, "© 2017 K3A")
			So(HTML2Text("&lt;printtag&gt;"), ShouldEqual, "<printtag>")
			So(HTML2Text(`would you pay in &cent;, &pound;, &yen; or &euro;?`),
				ShouldEqual, "would you pay in ¢, £, ¥ or €?")
			So(HTML2Text(`Tom & Jerry is not an entity`), ShouldEqual, "Tom & Jerry is not an entity")
			So(HTML2Text(`this &neither; as you see`), ShouldEqual, "this &neither; as you see")
			So(HTML2Text(`list of items<ul><li>One</li><li>Two</li><li>Three</li></ul>`), ShouldEqual, "list of items\r\n- One\r\n- Two\r\n- Three\r\n")
			So(HTML2Text(`list of items<ol><li>One</li><li>Two</li><li>Three</li></ol>`), ShouldEqual, "list of items\r\n- One\r\n- Two\r\n- Three\r\n")
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
		})

		Convey("Full HTML structure", func() {
			So(HTML2Text(``), ShouldEqual, "")
			So(HTML2Text(`<html><head><title>Good</title></head><body>x</body>`), ShouldEqual, "x")
			So(HTML2Text(`we are not <script type="javascript"></script>interested in scripts`),
				ShouldEqual, "we are not interested in scripts")
		})

		Convey("Switching Unix and Windows line breaks", func() {
			SetUnixLbr(true)
			So(HTML2Text(`two<br>line<br/>breaks`), ShouldEqual, "two\nline\nbreaks")
			So(HTML2Text(`<p>two</p><p>paragraphs</p>`), ShouldEqual, "two\n\nparagraphs")
			SetUnixLbr(false)
			So(HTML2Text(`two<br>line<br/>breaks`), ShouldEqual, "two\r\nline\r\nbreaks")
			So(HTML2Text(`<p>two</p><p>paragraphs</p>`), ShouldEqual, "two\r\n\r\nparagraphs")
		})

		Convey("Custom HTML Tags", func() {
			So(HTML2Text(`<aa>hello</aa>`), ShouldEqual, "hello")
			So(HTML2Text(`<aa >hello</aa>`), ShouldEqual, "hello")
			So(HTML2Text(`<aa x="1">hello</aa>`), ShouldEqual, "hello")
		})

		Convey("Split strings longer than 250", func() {
			So(HTML2Text(`<aa>UTCwUZ6sibdSsefMTrPg52Ows4uQSBXKjM0gJAvKxf7CeXYvKzf5lJd1I2nP3hXQ3soi9GZLYZmpumeIYlRIq0PJiDO0c9WHW6dU8OaHAFx1C5eF3QA22Cr6sVsWRbr1cBXPqtvWgfnM0KNgDgGKtg7afSLSIdTMNI6xNx8AFB0COO6zp17V2HWeCVVgwmWct5UEB1DlZG2m3TLfFtQ1hByB7fNgZNoRX67E4CdpDJezY7zoabasGDxFr2</aa>`),
				ShouldEqual, "UTCwUZ6sibdSsefMTrPg52Ows4uQSBXKjM0gJAvKxf7CeXYvKzf5lJd1I2nP3hXQ3soi9GZLYZmpumeIYlRIq0PJiDO0c9WHW6dU8OaHAFx1C5eF3QA22Cr6sVsWRbr1cBXPqtvWgfnM0KNgDgGKtg7afSLSIdTMNI6xNx8AFB0COO6zp17V2HWeCVVgwmWct5UEB1DlZG2m3TLfFtQ1hByB7fNgZNoRX67E4CdpDJezY7zoabasGDxFr2")
			So(HTML2Text(`<aa>UTCwUZ6sibdSsefMTrPg52Ows4uQSBXKjM0gJAvKxf7CeXYvKzf5lJd1I2nP3hXQ3soi9GZLYZmpumeIYlRIq0PJiDO0c9WHW6dU8OaHAFx1C5eF3QA22Cr6sVsWRbr1cBXPqtvWgfnM0KNgDgGKtg7afSLSIdTMNI6xNx8AFB0COO6zp17V2HWeCVVgwmWct5UEB1DlZG2m3TLfFtQ1hByB7fNgZNoRX67E4CdpDJezY7zoabasGDxFr2Rfino5IlkX1TyVRqx9dm2eTRdFAPW6xK4uPo8pVjmnqiQ4uOPtaJqkZhTwp60lX3fqrtzW1umqISS5Q6497tzBBghTwCoT2MC2GyRFJennWvvdjBsYQ9TQJbWfUxMLFKCG0PSis7Z3csFMJDwqd8kOImcrBvIShwm7nRa3zK16lZELCgNAtMVQPGpOTroS0o6w29tfX4C4S2KdzlCTiGu0QWO5Dmed4mIi8N2kml1dGNR6sHAoCqT5KWqE</aa>`),
				ShouldEqual, "UTCwUZ6sibdSsefMTrPg52Ows4uQSBXKjM0gJAvKxf7CeXYvKzf5lJd1I2nP3hXQ3soi9GZLYZmpumeIYlRIq0PJiDO0c9WHW6dU8OaHAFx1C5eF3QA22Cr6sVsWRbr1cBXPqtvWgfnM0KNgDgGKtg7afSLSIdTMNI6xNx8AFB0COO6zp17V2HWeCVVgwmWct5UEB1DlZG2m3TLfFtQ1hByB7fNgZNoRX67E4CdpDJezY7zoabasGDxFr2\r\nRfino5IlkX1TyVRqx9dm2eTRdFAPW6xK4uPo8pVjmnqiQ4uOPtaJqkZhTwp60lX3fqrtzW1umqISS5Q6497tzBBghTwCoT2MC2GyRFJennWvvdjBsYQ9TQJbWfUxMLFKCG0PSis7Z3csFMJDwqd8kOImcrBvIShwm7nRa3zK16lZELCgNAtMVQPGpOTroS0o6w29tfX4C4S2KdzlCTiGu0QWO5Dmed4mIi8N2kml1dGNR6sHAoCqT5KWqE")
		})

	})
}
