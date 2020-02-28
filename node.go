package main

//kill -9 $(lsof -t -i:8888)
//node should run via DNS
//nodexample.com

//basic protocol
//node receives tx messages
//adds tx messages to a pool
//block gets created every 10 secs

//getBlocks
//registerPeer
//pickRandomAccount
//storeBalance

//newWallet

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/ntwk"
)

var Peers []ntwk.Peer

//banned IPs
var nlog *log.Logger
var logfile_name = "node.log"

var blockTime = 10000 * time.Millisecond

var tchan chan string

type Configuration struct {
	PeerAddresses []string
	NodePort      int
	WebPort       int
}

//inbound
func addpeer(addr string, nodeport int) ntwk.Peer {
	//ignored

	p := ntwk.CreatePeer(addr, nodeport)
	Peers = append(Peers, p)
	nlog.Println("peers ", Peers)
	return p
}

//inbound
func setupPeer(addr string, nodeport int, conn net.Conn) {
	peer := addpeer(addr, nodeport)

	nlog.Println("setup channels for incoming requests")
	//TODO peers chan
	//TODO handshake
	go channelPeerNetwork(conn, peer)
}

// start listening on tcp and handle connection through channels
func ListenAll(node_port int) error {
	nlog.Println("listen all")
	var err error
	var listener net.Listener
	p := ":" + strconv.Itoa(node_port)
	nlog.Println("listen on ", p)
	listener, err = net.Listen("tcp", p)
	if err != nil {
		nlog.Println(err)
		return errors.Wrapf(err, "Unable to listen on port %d\n", node_port) //ntwk.Port
	}

	addr := listener.Addr().String()
	nlog.Println("Listen on", addr)

	//TODO check if peers are alive see
	//https://stackoverflow.com/questions/12741386/how-to-know-tcp-connection-is-closed-in-net-package
	//https://gist.github.com/elico/3eecebd87d4bc714c94066a1783d4c9c

	for {
		nlog.Println("Accept a connection request")

		//TODO peer handshake
		//TODO client handshake

		conn, err := listener.Accept()
		strRemoteAddr := conn.RemoteAddr().String()

		nlog.Println("accepted conn ", strRemoteAddr)
		if err != nil {
			nlog.Println("Failed accepting a connection request:", err)
			continue
		}

		setupPeer(strRemoteAddr, node_port, conn)

	}
}

func putMsg(msg_in_chan chan ntwk.Message, msg ntwk.Message) {
	msg_in_chan <- msg
}

func HandlePing() ntwk.Message {
	reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_PONG, ntwk.EMPTY_DATA)
	return reply
	//msg_out_chan <- reply
}

func HandleHandshake(msg_out_chan chan ntwk.Message) {
	reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_HANDSHAKE_STABLE, ntwk.EMPTY_DATA)
	msg_out_chan <- reply
}

