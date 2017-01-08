
../context/build/SignalServer : SignalServer.go Dockerfile
	docker build -t golang_signalserver .
	~/Desktop/project/dkrcp/dkrcp.sh golang_signalserver::/go/bin/SignalServer ../context/build/SignalServer
	docker rmi golang_signalserver
