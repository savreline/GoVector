package govec

import (
    "log"
    "net/rpc/jsonrpc"
    "net/rpc"
    "encoding/json"
    "net"
    "./../server/broker"
    "./../server/broker/nonce"
    "github.com/arcaneiceman/GoVector/govec/vclock"
)

// ******************************************
// Client side RPC Library Calls
// 

type GoPublisher struct {
    conn      net.Conn
    address   string
    rpcconn   *rpc.Client
    nonce     *nonce.Nonce
}

// TCP connection receives the nonce to use when creating the connection to broker
func NewGoPublisher(addr string, port string) *GoPublisher {
    url := addr + ":" + port
    log.Println("GoPublisher: Connecting to url: ", url)
    tcps, err := net.Dial("tcp", url)
    
    if err != nil {
        log.Fatal("GoPublisher tcp error: ", err)
    }
    
    var nonce nonce.Nonce
    
    d := json.NewDecoder(tcps)
    tcperr := d.Decode(&nonce)
    log.Println("GoPublisher: Received nonce: ", nonce)
    if tcperr != nil {
        log.Fatal("GoPublisher tcp error: ", tcperr)
    }
    
    jrpc := jsonrpc.NewClient(tcps)

    gp := &GoPublisher{
        conn: tcps,
        address: addr,
        rpcconn: jrpc,
        nonce: &nonce,
    }
    return gp
}

func (gp *GoPublisher) PublishLocalMessage(msg string, processID string, vclock vclock.VClock) error {
    message := brokervec.LocalMessage{
        Pid: processID, 
        Vclock: vclock.Bytes(),
        Message: msg,
        Nonce: gp.nonce.Nonce}
    var reply string
    err := gp.rpcconn.Call("PubManager.AddLocalMsg", message, &reply)
    
    if err != nil {
        log.Println("GoPublisher: PubMgr error: ", err)
        return err
    } else {
        log.Println("GoPublisher: Sent message, reply was: ", reply)
    }
    return nil
}

func (gp *GoPublisher) PublishNetworkMessage(msg string, processID string, vclock vclock.VClock) error {
    
    message := brokervec.NetworkMessage{
        Pid: processID, 
        Vclock: vclock.Bytes(),
        Message: msg,
        Nonce: gp.nonce.Nonce}
    var reply string
    err := gp.rpcconn.Call("PubManager.AddNetworkMsg", message, &reply)
    
    if err != nil {

        log.Println("GoPublisher: PubManager error: ", err)
        return err
    }
    
    return nil
}