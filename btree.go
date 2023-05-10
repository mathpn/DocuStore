package main

import (
	"fmt"
	"io"
)

type InvertedIndex struct {
	Token  string
	DocIDs []string
}

func NewInvIndex(token string, docIDs ...string) *InvertedIndex {
	return &InvertedIndex{
		token,
		docIDs,
	}
}

type docToken struct {
	docID string
	token string
}

type BinaryNode struct {
	Left  *BinaryNode
	Right *BinaryNode
	Data  *InvertedIndex
}

type BinaryTree struct {
	Root *BinaryNode
}

func (t *BinaryTree) InsertDoc(doc *DocSummary) {
	for token := range doc.TermFreqs {
		docToken := &docToken{doc.DocID, token}
		t.insert(docToken)
	}
}

func (t *BinaryTree) insert(Data *docToken) *BinaryTree {
	if t.Root == nil {
		t.Root = &BinaryNode{Data: NewInvIndex(Data.token, Data.docID), Left: nil, Right: nil}
	} else {
		t.Root.insert(Data)
	}
	return t
}

func (n *BinaryNode) insert(Data *docToken) {
	if n == nil {
		return
	} else if Data.token < n.Data.Token {
		if n.Left == nil {
			n.Left = &BinaryNode{Data: NewInvIndex(Data.token, Data.docID), Left: nil, Right: nil}
		} else {
			n.Left.insert(Data)
		}
	} else if Data.token > n.Data.Token {
		if n.Right == nil {
			n.Right = &BinaryNode{Data: NewInvIndex(Data.token, Data.docID), Left: nil, Right: nil}
		} else {
			n.Right.insert(Data)
		}
	} else {
		n.Data.DocIDs = append(n.Data.DocIDs, Data.docID)
	}
}

func (t *BinaryTree) Search(token string) *InvertedIndex {
	if t.Root.Data.Token == token {
		return t.Root.Data
	}
	return t.Root.search(token)
}

func (n *BinaryNode) search(token string) *InvertedIndex {
	if n.Data.Token == token {
		return n.Data
	}
	if token <= n.Data.Token {
		if n.Left == nil {
			return &InvertedIndex{
				token,
				make([]string, 0),
			}
		} else {
			return n.Left.search(token)
		}
	} else {
		if n.Right == nil {
			return &InvertedIndex{
				token,
				make([]string, 0),
			}
		} else {
			return n.Right.search(token)
		}
	}
}

func PrintTree(w io.Writer, node *BinaryNode, ns int, ch rune) {
	if node == nil {
		return
	}

	for i := 0; i < ns; i++ {
		fmt.Fprint(w, " ")
	}
	fmt.Fprintf(w, "%c:%v\n", ch, node.Data)
	PrintTree(w, node.Left, ns+2, 'L')
	PrintTree(w, node.Right, ns+2, 'R')
}
