package main

import (
	"fmt"
	"net"
)

/* A Simple function to verify error */
func PanicCheck(err error, moreContext string) {
	if err == nil {
		return
	}
	panicmess := "Panic: " + err.Error() + " " + moreContext
	fmt.Println("Error: " + err.Error() + " " + moreContext)
	panic(panicmess)
}

func MessageInfom(message string) {
	fmt.Println("Inform: " + message)
}

type SignalConn net.UDPConn

func SignalServer(port string) *SignalConn {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp4", ":"+port)
	PanicCheck(err, "resolving UDP addr. Port: "+port)
	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	PanicCheck(err, "establishing connection")

	return (*SignalConn)(ServerConn)
}

var clientIP net.UDPAddr

func (conn *SignalConn) SignalAccept(signal chan string) {

	buf := make([]byte, 1024)
	srvr := (*net.UDPConn)(conn)
	for {
		size, clientIPx, err := srvr.ReadFromUDP(buf)
		clientIP = *clientIPx
		PanicCheck(err, "size: "+string(size)+" IP add: "+string(clientIP.IP))
		signal <- string(buf[0:size])
	}
}

func main() {
	// create signal server
	sc := signalServerConfigure("10001")
	defer sc.Close()
	emttr, err := configEventEmttr(sc)
	PanicCheck(err, "configuring event emitter")
	st, err := configStatMchn()
	PanicCheck(err, "configuring state machine")
	MessageInfom("Running state machine")
	err = st.CommandRun(emttr)
	PanicCheck(err, "running initial state machine")
}

//-----------------------------------------------------------------------------
type eventLabel string

type eventPacket interface {
	toString() string
}

//-----------------------------------------------------------------------------
type eventPacketEmitterTest struct {
	emitFn func() string
}

func (e eventPacketEmitterTest) emit() (eventLabel, eventPacket, error) {
	msg := e.emitFn()
	pkt := new(eventPacketString)
	pkt.init(msg)
	return (eventLabel)(msg[0:4]), pkt, nil
}

//-----------------------------------------------------------------------------
func signalServerConfigure(port string) *SignalConn {
	MessageInfom("Starting Server")
	conn := SignalServer("10001")
	stStop.stStopDef(conn, &clientIP)
	return conn
}
func configEventEmttr(sc *SignalConn) (eventPacketEmitter, error) {
	em := new(eventPacketEmitterTest)
	signal := make(chan string)
	MessageInfom("Accepting Signals")
	go sc.SignalAccept(signal)
	em.emitFn = func() string {
		return <-signal
	}
	return em, nil
}

func configStatMchn() (stateMchn, error) {
	st := new(stateMchnTest)
	st.init()
	return st, nil
}

//-----------------------------------------------------------------------------
type eventPacketString struct {
	packet string
}

func (p *eventPacketString) init(packet string) {
	p.packet = packet
}
func (p *eventPacketString) toString() string {
	return p.packet
}

//-----------------------------------------------------------------------------
type eventPacketEmitter interface {
	emit() (eventLabel, eventPacket, error)
}

//-----------------------------------------------------------------------------
type stateMchnCmmd func(eventPacket) (stateMchn, error)

//-----------------------------------------------------------------------------
type stateMchn interface {
	CommandSelect(eventLabel) (stateMchnCmmd, error)
	CommandRun(eventPacketEmitter) error
}

//-----------------------------------------------------------------------------

type stateMchnTest struct {
	stMp map[eventLabel]stateMchnCmmd
}

func (st *stateMchnTest) CommandSelect(evLbl eventLabel) (stateMchnCmmd, error) {
	return st.stMp[evLbl], nil
}

func (st *stateMchnTest) init() error {
	st.stMp = make(map[eventLabel]stateMchnCmmd)
	var pumpCnt int = 0

	st.stMp["Pump"] = func(evP eventPacket) (stateMchn, error) {
		MessageInfom(evP.toString())
		pumpCnt += 1
		if pumpCnt > 23 {
			return &stStop, nil
		}
		return st, nil
	}
	return nil
}

func (stm *stateMchnTest) CommandRun(evPktEmttr eventPacketEmitter) error {
	for {
		evLbl, evPkt, err := evPktEmttr.emit()
		PanicCheck(err, "event emit")
		stmCmmd, err := stm.CommandSelect(evLbl)
		PanicCheck(err, "selected command")
		stmNext, err := stmCmmd(evPkt)
		PanicCheck(err, "ran command")
		if stm != stmNext {
			err := stmNext.CommandRun(evPktEmttr)
			PanicCheck(err, "running next state machine")
		}
	}
	return nil
}

var stStop stateMchnTest

func (st *stateMchnTest) stStopDef(sc *SignalConn, cIP *net.UDPAddr) {
	st.stMp = make(map[eventLabel]stateMchnCmmd)
	st.stMp["Pump"] = func(evp eventPacket) (stateMchn, error) {
		_, err := ((*net.UDPConn)(sc)).WriteToUDP([]byte("Stop"), cIP)
		MessageInfom("Sending stop message to client." + fmt.Sprintf(" %v", clientIP))
		return st, err
	}
}
