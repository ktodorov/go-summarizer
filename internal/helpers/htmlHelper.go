package helpers

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"regexp"

	"golang.org/x/net/html"
)

var negative, _ = regexp.Compile("/hidden|^hid$| hid$| hid |^hid |banner|combx|comment|com-|contact|foot|footer|footnote|masthead|media|meta|modal|outbrain|promo|related|scroll|share|shoutbox|sidebar|skyscraper|sponsor|shopping|tags|tool|widget/i,")
var positive, _ = regexp.Compile("/article|body|content|entry|hentry|h-entry|main|page|pagination|post|text|blog|story/i")

func getAttribute(node *html.Node, attributeName string) (attributeValue string, found bool) {
	var attributes = node.Attr
	for _, attr := range attributes {
		if attr.Key == attributeName {
			attributeValue = attr.Val
			found = true
			return
		}
	}

	return "", false
}

func getMainContentFromHTML(node *html.Node) (mainText string, images []string) {
	var pNodes = extractNodes(node, "p")
	var parents = make(map[*html.Node]float64)

	for _, pNode := range pNodes {
		var nodeParent = pNode.Parent
		var nodeGrandParent = nodeParent.Parent

		if _, exists := parents[nodeParent]; !exists {
			parents[nodeParent] = 0
		}

		if _, exists := parents[nodeGrandParent]; !exists {
			parents[nodeGrandParent] = 0
		}

		// Examine class attribute
		class, found := getAttribute(nodeParent, "class")
		if found {
			if negative.MatchString(class) {
				parents[nodeParent] -= 50
				parents[nodeGrandParent] -= 25
			} else if positive.MatchString(class) {
				parents[nodeParent] += 25
				parents[nodeGrandParent] += 12.5
			}
		}

		// Examine id attribute
		id, found := getAttribute(nodeParent, "id")
		if found {
			if negative.MatchString(id) {
				parents[nodeParent] -= 50
				parents[nodeGrandParent] -= 25
			} else if positive.MatchString(id) {
				parents[nodeParent] += 25
				parents[nodeGrandParent] += 12.5
			}
		}

		// Examine p tag length
		if len(pNode.Data) > 10 {
			parents[nodeParent]++
		}
	}

	var maxValue = 0.0
	var maxNodes []*html.Node

	for key, value := range parents {
		if value > maxValue {
			maxNodes = []*html.Node{key}
		} else if value == maxValue {
			maxNodes = append(maxNodes, key)
		}
	}

	if maxNodes == nil || len(maxNodes) == 0 {
		return "", nil
	}

	var textNodes = extractNodesFromMultipleParents(maxNodes, "p")
	mainText = ""

	for _, textNode := range textNodes {
		var nodeText = extractTextFromNode(textNode)
		mainText += "\n" + nodeText
	}

	images = []string{}

	var imageNodes = extractNodesFromMultipleParents(maxNodes, "img")
	for _, imageNode := range imageNodes {
		var imageSource, found = getAttribute(imageNode, "src")
		if found {
			images = append(images, imageSource)
		}
	}

	return mainText, images
}

func extractNodesFromMultipleParents(nodes []*html.Node, tag string) []*html.Node {
	var allNodes = []*html.Node{}

	for _, node := range nodes {
		var currentNodes = extractNodes(node, tag)
		allNodes = append(allNodes, currentNodes...)
	}

	return allNodes
}

func extractNodes(node *html.Node, tag string) []*html.Node {
	var allNodes = []*html.Node{}

	if node.Type == html.ElementNode && node.Data == tag {
		allNodes = append(allNodes, node)
	}

	for currentNode := node.FirstChild; currentNode != nil; currentNode = currentNode.NextSibling {
		var extractedNodes = extractNodes(currentNode, tag)
		if len(extractedNodes) > 0 {
			allNodes = append(allNodes, extractedNodes...)
		}
	}

	return allNodes
}

func extractNode(node *html.Node, nodeName string) (*html.Node, error) {
	if node.Type == html.ElementNode && node.Data == nodeName {
		var desiredNode = node
		return desiredNode, nil
	}

	for currentNode := node.FirstChild; currentNode != nil; currentNode = currentNode.NextSibling {
		var node, err = extractNode(currentNode, nodeName)
		if err == nil && node != nil {
			return node, nil
		}
	}

	return nil, errors.New("Missing <" + nodeName + "> in the node tree")
}

