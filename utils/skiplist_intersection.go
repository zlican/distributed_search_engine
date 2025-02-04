package utils

import "github.com/huandu/skiplist"

func SkipIntersection(lists ...*skiplist.SkipList) *skiplist.SkipList {
	if len(lists) == 0 {
		return nil
	}
	if len(lists) == 1 {
		return lists[0]
	}
	result := skiplist.New(skiplist.Uint64) //返回的交集也用跳表表示

	currNodes := make([]*skiplist.Element, len(lists)) //维护一个跳表各自指针数组
	for i, list := range lists {
		if list == nil || list.Len() == 0 {
			return nil
		}
		currNodes[i] = list.Front()
	}

	for {
		maxList := make(map[int]struct{}, len(currNodes))
		var maxValue uint64 = 0
		for i, currNode := range currNodes {
			if currNode.Key().(uint64) > maxValue {
				maxValue = currNode.Key().(uint64)
				maxList = map[int]struct{}{i: {}} //格式化
			} else if currNode.Key().(uint64) == maxValue {
				maxList[i] = struct{}{}
			}
		}

		if len(maxList) == len(currNodes) {
			result.Set(currNodes[0].Key(), currNodes[0].Value)
			for i, node := range currNodes {
				currNodes[i] = node.Next()
				if currNodes[i] == nil {
					return result
				}
			}
		} else {
			for i, currNode := range currNodes {
				if currNode.Key().(uint64) == maxValue {
					continue
				}

				currNodes[i] = currNodes[i].Next()
				if currNodes[i] == nil {
					return result
				}
			}
		}
	}
}
