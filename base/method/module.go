package method

import(
	"log"
	"errors"
)

type ModuleInterface interface {
	HandleRegister(*string, *string)error
	HandleDispatch(*string, *string)error
	HandleCall(*string, *string)error
}

type Module struct{
	Instance ModuleInterface
}

func (m *Module) Register(req *string, res * string) error {
	log.Println("Module register: ", *req)

	if m.Instance == nil {
		log.Println("Module interface is nil")
		return errors.New("Module interface is nil")
	}
	m.Instance.HandleRegister(req, res)
	return nil;
}

func (m *Module) Dispatch(req *string, res * string) error {
	log.Println("Module dispath : ", *req)

	if m.Instance == nil {
		log.Println("Module interface is nil")
		return errors.New("Module interface is nil")
	}
	m.Instance.HandleDispatch(req, res)
	return nil;
}

func (m *Module) Call(req *string, res * string) error {
	log.Println("Module call : ", *req)

	if m.Instance == nil {
		log.Println("Module interface is nil")
		return errors.New("Module interface is nil")
	}
	m.Instance.HandleCall(req, res)
	return nil;
}