func removeNodesFromNode(node *html.Node, nodeToRemove string) {
	var childrenToRemove = []*html.Node{}
	for currentNode := node.FirstChild; currentNode != nil; currentNode = currentNode.NextSibling {
		if currentNode.Type == html.ElementNode && currentNode.Data == nodeToRemove {
			childrenToRemove = append(childrenToRemove, currentNode)
		} else {
			removeNodesFromNode(currentNode, nodeToRemove)
		}
	}

	for _, nodeElement := range childrenToRemove {
		node.RemoveChild(nodeElement)
	}
}

func extractTextFromNode(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	var result = ""

	for currentNode := node.FirstChild; currentNode != nil; currentNode = currentNode.NextSibling {
		if currentNode.Type == html.TextNode {
			result = result + currentNode.Data
		} else if currentNode.Data == "br" {
			result = result + "\n"
		} else {
			var currentNodeText = extractTextFromNode(currentNode)
			result = result + currentNodeText
		}
	}

	result = strings.TrimSpace(result)

	return result
}

func renderNode(node *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, node)

	var nodeText = buf.String()
	var nodeData = node.Data

	nodeText = strings.Replace(nodeText, "<"+nodeData+">", "", 1)
	nodeText = strings.Replace(nodeText, "</"+nodeData+">", "", 1)
	nodeText = strings.Replace(nodeText, "  ", " ", -1)
	nodeText = strings.TrimSpace(nodeText)

	return nodeText
}

func getPageTitle(node *html.Node) string {
	var titleNodes = extractNodes(node, "title")
	var title = ""
	if len(titleNodes) > 0 {
		title = extractTextFromNode(titleNodes[0])
	}

	var originalTitle = title

	// We check if there is header equal to document title
	// If it is the same, then this is most likely the article title
	var headerNodes = extractNodes(node, "h1")
	var header2Nodes = extractNodes(node, "h2")
	headerNodes = append(headerNodes, header2Nodes...)

	if len(headerNodes) > 0 {
		for _, headerNode := range headerNodes {
			var headerText = extractTextFromNode(headerNode)
			if headerText == title {
				return headerText
			}
		}
	}

	if strings.Contains(originalTitle, ":") {
		var splitTitle = strings.Split(originalTitle, ":")
		if len(splitTitle) > 0 {
			title = splitTitle[len(splitTitle)-1]   // last title part
			if len(strings.Split(title, " ")) < 3 { // if the new title is too short
				title = originalTitle
			}
		}
	} else if len(originalTitle) < 15 || len(originalTitle) > 150 {
		// If the original title is too big or too small, we get the first header
		var firstHeaders = extractNodes(node, "h1")
		if len(firstHeaders) == 1 {
			title = extractTextFromNode(firstHeaders[0])
		}
	}

	title = strings.TrimSpace(title)
	return title
}

func replaceBrs(htmlBody *html.Node) (*html.Node, error) {
	for true {
		var brNode, _ = extractNode(htmlBody, "br")
		if brNode == nil {
			// If there are no more br nodes in the body, we exit
			return htmlBody, nil
		}

		var replaced = false
		for brNode.NextSibling.Data == "br" {
			replaced = true
			var brSibling = brNode.NextSibling
			brSibling.Parent.RemoveChild(brNode)
			brNode = brSibling
		}

		if replaced {
			var p = html.Node{Data: "p"}
			htmlBody.InsertBefore(&p, brNode)
			htmlBody.RemoveChild(brNode)

			var next = p.NextSibling
			for next != nil {
				// If we meet another <br><br> elements, we end the new p tag here
				if next.Data == "br" {
					var nextElem = next.NextSibling
					if nextElem != nil && nextElem.Data == "br" {
						break
					}
				}

				// Add this element as child to the new p tag
				var sibling = next.NextSibling
				p.AppendChild(next)
				next = sibling
			}
		} else {
			brNode.Parent.RemoveChild(brNode)
		}
	}

	return nil, errors.New("Error occured while parsing br nodes")
}

func getMainInfoFromHTML(htmlString string) (string, string, []string, error) {
	doc, _ := html.Parse(strings.NewReader(htmlString))
	bn, err := extractNode(doc, "body")
	if err != nil {
		return "", "", nil, err
	}
	removeNodesFromNode(bn, "script")
	removeNodesFromNode(bn, "style")
	removeNodesFromNode(bn, "form")

	bn, err = replaceBrs(bn)
	if err != nil {
		return "", "", nil, err
	}

	var mainText, mainImages = getMainContentFromHTML(bn)
	var title = getPageTitle(bn)

	return title, mainText, mainImages, nil
}
