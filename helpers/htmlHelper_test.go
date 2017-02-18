package helpers

import (
	"testing"

	"strings"

	"golang.org/x/net/html"
)

func createTestNode(t *testing.T, testHTMLString string, testHTMLTag string) *html.Node {
	htmlTestNode, _ := html.Parse(strings.NewReader(testHTMLString))
	testTag, err := extractNode(htmlTestNode, testHTMLTag)
	if err != nil {
		t.Error("expected no errors but received error: ", err.Error())
	}

	return testTag
}

func TestMissingAttributeGetter(t *testing.T) {
	var aTag = createTestNode(t, "<a>456</a>", "a")
	_, found := getAttribute(aTag, "test-attribute")
	if found {
		t.Error("getAttribute returned attribute found: ", found, ", when expected: false")
	}
}

func TestAvailableAttributeGetter(t *testing.T) {
	var aTag = createTestNode(t, "<a testattribute='123'>456</a>", "a")
	attrValue, found := getAttribute(aTag, "testattribute")

	if !found {
		t.Error("getAttribute returned attribute found: ", found, ", when expected: true")
	}

	if attrValue != "123" {
		t.Error("getAttribute returned attribute value: ", attrValue, ", when expected: 456")
	}
}

func TestExtractingMissingImagesFromNodes(t *testing.T) {
	var div1Tag = createTestNode(t, "<div><p></p><span></span></div>", "div")
	var div2Tag = createTestNode(t, "<div><p></p><span></span></div>", "div")
	var extractedImages = extractImagesFromNodes([]*html.Node{div1Tag, div2Tag})
	if extractedImages != nil && len(extractedImages) > 0 {
		t.Error("Expected no images but received: ", len(extractedImages))
	}
}

func TestExtractingAvailableImagesFromNodes(t *testing.T) {
	var div1Tag = createTestNode(t, "<div><p></p><span></span></div>", "div")
	var div2Tag = createTestNode(t, "<div><img src='testsrc' /><p></p><span></span><img src='test2src' /></div>", "div")
	var extractedImages = extractImagesFromNodes([]*html.Node{div1Tag, div2Tag})
	if extractedImages == nil || len(extractedImages) == 0 {
		t.Error("Expected images but received none")
	}

	if len(extractedImages) != 2 {
		t.Error("Expected 2 images but received ", len(extractedImages))
	}

	if extractedImages[0] != "testsrc" || extractedImages[1] != "test2src" {
		t.Error("Expected different images")
	}
}

func TestExtractingMissingNodesFromMultipleParents(t *testing.T) {
	var test1Tag = createTestNode(t, "<div><a href='testurl'>test</a><p>test text</p></div>", "div")
	var test2Tag = createTestNode(t, "<div><a href='testurl'>test</a><p>test text</p></div>", "div")
	var extractedNodes = extractNodesFromMultipleParents([]*html.Node{test1Tag, test2Tag}, "span")

	if extractedNodes != nil && len(extractedNodes) > 0 {
		t.Error("Expected no nodes but received: ", len(extractedNodes))
	}
}

func TestExtractingAvailableNodesFromMultipleParents(t *testing.T) {
	var test1Tag = createTestNode(t, "<div><a href='testurl'>test</a><p>test text</p></div>", "div")
	var test2Tag = createTestNode(t, "<div><a href='testurl'>test</a><p>test text</p></div>", "div")
	var extractedNodes = extractNodesFromMultipleParents([]*html.Node{test1Tag, test2Tag}, "p")

	if extractedNodes == nil && len(extractedNodes) == 0 {
		t.Error("Expected nodes but received none")
	}

	if len(extractedNodes) != 2 {
		t.Error("Expected 2 nodes but received ", len(extractedNodes))
	}
}

func TestExtractingMissingNode(t *testing.T) {
	var test1Tag = createTestNode(t, "<div><a href='testurl'>test</a><p>test text</p></div>", "div")
	var _, err = extractNode(test1Tag, "span")
	if err == nil {
		t.Error("Expected error for missing element but none received")
	}
}

func TestExtractingAvailableNode(t *testing.T) {
	var test1Tag = createTestNode(t, "<div><a href='testurl'>test</a><p>test text</p></div>", "div")
	var extractedNode, err = extractNode(test1Tag, "p")
	if err != nil {
		t.Error("Expected no error but received: ", err.Error())
	}

	if extractedNode == nil {
		t.Error("Expected node but received none")
	}
}

func TestRemovingNode(t *testing.T) {
	var test1Tag = createTestNode(t, "<div><a href='testurl'>test</a><p>test text</p></div>", "div")

	removeNodesFromNode(test1Tag, "p")
	var _, err = extractNode(test1Tag, "p")
	if err == nil {
		t.Error("Expected error for missing p tag but none received")
	}
}

func TestExtractingTextFromNode(t *testing.T) {
	var test1Tag = createTestNode(t, "<div><a href='testurl'>test</a><p>test text</p></div>", "div")

	var nodeText = extractTextFromNode(test1Tag)
	if nodeText != "testtest text" {
		t.Error("Expected to receive 'test test text' but got :", nodeText)
	}
}

func TestNodesTextDivFiltering(t *testing.T) {
	var test1Tag = createTestNode(t, "<body><div>test text</div></body>", "body")
	var resultTag = filterNodes(test1Tag)

	_, err := extractNode(resultTag, "div")
	if err == nil {
		t.Error("Expected error for missing div tag but none received")
	}

	pTag, err := extractNode(resultTag, "p")
	if err != nil {
		t.Error("Expected no error but received: ", err.Error())
	}

	if pTag == nil {
		t.Error("Expected div tag to be replaced by P tag but no p tags were found")
	}
}

func TestNodesDivWithPFiltering(t *testing.T) {
	var test1Tag = createTestNode(t, "<body><div><p>test text</p></div></body>", "body")
	var resultTag = filterNodes(test1Tag)

	_, err := extractNode(resultTag, "div")
	if err == nil {
		t.Error("Expected error for missing div tag but none received")
	}

	pTag, err := extractNode(resultTag, "p")
	if err != nil {
		t.Error("Expected no error but received: ", err.Error())
	}

	if pTag == nil {
		t.Error("Expected div tag to be replaced by P tag but no p tags were found")
	}
}

func TestBrReplacing(t *testing.T) {
	var testTag = createTestNode(t, "<body><br/><br/>first paragraph<br/><br/>second paragraph<br/><br/></body>", "body")
	var resultTag, err = replaceBrs(testTag)
	if err != nil {
		t.Error("Expected no error but received: ", err.Error())
	}

	var pTags = extractNodes(resultTag, "p")
	if pTags == nil {
		t.Error("Expected p tags but received none")
	} else if len(pTags) != 2 {
		t.Error("Expected 2 p tags but received: ", len(pTags))
	}
}
