package util

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/micro-in-cn/x-apisix/core/lib/encrypt"
)

//	解释proto文件后的结构体
type ProtoType struct {
	PackageName string
	ServiceList []string
	ProtoID     string
}

const (
	regexpPackage = `package\s+(?P<kaixin>.*)?;`
	regexpService = `service\s+(?P<kaixin>.*)?{`
	headerPackge  = `package.*`
	headerSyntax  = `syntax.*`
	headerComment = `//.*`
)

//	匹配package name
func MatchPackage(s []byte) (string, error) {
	re, err := regexp.Compile(regexpPackage)
	if err != nil {
		return "", err
	}

	matchs := re.FindStringSubmatch(string(s))
	if len(matchs) >= 1 {
		return strings.TrimPrefix(matchs[1], " "), nil
	}
	return "", fmt.Errorf("regexp not match")
}

//	匹配service list
func MatchServiceList(s []byte) ([]string, error) {
	var serviceList []string
	re, err := regexp.Compile(regexpService)
	if err != nil {
		return serviceList, err
	}

	matchs := re.FindAllStringSubmatch(string(s), -1)
	for _, m := range matchs {
		if len(matchs) >= 1 {
			serviceList = append(serviceList, strings.TrimSuffix(strings.TrimPrefix(m[1], " "), " "))
		}
	}
	return serviceList, nil
}

//	删除公共头
func DeleteProtoHead(c []byte) (string, error) {
	s := string(c)
	if pre, perr := regexp.Compile(headerPackge); perr == nil {
		s = pre.ReplaceAllString(s, "")
	} else {
		return "", perr
	}
	if sre, serr := regexp.Compile(headerSyntax); serr == nil {
		s = sre.ReplaceAllString(s, "")
	} else {
		return "", serr
	}
	return s, nil
}

//	压缩Proto file
func CompressProto(c []byte) (string, error) {
	s := string(c)
	if re, rerr := regexp.Compile(headerComment); rerr == nil {
		s = re.ReplaceAllString(s, "")
	} else {
		return "", rerr

	}
	//将匹配到的部分替换为"##.#"
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s, nil
}

//	生成proto ID
func MakeProtoID(c []byte) (*ProtoType, error) {
	services, serr := MatchServiceList(c)
	if serr != nil {
		return nil, fmt.Errorf("[msg]match servie list is error[error_info]%w", serr)
	}

	packageName, perr := MatchPackage(c)
	if perr != nil {
		return nil, fmt.Errorf("[msg]match package name is error[error_info]%w", perr)
	}
	sort.Strings(services)
	s := strings.Join(services, "-")
	s = fmt.Sprintf("%s-%s", packageName, s)
	md5, err := encrypt.GetMD5([]byte(s))
	if err != nil {
		return nil, fmt.Errorf("[msg]make md5 is error[package.name]%s[services]%+v[error_info]%w", packageName, services, err)
	}

	p := &ProtoType{
		PackageName: packageName,
		ServiceList: services,
		ProtoID:     encrypt.HexEncodeToString(md5),
	}
	return p, nil
}
