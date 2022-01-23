package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/netio"
)

//--- handlers ---

//alternative
// m := map[string]interface{}{
// 	"handlef": f,
// }
// 	switch fcall {
// 	case "f":
// 		func()
// 	}
// }

func HandleEcho(ins string) string {
	resp := "Echo:" + ins
	return resp
}

func HandlePing(msg netio.Message) netio.Message {
	// validRequest := msg.MessageType == netio.REQ && msg.Command == "PING"
	// if !validRequest{
	// 	//error
	// }
	reply_msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_PONG}
	return reply_msg
}

func HandleBlockheight(t *TCPNode, msg netio.Message) netio.Message {
	bh := len(t.Mgr.Blocks)
	//data := strconv.Itoa(bh)
	djson, _ := json.Marshal(bh)
	reply_msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_BLOCKHEIGHT, Data: []byte(djson)}
	//reply_msg := netio.EdnConstructMsgMapData(netio.REP, netio.CMD_BLOCKHEIGHT, data)
	//("BLOCKHEIGHT ", reply_msg)
	return reply_msg
}

//Standard Tx handler
//InteractiveTx also possible
//client requests tranaction <=> server response with challenge <=> client proves
func HandleTx(t *TCPNode, msg netio.Message) netio.Message {
	dataBytes := msg.Data

	var tx block.Tx

	if err := json.Unmarshal(dataBytes, &tx); err != nil {
		panic(err)
	}
	t.log(fmt.Sprintf("tx %v ", tx))

	reply := chain.HandleTx(t.Mgr, tx)
	return reply
}

func HandleBalance(t *TCPNode, msg netio.Message) string {
	dataBytes := msg.Data
	t.log(fmt.Sprintf("HandleBalance data %v %v", string(msg.Data), dataBytes))
	fmt.Println(fmt.Sprintf("HandleBalance data %v %v", string(msg.Data), dataBytes))

	//a := block.Account{AccountKey: string(msg.Data)}

	// var account block.Account

	// if err := json.Unmarshal(dataBytes, &account); err != nil {
	// 	panic(err)
	// }

	a := string(msg.Data)
	balance := t.Mgr.State.Accounts[a]
	fmt.Println("balance for ", a, balance, t.Mgr.State.Accounts)

	balJson, _ := json.Marshal(balance)
	reply_msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_BALANCE, Data: []byte(balJson)}
	reply_msg_json := netio.ToJSONMessage(reply_msg)
	//reply_msg := netio.EdnConstructMsgMapData(netio.REP, netio.CMD_BALANCE, data)

	return reply_msg_json
}

func HandlePeers(t *TCPNode, msg netio.Message) netio.Message {
	peersJson, _ := json.Marshal(t.Peers)
	reply_msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_GETPEERS, Data: []byte(peersJson)}
	return reply_msg
}

func HandleFaucet(t *TCPNode, msg netio.Message) netio.Message {
	//t.log(fmt.Sprintf("HandleFaucet"))

	// dataBytes := msg.Data
	// var account block.Account
	// if err := json.Unmarshal(dataBytes, &account); err != nil {
	// 	panic(err)
	// }

	//account := block.Account{AccountKey: string(msg.Data)}
	//t.log(fmt.Sprintf("faucet for ... %v", account.AccountKey))

	randNonce := 0
	amount := rand.Intn(10)

	keypair := chain.GenesisKeys()
	addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))
	//Genesis_Account := block.AccountFromString(addr)

	//tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: Genesis_Account, Receiver: account}
	a := string(msg.Data)
	tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: addr, Receiver: a}

	tx = crypto.SignTxAdd(tx, keypair)
	reply := chain.HandleTx(t.Mgr, tx)
	//t.log(fmt.Sprintf("resp > %s", reply_string))

	//reply := netio.EdnConstructMsgMapData(netio.REP, netio.CMD_FAUCET, reply_string)
	return reply
}

func HandleRegistername(t *TCPNode, peer *netio.Peer, msg netio.Message) netio.Message {
	fmt.Println("HandleRegistername")
	//TODO!
	data, _ := json.Marshal("ok")
	//set name
	newname := string(msg.Data)
	fmt.Println("set peer name to ", newname)
	peer.Name = newname
	fmt.Println("new name ", peer.Name)
	reply_msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_REGISTERNAME, Data: []byte(data)}
	//fmt.Println("name ", ntchan.Name)
	return reply_msg
}

func RequestReply(ntchan Ntchan, msg netio.Message) string {

	var reply_msg string
	//var reply_msg netio.Message

	fmt.Sprintf("Handle cmd %v", msg.Command)

	switch msg.Command {

	case netio.CMD_PING:
		reply_msg := "pong"
		return reply_msg
		//reply := HandlePing(msg)
		//msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_BALANCE, Data: []byte(balJson)}
		//reply_msg = netio.ToJSONMessage(reply)

	case netio.CMD_TIME:
		dt := time.Now()
		reply_msg := dt.String()
		return reply_msg

	case netio.CMD_REGISTERPEER:
		reply_msg := "todo"
		return reply_msg

	default:
		errormsg := "Error: not found command"
		fmt.Println(errormsg)
		xjson, _ := json.Marshal("")
		msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_ERROR, Data: []byte(xjson)}
		reply_msg = netio.ToJSONMessage(msg)
	}

	return reply_msg
}
