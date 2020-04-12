package util

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/micro-in-cn/x-apisix/monitor/task"
	"github.com/micro/go-micro/v2/registry"
)

//	自动生成详情key
const AUTO_PREFIX = "AUTO"

func add(a, b int) int {
	return a + b
}
func sub(a, b int) int {
	return a - b
}
func trimSuffix(str, suffix string) string {
	return strings.TrimSuffix(str, suffix)
}
func trimPrefix(str, prefix string) string {
	return strings.TrimPrefix(str, prefix)
}
func split(s, sep string) []string {
	return strings.Split(s, sep)
}

func MatchID(str string) string {
	i := strings.LastIndex(str, "/")
	l := len(str) - 1
	if i > -1 && i <= l {
		return str[i+1:]
	}
	return ""
}

//	匹配自动生成的前缀
func MatchAutoPrefixOfKey(str string) bool {
	return strings.HasPrefix(str, AUTO_PREFIX)
}

//  生成路由的唯一详情标记
func MakeRouteDesc(t *task.Svariable, endpoint *registry.Endpoint) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s-%s", AUTO_PREFIX, t.BU, t.Stype, t.Version, t.Module, endpoint.Name)
}

//  生成service唯一详情标记
func MakeServiceDesc(t *task.Svariable) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", AUTO_PREFIX, t.BU, t.Stype, t.Version, t.Module)
}

//  生成upstream唯一详情标记
func MakeUpstreamDesc(t *task.Svariable) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", AUTO_PREFIX, t.BU, t.Stype, t.Version, t.Module)
}

//	格式化gateway配置文件
func FormatGateway(wr io.Writer, tplName, tpl string, data interface{}) error {
	tmpl, err := template.New(tplName).
		Funcs(template.FuncMap{"add": add, "sub": sub, "trimSuffix": trimSuffix, "trimPrefix": trimPrefix, "split": split}).
		Parse(tpl)
	if err != nil {
		return err
	}
	// 数据驱动模板
	return tmpl.Execute(wr, data)
}

func JsonRequestBody(wr *bytes.Buffer, filename string, tplFormat, data interface{}) error {
	if err := FormatGatewayFile(wr, filename, data); err != nil {
		return err
	}
	switch tplFormat {
	case "json":
		return nil
	case "yaml", "yml":
		s, yerr := yaml.YAMLToJSON(wr.Bytes())
		if yerr != nil {
			return yerr
		}
		wr.Reset()
		wr.Write(s)
	}
	return nil
}

//	格式化gateway配置文件
func FormatGatewayFile(wr io.Writer, filename string, data interface{}) error {
	tmpl, err := template.New(filepath.Base(filename)).
		Funcs(template.FuncMap{"add": add, "sub": sub, "trimSuffix": trimSuffix, "trimPrefix": trimPrefix, "split": split}).ParseFiles(filename)
	if err != nil {
		return err
	}
	return tmpl.Execute(wr, data)
}

//	匹配proto文件列表
func GetProtoFileList(protoPath string, svariable *task.Svariable) ([]string, error) {
	protoPath = fmt.Sprintf("%s/%s/%s/%s/%s/*.proto",
		protoPath,
		svariable.BU,
		svariable.Module,
		svariable.Stype,
		svariable.Sver)
	return filepath.Glob(protoPath)
}
