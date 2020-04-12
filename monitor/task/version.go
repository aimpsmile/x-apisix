/*
@Time : 2020/3/26 下午5:25
@Author : songxiuxuan
@File : version.go
@Software: GoLand
*/
package task

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

var Operator = map[string]bool{
	//  大于
	">": true,
	//  小于
	"<": true,
	//  相等
	"=": true,
	//  所有的意思
	"*": true,
	//  版本范围
	"~": true,
	//  范围
	",": true,
	"":  true,
}

//  两个版本之间做比对
func VersionCompare(va, vb float64, operator string, vc float64) bool {
	switch operator {
	case ">":
		return decimal.NewFromFloat(va).GreaterThan(decimal.NewFromFloat(vc))
	case ">=":
		return decimal.NewFromFloat(va).GreaterThanOrEqual(decimal.NewFromFloat(vc))
	case "<":
		return decimal.NewFromFloat(va).LessThan(decimal.NewFromFloat(vc))
	case "<=":
		return decimal.NewFromFloat(va).LessThanOrEqual(decimal.NewFromFloat(vc))
	case "=":
		return decimal.NewFromFloat(va).Equal(decimal.NewFromFloat(vc))
	case "~":
		ln := len(strings.Split(fmt.Sprintf("%v", vc), ".")[1])
		return decimal.NewFromFloat(va).Truncate(int32(ln)).Equal(decimal.NewFromFloat(vc))
	case ",":
		a := decimal.NewFromFloat(va)
		return a.GreaterThanOrEqual(decimal.NewFromFloat(vb)) && a.LessThanOrEqual(decimal.NewFromFloat(vc))
	default:
		return decimal.NewFromFloat(va).Equal(decimal.NewFromFloat(vc))
	}

}

//  字符串转运行操作符与字符串
func StringConverOperator(str string) (start, op, end string) {
	idx := 0
	for k, v := range str {
		s := string(v)
		if _, ok := Operator[s]; ok {
			op += s
			idx = k + 1
		}
	}
	end = str[idx:]
	//	如果第一位是操作符，需要设置为空
	start = str[:idx-len(op)]
	return start, op, end
}

//  字符串版本转float64的字符串
func StringConverFloat64(str, flag string, bitsize int) (float64, error) {
	if str == "" || str == "*" {
		return 0.0, nil
	}
	if len(str) > 1 && str[:1] == "v" {
		str = str[1:]
	}
	index := strings.Index(str, flag)
	if index > 0 {
		index += len(flag)
		prex := str[:index]
		after := strings.ReplaceAll(str[index:], ".", "")
		str = prex + after
	}
	return strconv.ParseFloat(str, bitsize)
}
