package business

import "fmt"

func HandleMsg(args* string, reply *string) error {
	switch *args {
	case "new_address":
		fmt.Println(*args)
		*reply = *args
	case "withdrawal":
		*reply = *args
	}
	return nil
}
