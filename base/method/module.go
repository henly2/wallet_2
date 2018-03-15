package method

import(
	"log"
	"errors"
)

// 服务中心
type ServerCenterInterface interface {
	HandleRegister(*string, *string)error
	HandleDispatch(*string, *string)error
}

type ServerCenter struct{
	Instance ServerCenterInterface
}

// 服务节点
type ServerNodeInterface interface {
	HandleCall(*string, *string)error
}

type ServerNode struct{
	Instance ServerNodeInterface
}

// 服务中心方法
func (c *ServerCenter) Register(req *string, res * string) error {
	log.Println("ServerCenter register: ", *req)

	if c.Instance == nil {
		log.Println("ServerCenter interface is nil")
		return errors.New("ServerCenter interface is nil")
	}
	c.Instance.HandleRegister(req, res)
	return nil;
}

func (c *ServerCenter) Dispatch(req *string, res * string) error {
	log.Println("ServerCenter dispath : ", *req)

	if c.Instance == nil {
		log.Println("ServerCenter interface is nil")
		return errors.New("ServerCenter interface is nil")
	}
	c.Instance.HandleDispatch(req, res)
	return nil;
}

// 服务节点方法
func (n *ServerNode) Call(req *string, res * string) error {
	log.Println("ServerNode call: ", *req)

	if n.Instance == nil {
		log.Println("ServerNode interface is nil")
		return errors.New("ServerNode interface is nil")
	}
	n.Instance.HandleCall(req, res)
	return nil;
}
