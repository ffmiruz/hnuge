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
<!DOCTYPE html>
<html lang="en">
<head>
    <link rel="icon" href="https://via.placeholder.com/70x70">
    <meta charset="utf-8">
    <meta name="description" content="My description">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My title</title>
    <style type="text/css">
		aside {
		    border-left: 4px solid red;
		    padding: 0.01rem 0.8rem;
		}
    </style>
</head>

<body>                                              
  {{.Text}} by {{.By}} 
  {{template "comments" .Nodes}}
</body>
{{end}}

{{define "comments"}}
   {{- if . -}}
      <article>
      {{range . }}                                  
         <aside>                                         
           <div>                                  
             <div class="postTitle"><b>{{.By}}</b></div>   
           </div>
           <div class="postBody">{{.Text}}</div>
           {{template "comments" .Nodes}}
         </aside>
      {{end}}
      </article>
   {{- end -}}
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

	file, err := os.Create(rootStory.By + ".html")
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
