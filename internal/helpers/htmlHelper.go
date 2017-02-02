package helpers

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"regexp"

	"golang.org/x/net/html"
)

var negative, _ = regexp.Compile(".*comment.*|.*meta.*|.*footer.*|.*foot.*|.*cloud.*|.*head.*")
var positive, _ = regexp.Compile(".*post.*|.*hentry.*|.*entry.*|.*content.*|.*text.*|.*body.*")

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
	var parents = make(map[*html.Node]int)

	for _, pNode := range pNodes {
		var nodeParent = pNode.Parent
		if _, exists := parents[nodeParent]; !exists {
			parents[nodeParent] = 0
		}

		// Examine class attribute
		class, found := getAttribute(nodeParent, "class")
		if found {
			if negative.MatchString(class) {
				parents[nodeParent] -= 50
			} else if positive.MatchString(class) {
				parents[nodeParent] += 25
			}
		}

		// Examine id attribute
		id, found := getAttribute(nodeParent, "id")
		if found {
			if negative.MatchString(id) {
				parents[nodeParent] -= 50
			} else if positive.MatchString(id) {
				parents[nodeParent] += 25
			}
		}

		// Examine p tag length
		if len(pNode.Data) > 10 {
			parents[nodeParent]++
		}
	}

	var maxValue = 0
	var maxNode *html.Node

	for key, value := range parents {
		if value > maxValue {
			maxNode = key
		}
	}

	if maxNode == nil {
		return "", nil
	}

	var textNodes = extractNodes(maxNode, "p")
	mainText = ""

	for _, textNode := range textNodes {
		var nodeText = extractTextFromNode(textNode)
		mainText += "\n" + nodeText
	}

	images = []string{}

	var imageNodes = extractNodes(maxNode, "img")
	for _, imageNode := range imageNodes {
		var imageSource, found = getAttribute(imageNode, "src")
		if found {
			images = append(images, imageSource)
		}
	}

	return mainText, images
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

func getMainInfoFromHTML(htmlString string) (string, []string, error) {
	doc, _ := html.Parse(strings.NewReader(htmlString))
	bn, err := extractNode(doc, "body")
	if err != nil {
		return "", nil, err
	}
	removeNodesFromNode(bn, "script")
	removeNodesFromNode(bn, "style")

	var mainText, mainImages = getMainContentFromHTML(bn)

	return mainText, mainImages, nil
}
