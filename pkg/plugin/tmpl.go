package plugin

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	g "github/track/protoc-gen-examples/pkg/generator"
	"strings"
	"text/template"
)

type Method struct {
	Path   string
	Name   string
	Req    string
	Resp   string
	Action string
}

// 注册插件
func init() {
	g.RegisterPlugin(new(tmpl))
}

// 简写调用
func (t *tmpl) P(args ...interface{}) { t.gen.P(args...) }

// 插件:  tmpl-out  实现 组合Generator, 并实现generator/Plugin interface
type tmpl struct {
	gen   *g.Generator
	imp   map[string]string
	mainT *template.Template
}

// 初始化插件
func (t *tmpl) Init(g *g.Generator) {
	t.gen = g
	t.imp = make(map[string]string, 0)

	//初始化时间的导包路径
	t.imp["*time.Time"] = "time"
}
// main.Go文件, 仅有一个
func (t *tmpl) initMainFile(m []*Method) {
	if t.mainT != nil {
		return
	}
	t.mainT = template.New("cmd/main.go")
	t.mainT, _ = t.mainT.Parse(`
package main

import (
	"github/track/protoc-gen-examples/tmpl"
	"log"
	"net/http"
)
	
{{ range . }}
func {{ .Name }}(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != {{ .Action }} {
		writer.WriteHeader(http.StatusInternalServerError)
		_,_ = writer.Write([]byte("request method un support"))
	}
	var _ {{ .Req }}
	var _ {{ .Resp }}

	writer.WriteHeader(http.StatusOK)
	_,_=writer.Write([]byte("did "+method+" {{ .Name }} successful\n"))
}
{{ end }}

	
func main() {
	{{ range . }}
	http.HandleFunc("/{{ .Path }}", {{ .Name }})
	{{ end }}
	log.Fatal(http.ListenAndServe(":8000", nil))
}
			`)
	buf := new(bytes.Buffer)
	err := t.mainT.Execute(buf, m)
	if err == nil {
		// 写出文件
		t.gen.Response.File = append(t.gen.Response.File, &plugin_go.CodeGeneratorResponse_File{
			Name:    proto.String(t.mainT.Name()),
			Content: proto.String(buf.String()),
		})
	}
}
// 获取类型名
func (t *tmpl) typeName(str string) string {
	return t.gen.TypeName(t.gen.ObjectNamed(str))
}

// 生成具体文件
func (t *tmpl) Generate(file *g.FileDescriptor) {
	// 枚举类型生成
	for i, v := range file.GetEnumType() {
		ename := g.CamelCase(*v.Name)
		// 5 是枚举注释index
		t.gen.PrintComments(fmt.Sprintf("5,%d", i))
		t.P("type ", ename, " int32")
		t.P()
		t.P("const (")
		for j, e := range v.Value {
			// 5 2 这样子拼接就是子条目
			etorPath := fmt.Sprintf("5,%d,2,%d", i, j)
			t.gen.PrintComments(etorPath)
			// 拼接
			t.P(ename, "_", g.CamelCase(*e.Name), " ", ename, " = ", e.Number, " ")
		}
		t.P(")")
		t.P()
	}
	// 结构体类型
	for i, v := range file.GetMessageType() {
		// 4 是结构体注释index
		mname := g.CamelCase(*v.Name)
		t.gen.PrintComments(fmt.Sprintf("4,%d", i))
		t.P("type ", mname, " struct {")
		for j, f := range v.Field {
			// 4 2 这样子拼接就是子条目
			t.gen.PrintComments(fmt.Sprintf("4,%d,2,%d", i, j))
			fname := g.CamelCase(*f.Name)
			pname, _ := t.gen.GoTypePlugin(f)
			t.P(fmt.Sprintf("%s\t%s\t`json:%q`", fname, pname, f.GetJsonName()+",omitempty"))
		}
		t.P("}")
	}
	// 遍历Service
	for _, v := range file.GetService() {
		methods := make([]*Method, 0)
		// 遍历rpc函数
		for _, m := range v.GetMethod() {
			options := strings.Split(m.GetOptions().String(), ":")
			if len(options) < 2 { // 没有option不处理
				continue
			}
			// 组装函数结构
			methods = append(methods, &Method{
				Path:  strings.ToLower( m.GetName()),
				Name:   m.GetName(),
				Req:    file.GetPackage() + "." + t.typeName(m.GetInputType()),
				Resp:   file.GetPackage() + "." + t.typeName(m.GetOutputType()),
				Action: options[1],
			})
		}
		t.initMainFile(methods)
	}

}

// 插件导入包
func (t *tmpl) GenerateImports(file *g.FileDescriptor) {
	t.P("import (")
	for _, v := range file.GetMessageType() {
		for _, f := range v.Field {
			goType, _ := t.gen.GoTypePlugin(f)
			if value, exists := t.imp[goType]; exists {
				t.P(fmt.Sprintf(`"%s"`, value))
			}
		}
	}
	t.P(")")
	t.P()
}

// 返回插件名称
func (*tmpl) Name() string {
	return "tmpl"
}