func HandleReqMsg(msg ntwk.Message) ntwk.Message {
	nlog.Println("Handle ", msg.Command)

	switch msg.Command {

	case ntwk.CMD_PING:
		nlog.Println("PING PONG")
		return HandlePing()

		// case ntwk.CMD_HANDSHAKE_HELLO:
		// 	nlog.Println("handshake")
		// 	HandleHandshake(rep_chan)

		// case ntwk.CMD_BALANCE:
		// 	nlog.Println("Handle balance")

		// 	dataBytes := msg.Data
		// 	nlog.Println("data ", dataBytes)
		// 	var account block.Account

		// 	if err := json.Unmarshal(dataBytes, &account); err != nil {
		// 		panic(err)
		// 	}
		// 	nlog.Println("get balance for account ", account)

		// 	balance := chain.Accounts[account]
		// 	//s := strconv.Itoa(balance)
		// 	data, _ := json.Marshal(balance)
		// 	reply := ntwk.EncodeMsgBytes(ntwk.REP, ntwk.CMD_BALANCE, data)
		// 	log.Println(">> ", reply)

		// 	rep_chan <- reply

		// case ntwk.CMD_FAUCET:
		// 	//send money to specified address

		// 	dataBytes := msg.Data
		// 	var account block.Account
		// 	if err := json.Unmarshal(dataBytes, &account); err != nil {
		// 		panic(err)
		// 	}
		// 	nlog.Println("faucet for ... ", account)

		// 	randNonce := 0
		// 	amount := 10

		// 	keypair := chain.GenesisKeys()
		// 	addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))
		// 	Genesis_Account := block.AccountFromString(addr)

		// 	tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: Genesis_Account, Receiver: account}

		// 	tx = crypto.SignTxAdd(tx, keypair)
		// 	reply_string := chain.HandleTx(tx)
		// 	nlog.Println("resp > ", reply_string)

		// 	reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_FAUCET, reply_string)

		// 	rep_chan <- reply

		// case ntwk.CMD_BLOCKHEIGHT:

		// 	data, _ := json.Marshal(len(chain.Blocks))
		// 	reply := ntwk.EncodeMsgBytes(ntwk.REP, ntwk.CMD_BLOCKHEIGHT, data)
		// 	log.Println("CMD_BLOCKHEIGHT >> ", reply)

		// 	rep_chan <- reply

		// case ntwk.CMD_TX:
		// 	nlog.Println("Handle tx")

		// 	dataBytes := msg.Data

		// 	var tx block.Tx

		// 	if err := json.Unmarshal(dataBytes, &tx); err != nil {
		// 		panic(err)
		// 	}
		// 	nlog.Println(">> ", tx)

		// 	resp := chain.HandleTx(tx)
		// 	msg := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_TX, resp)
		// 	rep_chan <- msg

		// case ntwk.CMD_GETTXPOOL:
		// 	nlog.Println("get tx pool")

		// 	//TODO
		// 	data, _ := json.Marshal(chain.Tx_pool)
		// 	msg := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_GETTXPOOL, string(data))
		// 	rep_chan <- msg

		//var Tx_pool []block.Tx

		// case ntwk.CMD_RANDOM_ACCOUNT:
		// 	nlog.Println("Handle random account")

		// 	txJson, _ := json.Marshal(chain.RandomAccount())

	default:
		// 	nlog.Println("unknown cmd ", msg.Command)
		resp := "ERROR UNKONWN CMD"

		msg := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_TX, resp)
		return msg

	}
}

//handle messages
// func HandleMsg(req_chan chan ntwk.Message, rep_chan chan ntwk.Message) {
// 	req_msg := <-req_chan
// 	//fmt.Println("handle msg string ", msgString)

// 	fmt.Println("msg type ", req_msg.MessageType)

// 	if req_msg.MessageType == ntwk.REQ {
// 		HandleReqMsg(req_msg, rep_chan)
// 	} else if req_msg.MessageType == ntwk.REP {
// 		nlog.Println("handle reply")
// 	}
// }

//old
// func ReplyLoop(ntchan ntwk.Ntchan, req_chan chan ntwk.Message, rep_chan chan ntwk.Message) {

// 	//continously read for requests and respond with reply
// 	for {

// 		// read from network
// 		msgString := ntwk.NetworkReadMessage(ntchan)
// 		if msgString == ntwk.EMPTY_MSG {
// 			//log.Println("empty message, ignore")
// 			time.Sleep(500 * time.Millisecond)
// 			continue
// 		}

// 		msg := ntwk.ParseMessage(msgString)
// 		nlog.Print("Receive message over network ", msgString)

// 		//put in the channel
// 		go putMsg(req_chan, msg)

// 		//handle in channel and put reply in msg_out channel
// 		go HandleMsg(req_chan, rep_chan)

// 		//take from reply channel and send over network
// 		reply := <-rep_chan
// 		fmt.Println("msg out ", reply)
// 		ntwk.ReplyNetwork(ntchan, reply)

// 	}
// }

//old
func ReqLoop(ntchan ntwk.Ntchan, out_req_chan chan ntwk.Message, out_rep_chan chan ntwk.Message) {

	//TODO
	for {
		request := <-out_req_chan
		log.Println("request ", request)
	}
}

//old
func PubLoop(ntchan ntwk.Ntchan, Pub_chan chan ntwk.Message) {
	for {
		pubmsg := <-Pub_chan
		log.Println("pubchan.. todo relay ", pubmsg)
		var msg ntwk.Message
		msg.MessageType = ntwk.PUB
		msg.Command = "test"
		//msg.Data := bytes("")
		ntwk.PubNetwork(ntchan, msg)
	}
}

