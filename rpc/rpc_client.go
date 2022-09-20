package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"

	"GoUserManaSys/utils"
)

//rpc client,contains pool
type Client struct {
	pool chan net.TCPConn
}

//client creat n connections connect to address
func NewClient(n int, add string) (Client, error) {
	//parse tcp server address
	tAddr, err := net.ResolveTCPAddr("tcp4", add)
	if err != nil {
		return Client{}, err
	}
	//creat connections and save to pool
	newPool := make(chan net.TCPConn, n)
	for i := 0; i < n; i++ {
		connect, err := net.DialTCP("tcp4", nil, tAddr)
		if err != nil {
			return Client{}, errors.New("rpc_client:init failed")
		}
		newPool <- *connect
	}
	//fmt.Printf("rpc_client:init success...")
	return Client{pool: newPool}, nil
}

// close pool
func (c *Client) Close() {
	close(c.pool)
	for cn := range c.pool {
		cn.Close()
		//fmt.Println("close conn")
	}
}

//caller use Call function to connect server
func (c *Client) Call(name string, req interface{}, res interface{}) error {
	//get an idle conn
	connect := c.getC()
	//put conn to pool
	defer c.putC(connect)
	//pack the request
	reqJ, err := c.packReq(name, req)
	if err != nil {
		return err
	}
	//send data to server
	//handle err
	connect.Write(reqJ)

	dL := make([]byte, TcpHeadMaxSize)
	n, err := connect.Read(dL)
	if err != io.EOF && err != nil {
		return err
	}
	if n <= 0 {
		return errors.New("rpc_client: no data")
	}
	//convert dL to 64 bitSize, save to len
	l, err := strconv.ParseInt(string(dL[:TcpHeadMaxSize]), 10, 64)
	if err != nil {
		return nil
	}
	//make b size buffer to accept response data
	b := make([]byte, l)
	//read data and save to b
	n, err = connect.Read(b)
	if err != io.EOF && err != nil {
		return err
	}
	if n <= 0 {
		return errors.New("rpc_client: no data")
	}
	//parse json data and save to res
	err = c.unpackRes(res, b)
	if err != nil {
		return err
	}
	return nil
}

//get an idle conn and return
func (c *Client) getC() (connect net.TCPConn) {
	select {
	case connect := <-c.pool:
		return connect
	}
}

//put the conn to pool
func (c *Client) putC(Connect net.TCPConn) {
	select {
	case c.pool <- Connect:
		return
	}
}

//pack the request
func (c *Client) packReq(name string, f interface{}) ([]byte, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	req := Request{ReqName: name, ReqData: data}
	request, err := toJsons(req)
	if err != nil {
		return nil, err
	}
	return request, nil
}

//unpack the response data
func (c *Client) unpackRes(resI interface{}, resB []byte) error {
	err := json.Unmarshal(resB, resI)
	if err != nil {
		return err
	}
	return nil
}

//rcp client re-call
func (c *Client) ReCall(name string, req interface{}, res interface{}) error {
	fmt.Println("Disconnected, re-establishing connection...")
	//close pool
	c.Close()
	//re-creat rpc client
	cli, err := NewClient(utils.ClientPoolSize, utils.ServerPort)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully re-established connection...")
	err = cli.Call(name, req, &res)
	return err
}
