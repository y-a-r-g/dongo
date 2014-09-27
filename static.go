package dongo

import (
	"github.com/flosch/pongo2"
	"bytes"
	"net/http"
	"fmt"
)

var staticUrlPrefix = "/static/"

type tagStaticNode struct {
	path pongo2.IEvaluator
}

func tagStaticParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, error) {
	staticNode := &tagStaticNode{}

	path, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	staticNode.path = path

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Static path is malformed.", nil)
	}

	return staticNode, nil
}

func (node *tagStaticNode) Execute(ctx *pongo2.ExecutionContext, buffer *bytes.Buffer) (error) {
	result, err := node.path.Evaluate(ctx)
	if err != nil {
		return err
	}

	buffer.WriteString(staticUrlPrefix+result.String())

	return nil
}

func init() {
	pongo2.RegisterTag("static", tagStaticParser)
}

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.RequestURI)
}

func ServeStatic(urlPrefix string, rootPath string) {
	staticUrlPrefix = urlPrefix

	fs := http.FileServer(http.Dir(rootPath))
	http.Handle(urlPrefix, http.StripPrefix(urlPrefix, fs))
}
