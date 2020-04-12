package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

var content []byte

func init() {
	protoPath := "/storage/code/aimgo.proto/aimgo/passport/http2/v1/"
	protoPath = fmt.Sprintf("%s/kaixin.proto", strings.TrimSuffix(protoPath, "/"))
	c, err := ioutil.ReadFile(protoPath)
	if err != nil {
		log.Fatal(err)
	}
	content = c
}
func TestKaixin(t *testing.T) {
	t.Log(MakeProtoID(content))
	t.Log(DeleteProtoHead(content))
	t.Log(CompressProto(content))
}
