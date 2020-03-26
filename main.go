package main

import (
	"log"
	"os"
	"text/template"

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

var t = template.Must(template.New("").Parse(`
{{define "mainStory"}}
<div id="post">                                              
  {{.Text}} by {{.By}} 
  {{template "comments" .Nodes}}
</div>
{{end}}

{{define "comments"}}
   {{if .}}
      <ul>
      {{range . }}                                  
         <li class="post">                                         
           <div class="postHead">                                  
             <div class="postTitle"><b>{{.By}}</b></div>   
           </div>
           <div class="postBody">{{.Text}}</div>
           {{template "comments" .Nodes}}
         </li>
      {{end}}
      </ul>
   {{end}}
{{end}}
`))

func main() {
	client := hn.NewClient()
	astory, _ := client.GetItem(22675078)
	rootStory := &Node{
		Id:    astory.Id,
		Kids:  astory.Kids,
		Text:  astory.Title,
		Url:   astory.Url,
		By:    astory.By,
		Nodes: make([]*Node, len(astory.Kids)),
	}
	fillNode(rootStory, client)

	file, err := os.Create(rootStory.Url + ".html")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err := t.ExecuteTemplate(file, "mainStory", rootStory); err != nil {
		log.Fatalln(err)
	}
}

func fillNode(node *Node, client *hn.Client) {
	for i, v := range node.Kids {
		item, _ := client.GetItem(v)

		node.Nodes[i] = ItemToNode(item)
		if len(item.Kids) > 0 {
			node.Nodes[i].Nodes = make([]*Node, len(item.Kids))
			fillNode(node.Nodes[i], client)
		}
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
