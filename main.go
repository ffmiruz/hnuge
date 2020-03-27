package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
	"text/template"
	"unicode"

	hn "github.com/montanaflynn/hn/hnclient"
)

// Comment node
type Node struct {
	Id    int
	By    string
	Time  int
	Url   string
	Text  string
	Level int
	Title string
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
    <title>{{.Title}}</title>
    <style type="text/css">
		aside {
		    border-left: 4px solid  #DAF7A6 ;
		    padding: 0.01rem 0.6rem;
		    border-bottom: 6px solid white;
		}
		article {
			margin-bottom: 6px;
		}
		p{margin: 6px auto;}
		.\31 {
        border-left-color:  #581845 ;
	    }
	    .\33 {
	    border-left-color: #C70039; 
	    }
	    .\34 {
	    border-left-color:  #FF5733 ;
	    }
	    .\35 {
	    border-left-color: #FFC300;
	    }
	    .\32 {
	    border-left-color: #900C3F;
	    }
	    code {
	    	display:block;overflow-x:hidden;
	    }
    </style>
</head>

<body>
  <h2>{{.Title}} </h2>
  {{.Text}} 
  {{template "comments" .Nodes}}
</body>
{{end}}

{{define "comments"}}
   {{- if . -}}
      {{range . }}                                  
         <aside class="{{.Level}}">                                         
           <div>                                  
             <div class="postTitle"><b>{{.By}}</b></div>   
           </div>
           <article>{{.Text}}</article>
           {{template "comments" .Nodes}}
         </aside>
      {{end}}
   {{- end -}}
{{end}}
`))

var client = hn.NewClient()

func main() {
	astory, _ := client.GetItem(22696377)
	rootStory := &Node{
		Id:    astory.Id,
		Kids:  astory.Kids,
		Text:  astory.Text,
		Url:   astory.Url,
		By:    astory.By,
		Title: astory.Title,
		Level: 0,
		Nodes: make([]*Node, len(astory.Kids)),
	}
	rootStory.fillNode()

	var b bytes.Buffer

	if err := t.ExecuteTemplate(&b, "mainStory", rootStory); err != nil {
		log.Fatalln(err)
	}
	a := StringMinifier(b.String())
	ioutil.WriteFile(strconv.Itoa(rootStory.Id)+".html", []byte(a), 0644)
}

func (node *Node) fillNode() {
	wg := &sync.WaitGroup{}
	wg.Add(len(node.Kids))
	for i, v := range node.Kids {
		go node.getComment(i, v, wg)

	}
	wg.Wait()
}

func (node *Node) getComment(i, v int, wg *sync.WaitGroup) {
	defer wg.Done()
	item, _ := client.GetItem(v)

	node.Nodes[i] = ItemToNode(item)
	node.Nodes[i].Level = node.Level + 1
	if len(item.Kids) > 0 {
		node.Nodes[i].Nodes = make([]*Node, len(item.Kids))
		node.Nodes[i].fillNode()
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

func StringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}