func pubhearbeat(ntchan ntwk.Ntchan) {
	heartbeat_time := 1000 * time.Millisecond

	for {
		msg := ntwk.EncodeHeartbeat("peer1")
		ntchan.Writer_queue <- msg
		time.Sleep(heartbeat_time)
	}

}

//process requests
func Reqprocessor(ntchan ntwk.Ntchan) {
	msg_string := <-ntchan.REQ_chan_in
	//logmsgd("REQ processor ", x)

	//reply_string := "reply"
	//reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_PONG, ntwk.EMPTY_DATA)
	msg := ntwk.ParseMessage(msg_string)
	reply := HandleReqMsg(msg)

	reply_string := ntwk.MsgString(reply)

	ntchan.Writer_queue <- reply_string
}

func Repprocessor(ntchan ntwk.Ntchan) {

}

//setup the network of channels
//the main junction for managing message flow between types of messages
func channelPeerNetwork(conn net.Conn, peer ntwk.Peer) {

	ntchan := ntwk.ConnNtchan(conn, peer.Address)

	//main reader and writer setup
	go ntwk.ReaderWriterConnector(ntchan)

	go Reqprocessor(ntchan)

	go pubhearbeat(ntchan)

	//Request handler
	// go func() {
	// 	log.Println("handler ")
	// 	msg := <-ntchan.Req_chan
	// 	log.Println(">>> request ", msg)
	// }()

	//could add max listen
	//timeoutDuration := 5 * time.Second
	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	//when close?
	//defer conn.Close()

	//go ReplyLoop(ntchan, peer.Req_chan, peer.Rep_chan)

	//go ReqLoop(ntchan, peer.Out_req_chan, peer.Out_rep_chan)

	//----------------
	//TODO pubsub
	//go publishLoop(msg_in_chan, msg_out_chan)

	// go ntwk.Subtime(tchan, "peer1")

	// go ntwk.Subout(tchan, "peer1", peer.Pub_chan)

	// go PubLoop(rw, peer.Pub_chan)

}

//channel network optimised for client
func channelPeerNetworkClient(conn net.Conn, peer ntwk.Peer) {

	ntchan := ntwk.ConnNtchan(conn, peer.Address)

	ntwk.ReaderWriterConnector(ntchan)

}

//basic threading helper
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func doEveryX(d time.Duration, f func() <-chan string) {
	for _ = range time.Tick(d) {
		f()
	}
}

func connect_peers(node_port int, PeerAddresses []string) {

	//TODO

	for _, peer := range PeerAddresses {
		addr := peer + strconv.Itoa(node_port)

		ntchan := ntwk.OpenNtchan(addr)

		out_req := make(chan ntwk.Message)
		out_rep := make(chan ntwk.Message)
		ReqLoop(ntchan, out_req, out_rep)
		//log.Println("ping ", peer)
		//MakePingOld(req_chan, rep_chan)

	}
}

func run_node(config Configuration) {
	nlog.Println("run node ", config.NodePort)

	nlog.Println("PeerAddresses: ", config.PeerAddresses)

	//TODO signatures of genesis
	chain.InitAccounts()

	success := chain.ReadChain()

	nlog.Printf("block height %d", len(chain.Blocks))

	//nlog.Println("genesis block ", chain.Blocks[0])
	//chain.WriteGenBlock(chain.Blocks[0])

	//create new genesis block (demo)
	createDemo := !success
	if createDemo {
		genBlock := chain.MakeGenesisBlock()
		chain.ApplyBlock(genBlock)
		chain.AppendBlock(genBlock)
	}

	//if file exists read the chain

	// create block every 10sec

	delegation_enabled := false
	if delegation_enabled {
		go doEvery(blockTime, chain.MakeBlock)
	}

	//connect_peers(configuration.PeerAddresses)

	go ListenAll(config.NodePort)

}

func setupLogfile() *log.Logger {
	//setup log file

	logFile, err := os.OpenFile(logfile_name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		nlog.Fatal(err)
	}

	//defer logfile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	logger := log.New(logFile, "node ", log.LstdFlags)
	logger.SetOutput(mw)

	//log.SetOutput(file)

	nlog = logger
	return logger

}

