package dongo

import (
	"github.com/flosch/pongo2"
	"bytes"
	"net/http"
)

var views = make(map[string]string)

type tagUrlNode struct {
	name pongo2.IEvaluator
}

func tagUrlParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, error) {
	urlNode := &tagUrlNode{}

	name, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	urlNode.name = name

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Static name is malformed.", nil)
	}

	return urlNode, nil
}

func (node *tagUrlNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) (error) {
	viewName, err := node.name.Evaluate(ctx)
	if err != nil {
		return err
	}

	buffer.WriteString(views[viewName.String()])

	return nil
}

func init() {
	pongo2.RegisterTag("url", tagUrlParser)
}

func ServeView(name, path string, handler func(writer http.ResponseWriter, request *http.Request)) {
	views[name] = path
	http.HandleFunc(path, handler)
}
