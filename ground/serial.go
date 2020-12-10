package main

import (
	"encoding/binary"
	"io"
	"strconv"
)

func listenPort(stream io.ReadWriteCloser, bufferSize int, result chan *hamsterTongueMessage) {
	buffer := make([]byte, bufferSize)
	var (
		markerFound bool
		appendCount byte
		message     *hamsterTongueMessage = &hamsterTongueMessage{Payload: []byte{}}
	)
	for {
		readCount, err := stream.Read(buffer)
		if readCount > 0 {
			for _, b := range buffer[0:readCount] {
				if markerFound {
					appendCount++
					switch appendCount {
					case 1:
						message.Length = b
					case 2:
						message.Verb = b
					case 3:
						message.Noun = b
					default:
						if appendCount <= message.Length {
							message.Payload = append(message.Payload, b)
						} else {
							message.CRC = b
							result <- message
							message = &hamsterTongueMessage{Payload: []byte{}}
							markerFound = false
							appendCount = 0
						}
					}
				} else {
					if b == 0xFF {
						markerFound = true
					}
				}
			}
		}
		if err != nil {
			globalLogger.WithField("error", err).Errorln("encounter error while read serial port")
			return
		}
	}
}

func decodeMessage(msgchan chan *hamsterTongueMessage) {
	for {
		select {
		case msg := <-msgchan:
			globalLogger.Debugf("message income : %#v", msg)
			switch msg.Verb {
			case 0: //Heartbeat
			case 1: //Value
				Value[strconv.Itoa(int(msg.Noun))] = binary.LittleEndian.Uint32(msg.Payload)
			}
		}
	}
}