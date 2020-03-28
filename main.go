package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
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

var tpl = template.Must(template.New("").Parse(`
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
			border: 1px solid #e9e9e9;
		    border-radius: 5px;
		    box-shadow: 2px 2px 10px #f4f4f4;
		    padding: 0.5rem;
			margin-bottom: 6px;
		}
		p{margin: 5px auto;}
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
	    a {
	    	color:#b00020;
	    }
	    .by{font-style:italic}
    </style>
</head>

<body>
  <header>
  <div>by  <b>{{.By}}</b></div>
   <a href="{{.Url}}"><h2>{{.Title}}</h2></a>
   {{if .Text}}
   <article>{{.Text}}</article>
   {{end}}
  </header>
  {{template "comments" .Nodes}}
</body>
{{end}}

{{define "comments"}}
   {{- if . -}}
      {{range . }}                                  
         <aside class="{{.Level}}">                                         
           <div class="by"><b>{{.By}}</b></div>
           {{.Text}}
           {{template "comments" .Nodes}}
         </aside>
      {{end}}
   {{- end -}}
{{end}}
`))

var client = hn.NewClient()

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage:", os.Args[0], "ITEM_ID")
	}
	id, _ := strconv.Atoi(os.Args[1])

	astory, err := client.GetItem(id)
	if err != nil {
		log.Println("Failed to fetch a story")
	}
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

	if err := tpl.ExecuteTemplate(&b, "mainStory", rootStory); err != nil {
		log.Fatalln(err)
	}
	a := minifier(b.String())
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

	node.Nodes[i] = itemToNode(item)
	node.Nodes[i].Level = node.Level + 1
	if len(item.Kids) > 0 {
		node.Nodes[i].Nodes = make([]*Node, len(item.Kids))
		node.Nodes[i].fillNode()
	}
}

func itemToNode(item hn.Item) *Node {
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

func minifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else if '>' == c {
			out = out + string(c)
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}
