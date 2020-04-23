package sortedset

type SortedSetLevel struct {
	forward *SortedSetNode
	span    int64
}

// Node in skip list
type SortedSetNode struct {
	key      string      // unique key of this node
	Value    interface{} // associated data
	score    SCORE       // score to determine the order of this node in the set
	backward *SortedSetNode
	level    []SortedSetLevel
}

// Get the key of the node
func (this *SortedSetNode) Key() string {
	return this.key
}

// Get the node of the node
func (this *SortedSetNode) Score() SCORE {
	return this.score
}
