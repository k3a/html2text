package html2text

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestHTML2Text(t *testing.T) {
	Convey("HTML2Text should work", t, func() {

		Convey("Links", func() {
			So(HTML2Text(`<div></div>`), ShouldEqual, "")
			So(HTML2Text(`<div>simple text</div>`), ShouldEqual, "simple text")
			So(HTML2Text(`click <a href="test">here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a class="x" href="test">here</a>`), ShouldEqual, "click test")
			So(HTML2Text(`click <a href="ents/&apos;x&apos;">here</a>`), ShouldEqual, "click ents/'x'")
			So(HTML2Text(`click <a href="javascript:void(0)">here</a>`), ShouldEqual, "click ")
		})

		Convey("Inlines", func() {
			So(HTML2Text(`strong <strong>text</strong>`), ShouldEqual, "strong text")
			So(HTML2Text(`some <div id="a" class="b">div</div>`), ShouldEqual, "some div")
		})

		Convey("Line breaks", func() {
			So(HTML2Text("should \nignore \r\nnew lines"), ShouldEqual, "should ignore new lines")
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
			So(HTML2Text(`two&nbsp;&nbsp;spaces`), ShouldEqual, "two  spaces")
			So(HTML2Text(`&copy; 2017 K3A`), ShouldEqual, "© 2017 K3A")
			So(HTML2Text("&lt;printtag&gt;"), ShouldEqual, "<printtag>")
			So(HTML2Text(`would you pay in &cent;, &pound;, &yen; or &euro;?`),
				ShouldEqual, "would you pay in ¢, £, ¥ or €?")
			So(HTML2Text(`Tom & Jerry is not an entity`), ShouldEqual, "Tom & Jerry is not an entity")
			So(HTML2Text(`this &neither; as you see`), ShouldEqual, "this &neither; as you see")
			So(HTML2Text(`fish &amp; chips`), ShouldEqual, "fish & chips")
			So(HTML2Text(`&quot;I'm sorry, Dave. I'm afraid I can't do that.&quot; – HAL, 2001: A Space Odyssey`), ShouldEqual, "\"I'm sorry, Dave. I'm afraid I can't do that.\" – HAL, 2001: A Space Odyssey")
			So(HTML2Text(`Google &reg;`), ShouldEqual, "Google ®")
		})

		Convey("Large Entity", func() {
			So(HTMLEntitiesToText("&abcdefghij;"), ShouldEqual, "&abcdefghij;")
		})

		Convey("Full HTML structure", func() {
			So(HTML2Text(``), ShouldEqual, "")
			So(HTML2Text(`<html><head><title>Good</title></head><body>x</body>`), ShouldEqual, "x")
			So(HTML2Text(`we are not <script type="javascript"></script>interested in scripts`),
				ShouldEqual, "we are not interested in scripts")
		})
	})
}
