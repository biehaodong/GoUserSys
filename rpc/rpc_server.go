package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
)

const TcpHeadMaxSize int = 4

//use map struct to reflect the handle name and handle function
type Server struct {
	Se map[string]Handler
}

type serveFunc func(interface{}) interface{}

//request struct
type Handler struct {
	intfc   serveFunc    //request func
	reqType reflect.Type //request type
	resType reflect.Type //response type
}

//client request struct
type Request struct {
	ReqName string `json:"reqName"` //client request name
	ReqData []byte `json:"reqDate"` //client request data
}

//init Server
func NewServer() Server {
	return Server{make(map[string]Handler)}
}

//rpc register function, get handle by handler ,get actual args type by serv.
func (s *Server) Register(name string, handler serveFunc, serv interface{}) error {
	//get type of service by reflect
	servType := reflect.TypeOf(serv)
	//check the request type
	if err := s.checkType(servType); err != nil {
		return err
	}
	//get type of args and response
	reqType := servType.In(0)
	resType := servType.Out(0)
	//save [name,Handle] map
	s.Se[name] = Handler{intfc: handler, reqType: reqType, resType: resType}
	return nil
}

//rpc listen function, listen the address
func (s *Server) Listen(address string) (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}
	listen, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		return nil, err
	}
	return listen, nil
}

//serve function
func (s *Server) Serve(l *net.TCPListener) error {
	if err := s.accept(l); err != nil {
		return err
	}
	return nil
}

//read the data from rpc client, handle and return
func (s *Server) Handle(conn *net.TCPConn) error {
	if conn == nil {
		return errors.New("rpc_server:connect is null")
	}
	dLen := make([]byte, TcpHeadMaxSize)
	for {
		n, err := conn.Read(dLen)
		if err != io.EOF && err != nil {
			errors.New("rpc_server:read connect failed")
		}
		if n <= 0 {
			errors.New("rpc_server:no data")
		}
		l, err := strconv.ParseInt(string(dLen[:TcpHeadMaxSize]), 10, 64)
		if err != nil {
			return err
		}
		//get the data of l length pack
		b := make([]byte, l)
		n, err = conn.Read(b)
		if err != nil {
			return errors.New("rpc_server: connection read body failed")
		}
		if n <= 0 {
			return errors.New("rpc_server: no data")
		}
		//handle the request
		resp, err := s.find(b)
		if err != nil {
			return err
		}
		//pack the response data to the json format
		respJ, err := s.ToJson(resp)
		if err != nil {
			return err
		}
		conn.Write(respJ) //send the result
	}
}

//find the Handle interface to handle the request by its name
func (s *Server) find(data []byte) (interface{}, error) {
	//parse the interface name
	var req Request
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}
	//get the Handler by req name
	r, ok := s.Se[req.ReqName]
	if !ok {
		return nil, errors.New("rpc_server:don't find handler")
	}
	//pares data type, get data by this type. save to reqType, reqType is Handle args
	reqType := reflect.New(r.reqType).Interface()
	err = json.Unmarshal(req.ReqData, reqType)
	if err != nil {
		return nil, err
	}
	return r.intfc(reqType), nil
}
func (s *Server) ToJson(f interface{}) ([]byte, error) {
	return toJsons(f)
}

//pack the interface to json format and return
func toJsons(f interface{}) ([]byte, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	l := len(b)
	//convert int to string
	lStr := strconv.Itoa(l)

	iL := len(lStr)
	//l := len(strconv.Itoa(len(b)))
	//l size,compare with teh TcpHeadMaxSize
	if iL > TcpHeadMaxSize {
		return nil, errors.New("rpc_server: header is too big")
		fmt.Println("rpc_server: header is too big")
	}
	iL = iL - 1
	tb := []byte("0")
	j := make([]byte, TcpHeadMaxSize+l)
	//request head
	for i := TcpHeadMaxSize - 1; i >= 0; i-- {
		if iL < 0 {
			j[i] = tb[0]
		} else {
			j[i] = []byte(lStr)[iL]
			iL = iL - 1
		}
	}
	//request body
	for i := 0; i < l; i++ {
		j[TcpHeadMaxSize+i] = b[i]
	}
	return j, nil
}

//accept new request
func (s *Server) accept(listen *net.TCPListener) error {
	defer listen.Close()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			return err
		}
		defer conn.Close()
		go s.Handle(conn)
	}
}

//check handler type
func (s *Server) checkType(handlerType reflect.Type) error {
	// whether is func type
	if handlerType.Kind() != reflect.Func {
		return errors.New("rpc.Register: handler is not func")
	}
	// parameter amount.
	if handlerType.NumIn() != 1 {
		return errors.New("rpc.Register: handler input parameters number is wrong, need one")
	}
	// response data amount.
	if handlerType.NumOut() != 1 {
		return errors.New("rpc.Register: handler output parameters number is wrong, need one")
	}
	// judge the parameter type and response data type.
	if handlerType.In(0).Kind() != reflect.Struct || handlerType.Out(0).Kind() != reflect.Struct {
		return errors.New("rpc.Register: parameters must be Struct")
	}
	return nil
}

//shutdown
func (s *Server) Shutdown(listen *net.TCPListener) {
	defer listen.Close()
}
