package common

const(
	Method_Module_Register = "Module.Register"
	Method_Module_Dispatch = "Module.Dispatch"
	Method_Module_Call = "Module.Call"
)

type ModuleRegisterData struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type ModuleDispatchData struct{
	Name string `json:"name"`
}