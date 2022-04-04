package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	TA "github.com/GGroups/rttm_svtags/tagmd"

	COMM "github.com/GGroups/rttm_login/comm"
	log "github.com/cihub/seelog"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	log.Info("#->Login service started!")

	TA.InitTag()

	taobj := TA.TagM{}
	ep1 := TA.MakeEndPointAddOneTag(taobj)
	ep2 := TA.MakeEndPointGetAllTag(taobj)
	ep3 := TA.MakeEndPointSetOneTag(taobj)

	svr1 := httpTransport.NewServer(ep1, TA.DecodeRequestAddOneTag, COMM.CommEncodeResponse)
	svr2 := httpTransport.NewServer(ep2, TA.DecodeRequestEmptyReq, COMM.CommEncodeResponse)
	svr3 := httpTransport.NewServer(ep3, TA.DecodeRequestAddOneTag, COMM.CommEncodeResponse)

	routeSvr := mux.NewRouter()

	routeSvr.Handle(`/rttm/tags/AddOne`, svr1).Methods("POST")
	routeSvr.Handle(`/rttm/tags/GetAll`, svr2).Methods("POST")
	routeSvr.Handle(`/rttm/tags/SetOne`, svr3).Methods("POST")

	//min loop
	ch := make(chan error, 2)
	go func() {
		log.Info("0.0.0.0:18001", `/rttm/tags**`)
		ch <- http.ListenAndServe("0.0.0.0:18001", routeSvr)
	}()
	go func() {
		log.Info("##", "wait for exit sigint...")
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		ch <- fmt.Errorf("%s", <-c)
	}()

	log.Info("MainSvr Terminated", <-ch)
}
