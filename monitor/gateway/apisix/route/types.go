package route

type Value struct {
	Desc            string `json:"desc"`
	ServiceID       string `json:"service_id"`
	UpstreamID      string `json:"upstream_id"`
	ServiceProtocol string `json:"service_protocol"`
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
