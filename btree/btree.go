/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package btree

import "os"

type BNode struct {
	children []*BNode
	ptr      *os.File
}

type BTree BNode

func NewBTree() *BTree {
	return &BTree{}
}

func (bt *BTree) Create(node *BNode) {

}

