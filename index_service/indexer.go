package indexservice

import (
	"bytes"
	"encoding/gob"
	"engine/internal/kvdb"
	reverseindex "engine/internal/reverse_index"
	"engine/types"
	"fmt"
	"strings"
	"sync/atomic"
)

// 将正排和倒排结合起来
type Indexer struct {
	forwardIndex kvdb.IKeyValueDB
	reverseIndex reverseindex.IReverseIndexer
	maxIntId     uint64
}

func (indexer *Indexer) Init(DocNumEstimate int, dbtype int, DataDir string) error {
	db, err := kvdb.GetKvDb(dbtype, DataDir)
	if err != nil {
		return err
	}
	indexer.forwardIndex = db
	indexer.reverseIndex = reverseindex.NewSkipListReverseIndex(DocNumEstimate)
	return nil
}

func (indexer *Indexer) Close() error {
	return indexer.forwardIndex.Close()
}

// 往正排和倒排中写入doc
func (indexer *Indexer) AddDoc(doc types.Document) (int, error) {
	docId := strings.TrimSpace(doc.Id)
	if len(docId) == 0 {
		return 0, nil
	}

	indexer.DeleteDoc(docId) //删除旧的doc
	doc.IntId = atomic.AddUint64(&indexer.maxIntId, 1)

	var value bytes.Buffer
	encoder := gob.NewEncoder(&value)
	if err := encoder.Encode(doc); err == nil {
		indexer.forwardIndex.Set([]byte(docId), value.Bytes()) //将整个doc序列化成value
	} else {
		return 0, nil
	}

	indexer.reverseIndex.Add(doc)
	return 1, nil
}

func (indexer *Indexer) DeleteDoc(docId string) int {
	n := 0
	docBs, err := indexer.forwardIndex.Get([]byte(docId))
	if err == nil {
		reader := bytes.NewReader([]byte{})
		if len(docBs) > 0 {
			n = 1
			reader.Reset(docBs)
			decoder := gob.NewDecoder(reader)
			var doc types.Document
			err := decoder.Decode(&doc)
			if err == nil {
				for _, kw := range doc.Keywords {
					indexer.reverseIndex.Delete(doc.IntId, kw)
				}
			}
		}
	}
	indexer.forwardIndex.Delete([]byte(docId))
	return n
}

// 当系统重启时，直接从索引文件加载数据
func (indexer *Indexer) LoadFormIndexFile() int {
	reader := bytes.NewReader([]byte{})
	n := indexer.forwardIndex.IterDB(func(k, v []byte) error {
		reader.Reset(v)
		decoder := gob.NewDecoder(reader)
		var doc types.Document
		err := decoder.Decode(&doc)
		if err == nil {
			indexer.AddDoc(doc)
		}
		return err
	})
	return int(n)
}

func (indexer *Indexer) Search(q *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []*types.Document {
	docIds := indexer.reverseIndex.Search(q, onFlag, offFlag, orFlags)
	if len(docIds) == 0 {
		return nil
	}

	keys := make([][]byte, 0, len(docIds))
	for _, docId := range docIds {
		keys = append(keys, []byte(docId))
	}
	data, err := indexer.forwardIndex.BatchGet(keys)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	result := make([]*types.Document, 0, len(data))
	reader := bytes.NewReader(([]byte{}))
	for _, value := range data {
		reader.Reset(value)
		decoder := gob.NewDecoder(reader)
		var doc types.Document
		err := decoder.Decode(&doc)
		if err == nil {
			result = append(result, &doc)
		}
	}
	return result
}
