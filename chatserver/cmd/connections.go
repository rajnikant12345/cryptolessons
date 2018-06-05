package cmd

import (
	"net"
	"cryptolessons/chatserver/config"
	"cryptolessons/chatserver/model"
	"encoding/json"
	"io"
	"log"
	"cryptolessons/chatserver/processors"
)

var ch chan channelData

type channelData struct {
	m model.CommonMessage
	encoder net.Conn
}

func processMessage(ch chan channelData) {
	for m := range ch {
		processors.ProcessMessage(m.m , m.encoder)
	}
}


func StartWorkers() {
	chansz := config.CFG.ChannelSize

	if chansz == 0 {
		chansz = 1000
	}

	ch = make(chan channelData , chansz)

	for i:=0;i<config.CFG.Workers;i++ {
		go processMessage(ch)
	}
}

//TODO: Handle timeout of an idle connection
//TODO: remove termonated connection from usermap

func HandleConnections(c net.Conn, conn string) {

	for {
		m := model.CommonMessage{}
		j := json.NewDecoder(c)
		e := j.Decode(&m)
		if e != nil {
			if e == io.EOF {
				log.Println("Conection terminated.")
				return
			}else {
				log.Println(e.Error())
				continue
			}
		}
		//validate conection id
		if m.Conn != conn {
			c.Close()
			break
		}

		cdata := channelData{m,c}

		ch <- cdata
	}
}