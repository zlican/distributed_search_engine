package types

import "strings"

type TermQueryV0 struct {
	Should  []TermQueryV0 //多个TermQuery与
	Must    []TermQueryV0 //多个TermQuery或
	Keyword string        //叶子节点本身TermQuery
}

// (A|B|C)&E|F
func (v TermQueryV0) Empty() bool {
	return len(v.Should) == 0 && len(v.Must) == 0 && len(v.Keyword) == 0
}

func MustExpression(exps ...TermQueryV0) TermQueryV0 {
	if len(exps) == 0 {
		return TermQueryV0{}
	}

	array := make([]TermQueryV0, 0, len(exps))
	for _, exp := range exps {
		if !exp.Empty() {
			array = append(array, exp)
		}
	}

	return TermQueryV0{Must: array}
}

func ShouldExpression(exps ...TermQueryV0) TermQueryV0 {
	if len(exps) == 0 {
		return TermQueryV0{}
	}

	array := make([]TermQueryV0, 0, len(exps))
	for _, exp := range exps {
		if !exp.Empty() {
			array = append(array, exp)
		}
	}
	return TermQueryV0{Should: array}
}

func StringExpression(str string) TermQueryV0 {
	return TermQueryV0{Keyword: str}
}

func (exp TermQueryV0) ToString() string {
	if len(exp.Keyword) > 0 {
		return exp.Keyword
	} else if len(exp.Must) > 0 {
		if len(exp.Must) == 1 {
			return exp.Must[0].ToString() //递归
		} else {
			sb := strings.Builder{}
			sb.WriteString("(")
			for _, e := range exp.Must {
				sb.WriteString(e.ToString())
				sb.WriteByte('&')
			}
			s := sb.String()
			s = s[0:len(s)-1] + ")"
			return s
		}
	} else if len(exp.Should) > 0 {
		if len(exp.Should) == 1 {
			return exp.Should[0].ToString()
		} else {
			sb := strings.Builder{}
			sb.WriteString("(")
			for _, e := range exp.Should {
				sb.WriteString(e.ToString())
				sb.WriteByte('|')
			}
			s := sb.String()
			s = s[0:len(s)-1] + ")"
			return s
		}
	}
	return exp.ToString()
} //核心：将层层嵌套的TermQuery转化为string
