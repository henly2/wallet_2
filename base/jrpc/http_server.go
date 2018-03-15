package jrpc

import (
	"net/http"
	"io"
	"log"
	"bytes"
	"net/rpc/jsonrpc"
)

// rpcRequest represents a RPC request.
// rpcRequest implements the io.ReadWriteCloser interface.
type rpcRequest struct {
	r    io.Reader     // holds the JSON formated RPC request
	rw   io.ReadWriter // holds the JSON formated RPC response
	done chan bool     // signals then end of the RPC request
}

// Read implements the io.ReadWriteCloser Read method.
func (r *rpcRequest) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

// Write implements the io.ReadWriteCloser Write method.
func (r *rpcRequest) Write(p []byte) (n int, err error) {
	return r.rw.Write(p)
}

// Close implements the io.ReadWriteCloser Close method.
func (r *rpcRequest) Close() error {
	r.done <- true
	return nil
}

// Call invokes the RPC request, waits for it to complete, and returns the results.
func (r *rpcRequest) Call() io.Reader {
	go jsonrpc.ServeConn(r)
	<-r.done
	return r.rw
}

// NewRPCRequest returns a new rpcRequest.
func NewRPCRequest(r io.Reader) *rpcRequest {
	var buf bytes.Buffer
	done := make(chan bool)
	return &rpcRequest{r, &buf, done}
}

// Start a JRPC Http Server
// @parameter: port string, like ":8080"
// @return: error
func StartJRPCHttpServer(port string) error {
	log.Println("Start JRPC Http server...", port)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("JRPC Http server Accept a client: ", req.RemoteAddr)

		defer req.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		res := NewRPCRequest(req.Body).Call()
		io.Copy(w, res)
	})

	log.Println("Start JRPC Http server successfully, listen on port: ", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}

	return nil
}