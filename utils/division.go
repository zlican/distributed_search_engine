package utils

type Doc struct {
	Id       int
	KeyWords []string
}

func Division(docs []*Doc) map[string][]int {
	index := make(map[string][]int, 100)

	for _, doc := range docs {
		for _, keyWord := range doc.KeyWords {
			index[keyWord] = append(index[keyWord], doc.Id)
		}
	}

	return index
}
