package types

import "strings"

type TermQuery struct {
	Should  []*TermQuery //多个TermQuery与
	Must    []*TermQuery //多个TermQuery或
	Keyword string       //叶子节点本身TermQuery
}

func Be(str string) *TermQuery {
	return &TermQuery{Keyword: str}
}
func (v TermQuery) Empty() bool {
	return len(v.Should) == 0 && len(v.Must) == 0 && len(v.Keyword) == 0
}

func (q *TermQuery) And(querys ...*TermQuery) *TermQuery {
	if len(querys) == 0 {
		return q
	}

	array := make([]*TermQuery, 0, 1+len(querys)) //包含本身
	if !q.Empty() {
		array = append(array, q) //如果q为空则排除
	}
	for _, exp := range querys {
		if !exp.Empty() {
			array = append(array, exp)
		}
	}
	return &TermQuery{Must: array}
}

func (q *TermQuery) Or(querys ...*TermQuery) *TermQuery {
	if len(querys) == 0 {
		return q
	}

	array := make([]*TermQuery, 0, 1+len(querys)) //包含本身
	if !q.Empty() {
		array = append(array, q) //如果q为空则排除
	}
	for _, exp := range querys {
		if !exp.Empty() {
			array = append(array, exp)
		}
	}
	return &TermQuery{Should: array}
}

func (exp *TermQuery) ToString() string {
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
}
