package service

//	service对应的负载均衡
type Upstream struct {
	Nodes map[string]int `json:"nodes"`
}

type Value struct {
	Desc     string    `json:"desc"`
	Upstream *Upstream `json:"upstream"`
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
