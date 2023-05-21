package main

import (
	"fmt"
	"io"
	"sync"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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
	Left   *BinaryNode
	Right  *BinaryNode
	Data   *InvertedIndex
	height int
}

func (n *BinaryNode) Height() int {
	if n == nil {
		return 0
	}
	return n.height
}

func (n *BinaryNode) Bal() int {
	return n.Right.Height() - n.Left.Height()
}

func (n *BinaryNode) insert(Data *docToken) *BinaryNode {
	if n == nil {
		return &BinaryNode{
			Data:   NewInvIndex(Data.token, Data.docID),
			height: 1,
		}
	}
	if Data.token == n.Data.Token {
		n.Data.DocIDs = append(n.Data.DocIDs, Data.docID)
		return n
	}
	if Data.token < n.Data.Token {
		n.Left = n.Left.insert(Data)
	} else {
		n.Right = n.Right.insert(Data)
	}
	n.height = max(n.Left.Height(), n.Right.Height()) + 1
	return n.rebalance()
}

func (n *BinaryNode) rotateLeft() *BinaryNode {
	r := n.Right
	n.Right = r.Left
	r.Left = n

	n.height = max(n.Left.Height(), n.Right.Height()) + 1
	r.height = max(r.Left.Height(), r.Right.Height()) + 1
	return r
}

func (n *BinaryNode) rotateRight() *BinaryNode {
	l := n.Left
	n.Left = l.Right
	l.Right = n
	n.height = max(n.Left.Height(), n.Right.Height()) + 1
	l.height = max(l.Left.Height(), l.Right.Height()) + 1
	return l
}

func (n *BinaryNode) rotateRightLeft() *BinaryNode {
	n.Right = n.Right.rotateRight()
	n = n.rotateLeft()
	n.height = max(n.Left.Height(), n.Right.Height()) + 1
	return n
}

func (n *BinaryNode) rotateLeftRight() *BinaryNode {
	n.Left = n.Left.rotateLeft()
	n = n.rotateRight()
	n.height = max(n.Left.Height(), n.Right.Height()) + 1
	return n
}

func (n *BinaryNode) rebalance() *BinaryNode {
	switch {
	case n.Bal() < -1 && n.Left.Bal() == -1:
		return n.rotateRight()
	case n.Bal() > 1 && n.Right.Bal() == 1:
		return n.rotateLeft()
	case n.Bal() < -1 && n.Left.Bal() == 1:
		return n.rotateLeftRight()
	case n.Bal() > 1 && n.Right.Bal() == -1:
		return n.rotateRightLeft()
	}
	return n
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

// tree rebalancing code adapted from https://appliedgo.net/balancedtree/

func (t *BinaryTree) insert(Data *docToken) {
	if t.Root == nil {
		t.Root = &BinaryNode{
			Data: NewInvIndex(Data.token, Data.docID),
		}
	} else {
		t.Root = t.Root.insert(Data)
	}
	if t.Root.Bal() < -1 || t.Root.Bal() > 1 {
		t.rebalance()
	}
}

func (t *BinaryTree) rebalance() {
	if t == nil || t.Root == nil {
		return
	}
	t.Root = t.Root.rebalance()
}

func (t *BinaryTree) Search(token string) *InvertedIndex {
	if t.Root == nil {
		return nil
	}
	if t.Root.Data.Token == token {
		return t.Root.Data
	}
	return t.Root.search(token)
}

func (t *BinaryTree) SearchTokens(tokens []string) []string {
	results := make([]*InvertedIndex, len(tokens))
	var wg sync.WaitGroup
	for i := 0; i < len(tokens); i++ {
		wg.Add(1)
		go func(i int, token string) {
			defer wg.Done()
			res := t.Search(token)
			if res != nil {
				results[i] = res
			}
		}(i, tokens[i])
	}
	wg.Wait()
	docMap := make(map[string]bool)
	out := make([]string, 0)
	for i := 0; i < len(results); i++ {
		invInd := results[i]
		if invInd != nil {
			for _, docID := range invInd.DocIDs {
				if _, val := docMap[docID]; !val {
					docMap[docID] = true
					out = append(out, docID)
				}
			}
		}
	}
	return out
}

func (n *BinaryNode) search(token string) *InvertedIndex {
	if n.Data.Token == token {
		return n.Data
	}
	if token <= n.Data.Token {
		if n.Left == nil {
			return nil
		}
		return n.Left.search(token)
	}
	if n.Right == nil {
		return nil
	}
	return n.Right.search(token)
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
