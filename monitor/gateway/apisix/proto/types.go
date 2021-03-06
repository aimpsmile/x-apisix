package proto

type Value struct {
	Content string `json:"content"`
}

//	节点
type Node struct {
	Value *Value `json:"value"`
	Key   string `json:"key"`
}

//	节点列表
type Nodes struct {
	Nodes []*Node `json:"nodes"`
}

//	写入列表
type OneRsp struct {
	Node *Node `json:"node"`
}

//	读取列表
type ListRsp struct {
	Node   *Nodes `json:"node"`
	Action string `json:"action"`
}
