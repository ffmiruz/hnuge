package main

import (
	"fmt"
	hn "github.com/montanaflynn/hn/hnclient"
)

// Comment node
type Node struct {
	Id    int
	By    string
	Time  int
	Url   string
	Text  string
	Kids  []int
	Nodes []*Node
}

func main() {
	client := hn.NewClient()
	astory, _ := client.GetItem(22675078)
	rootStory := &Node{
		Kids:  astory.Kids,
		Nodes: make([]*Node, len(astory.Kids)),
	}
	fillNode(rootStory, client)

	fmt.Println(rootStory.Nodes[0].Nodes[1].Text)
}

func fillNode(node *Node, client *hn.Client) {
	for i, v := range node.Kids {
		item, _ := client.GetItem(v)

		node.Nodes[i] = ItemToNode(item)
		if len(item.Kids) > 0 {
			node.Nodes[i].Nodes = make([]*Node, len(item.Kids))
			fmt.Println(node.Nodes[i].Kids)
			fillNode(node.Nodes[i], client)
		}
		fmt.Println(node.Nodes[i].Nodes)
	}
}

func ItemToNode(item hn.Item) *Node {
	node := &Node{
		By:   item.By,
		Id:   item.Descendants,
		Time: item.Time,
		Text: item.Text,
		Url:  item.Url,
		Kids: item.Kids,
	}
	return node
}