func LoadConfiguration(file string) Configuration {
	var config Configuration
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

//HTTP
func LoadContent(peers []ntwk.Peer) string {
	content := ""

	content += fmt.Sprintf("<h2>Peers</h2>Peers: %d<br>", len(peers))
	for i := 0; i < len(peers); i++ {
		content += fmt.Sprintf("peer ip address: %s<br>", peers[i].Address)
	}

	content += fmt.Sprintf("<h2>TxPool</h2>%d<br>", len(chain.Tx_pool))

	for i := 0; i < len(chain.Tx_pool); i++ {
		//content += fmt.Sprintf("Nonce %d, Id %x<br>", chain.Tx_pool[i].Nonce, chain.Tx_pool[i].Id[:])
		ctx := chain.Tx_pool[i]
		content += fmt.Sprintf("%d from %s to %s %x<br>", ctx.Amount, ctx.Sender, ctx.Receiver, ctx.Id)
	}

	content += fmt.Sprintf("<h2>Accounts</h2>number of accounts: %d<br><br>", len(chain.Accounts))

	for k, v := range chain.Accounts {
		content += fmt.Sprintf("%s %d<br>", k, v)
	}

	content += fmt.Sprintf("<br><h2>Blocks</h2><i>number of blocks %d</i><br>", len(chain.Blocks))

	for i := 0; i < len(chain.Blocks); i++ {
		t := chain.Blocks[i].Timestamp
		tsf := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		//summary
		content += fmt.Sprintf("<br><h3>Block %d</h3>timestamp %s<br>hash %x<br>prevhash %x\n", chain.Blocks[i].Height, tsf, chain.Blocks[i].Hash, chain.Blocks[i].Prev_Block_Hash)

		content += fmt.Sprintf("<h4>Number of Tx %d</h4>", len(chain.Blocks[i].Txs))
		for j := 0; j < len(chain.Blocks[i].Txs); j++ {
			ctx := chain.Blocks[i].Txs[j]
			content += fmt.Sprintf("%d from %s to %s %x<br>", ctx.Amount, ctx.Sender, ctx.Receiver, ctx.Id)
		}
	}

	return content
}

func Runweb(webport int) {
	//webserver to access node state through browser
	// HTTP
	nlog.Println("start webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := LoadContent(Peers)
		//nlog.Print(p)
		fmt.Fprintf(w, "<h1>Polygon chain</h1><div>%s</div>", p)
	})

	nlog.Fatal(http.ListenAndServe(":"+strconv.Itoa(webport), nil))

}

func rungin(webport int) {
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	//router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// Setup route group for the API
	example := "test"
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, example)
		})

		api.GET("/peers", func(c *gin.Context) {
			c.JSON(http.StatusOK, Peers)
		})

		api.GET("/txpool", func(c *gin.Context) {
			c.JSON(http.StatusOK, chain.Tx_pool)
		})

		api.GET("/blockheight", func(c *gin.Context) {
			c.JSON(http.StatusOK, len(chain.Blocks))
		})

		//json: unsupported type: map[block.Account]int
		api.GET("/accounts", func(c *gin.Context) {
			c.JSON(http.StatusOK, chain.Accounts)
		})

		//blocks

	}

	// Start and run the server
	router.Run(":" + strconv.Itoa(webport))
}

func printsub(t time.Time) {
	log.Println("printsub")
	x := ntwk.TakeChan()

	nlog.Println("> ", x)
	//nlog.Println(len(ntwk.Timechan))

	// select {
	// case <-ntwk.Timechan:
	// 	fmt.Println("Received from Timechan. bufsize ", len(ntwk.Timechan))
	//}
}

func main() {

	//setup publisher
	//tchan = ntwk.Publisher()
	//go ntwk.Pubtime(tchan)

	setupLogfile()

	config := LoadConfiguration("nodeconf.json")

	nlog.Println("run node with config ", config)

	run_node(config)

	log.Println("run web on ", config.WebPort)

	//rungin(config.WebPort)

	Runweb(config.WebPort)

}
