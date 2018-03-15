package function

type Module int

func (t *Module) Do(args *string, reply *string) error {
	*reply = "module real do"
	return nil
}
