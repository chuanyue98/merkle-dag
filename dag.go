package merkledag

import (
	"hash"
	"sort"
	"bytes"
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

// Add 将 Node 中的数据保存在 KVStore 中，并计算出 Merkle Root
func Add(store KVStore, node Node, h hash.Hash) []byte {
	// 如果节点是文件，则直接将其数据写入KVStore
	if node.Type() == FILE {
		fileData := node.(File).Bytes()
		h.Reset()
		h.Write(fileData)
		hashValue := h.Sum(nil)
		store.Put(hashValue, fileData)
		return hashValue
		}
	// 如果节点是目录，则递归处理子节点
	if node.Type() == DIR{
		var object Object
		var childHashs [][]byte
	   it := node.(Dir).It()
	   for it.Next() {
	       childNode := it.Node()
	       childHash := Add(store, childNode, h)
		   childHashs = append(childHashs, childHash)
		   object.Links = append(object.Links, Link{
			   Name: childNode.Name(),
			   Hash: childHash,
			   Size: int(childNode.Size()),
		   })
		   
	   }
	   sort.Slice(childHashs, func(i, j int) bool {
		return bytes.Compare(childHashs[i], childHashs[j]) < 0
		})
		h.Reset()
		for _, hash := range childHashs{
		h.Write(hash)
		}
		return h.Sum(nil)
	}
	return nil
}
