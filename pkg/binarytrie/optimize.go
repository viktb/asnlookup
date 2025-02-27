package binarytrie

// Optimize performs level compression and path compression on the trie.
//
// This operation makes the trie immutable.
func (t *NaiveTrie) Optimize(fillFactor float32) error {
	if !t.mutable {
		return ErrTrieImmutable
	}

	t.root.propagateValues(0)
	t.root.removeRedundancies()
	t.root.compressLevels(fillFactor)
	t.root.compressPaths(nil, 0, 0)
	t.mutable = false
	return nil
}

func (n *naiveTrieNode) propagateValues(value uint32) {
	for i, child := range n.children {
		if child == nil {
			n.children[i] = &naiveTrieNode{parent: n}
		}
	}
	if n.value == 0 {
		n.value = value
	}

	// Preorder traversal.
	for _, child := range n.children {
		child.propagateValues(n.value)
	}
}

func (n *naiveTrieNode) removeRedundancies() {
	// Postorder traversal.
	for _, child := range n.children {
		child.removeRedundancies()
	}

	if n.isLeaf() {
		return
	}
	if n.children[0].isLeaf() && n.children[1].isLeaf() && n.children[0].value == n.children[1].value {
		n.value = n.children[0].value
		n.branchingFactor = 0
		n.children = nil
	}
}

func (n *naiveTrieNode) compressPaths(firstNode *naiveTrieNode, depth uint8, prefix uint32) {
	var nextNode *naiveTrieNode
	var nextNodeIndex int
	for i, child := range n.children {
		if child.isLeaf() {
			if child.value == n.value {
				// Ignore trivial leaves introduced by normalization.
				continue
			}
			// Any other leaf breaks the path.
			nextNode = nil
			break
		}
		if nextNode != nil {
			// The path ends if there is more than 1 nontrivial child.
			nextNode = nil
			break
		}
		nextNode = child
		nextNodeIndex = i
	}

	if firstNode != nil && (nextNode == nil || depth+n.branchingFactor > 32) {
		// The path ends.
		firstNode.skipValue = depth
		firstNode.children = n.children
		firstNode.branchingFactor = n.branchingFactor
		for _, child := range n.children {
			child.skippedBits = prefix
			child.parent = firstNode
		}
	} else if firstNode != nil {
		// The path continues.
		nextNode.compressPaths(firstNode, depth+n.branchingFactor, (prefix<<n.branchingFactor)|uint32(nextNodeIndex))
		return
	} else if nextNode != nil {
		// The path begins.
		nextNode.compressPaths(n, n.branchingFactor, uint32(nextNodeIndex))
		return
	}

	for _, child := range n.children {
		child.compressPaths(nil, 0, 0)
	}
}

func (n *naiveTrieNode) compressLevels(fillFactor float32) {
	fakeNodes := make(map[*naiveTrieNode]bool)

	depth := uint8(0)
	nodes := []*naiveTrieNode{n}
	for {
		nextNodes := make([]*naiveTrieNode, 0, 1<<(depth+1))
		for _, node := range nodes {
			if !node.isLeaf() {
				nextNodes = append(nextNodes, node.children...)
				continue
			}

			fakeChildren := []*naiveTrieNode{
				{parent: node, value: node.value},
				{parent: node, value: node.value},
			}
			nextNodes = append(nextNodes, fakeChildren...)
			fakeNodes[fakeChildren[0]] = true
			fakeNodes[fakeChildren[1]] = true
		}

		realNodeCount := 0
		for _, node := range nextNodes {
			if !fakeNodes[node] {
				realNodeCount++
			}
		}
		if depth >= 31 || float32(realNodeCount)/float32(len(nextNodes)) < fillFactor {
			break
		}
		depth++
		nodes = nextNodes
	}

	if depth > 1 {
		n.branchingFactor = depth
		n.children = nodes
	}

	for _, child := range n.children {
		child.compressLevels(fillFactor)
	}
}
