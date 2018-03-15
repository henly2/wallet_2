package method

import(
	"log"
	"errors"
)

// 服务中心
type ServiceCenterInterface interface {
	HandleRegister(*string, *string)error
	HandleDispatch(*string, *string)error
}

type ServiceCenter struct{
	Instance ServiceCenterInterface
}

// 服务节点
type ServiceNodeInterface interface {
	HandleCall(*string, *string)error
}

type ServiceNode struct{
	Instance ServiceNodeInterface
}

// 服务中心方法
func (c *ServiceCenter) Register(req *string, res * string) error {
	log.Println("ServiceCenter register: ", *req)

	if c.Instance == nil {
		log.Println("ServiceCenter interface is nil")
		return errors.New("ServiceCenter interface is nil")
	}
	c.Instance.HandleRegister(req, res)
	return nil;
}

func (c *ServiceCenter) Dispatch(req *string, res * string) error {
	log.Println("ServiceCenter dispatch : ", *req)

	if c.Instance == nil {
		log.Println("ServiceCenter interface is nil")
		return errors.New("ServiceCenter interface is nil")
	}
	c.Instance.HandleDispatch(req, res)
	return nil;
}

// 服务节点方法
func (n *ServiceNode) Call(req *string, res * string) error {
	log.Println("ServiceNode call: ", *req)

	if n.Instance == nil {
		log.Println("ServiceNode interface is nil")
		return errors.New("ServiceNode interface is nil")
	}
	n.Instance.HandleCall(req, res)
	return nil;
}
