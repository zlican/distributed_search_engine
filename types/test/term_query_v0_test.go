package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zlican/engine/types"
)

// ((A|B|C)&E|(D|(A&B))

func should(s ...string) string {
	if len(s) == 0 { //排除空调用
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("(")
	for _, ele := range s {
		if len(ele) > 0 { //排除空的元素
			sb.WriteString(ele + "|")
		}
	}
	rect := sb.String()
	return rect[0:len(rect)-1] + ")"
	//return "(" + strings.Join(s, "|") + ")"
}

func must(s ...string) string {
	if len(s) == 0 { //排除空调用
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("(")
	for _, ele := range s {
		if len(ele) > 0 { //排除空的元素
			sb.WriteString(ele + "&")
		}
	}
	rect := sb.String()
	return rect[0:len(rect)-1] + ")"
	//return "(" + strings.Join(s, "&") + ")"
}

func TestN(t *testing.T) {
	fmt.Println(should(must(should("A", "B", "C"), "E"), should("D", must("A", "B"))))
}

func TestTradisionalTermQuery(t *testing.T) {
	//(A|B|C)&(E|(A&B))
	resultExp := types.MustExpression(types.ShouldExpression(types.StringExpression("A"), types.StringExpression("B"),
		types.StringExpression("C")), types.ShouldExpression(types.StringExpression("E"), types.MustExpression(types.StringExpression("A"),
		types.StringExpression("B"))))
	fmt.Println(resultExp.ToString())
} //嵌套地狱
