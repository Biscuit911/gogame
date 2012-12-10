package main

import (
	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/Nightgunner5/netchan"
	"net"
)

const (
	MaxQueue = 8
)

var (
	broadcast chan<- packet.Packet
	connected chan<- chan<- packet.Packet
)

func init() {
	broadcast_ := make(chan packet.Packet)
	broadcast = broadcast_

	connected_ := make(chan chan<- packet.Packet)
	connected = connected_

	go func() {
		connections := make(map[chan<- packet.Packet]bool)
		for {
			select {
			case p := <-broadcast_:
				for c := range connections {
					select {
					case c <- p:
					default:
					}
				}
			case c := <-connected_:
				connections[c] = true
			}
		}
	}()
}

func main() {
	ln, err := net.Listen("tcp", ":7031")
	if err != nil {
		panic(err)
	}

	netchan.Listen(ln, func(addr net.Addr, c *netchan.Chan) {
		recv := c.ChanRecv().(<-chan packet.Packet)
		send := c.ChanSend().(chan<- packet.Packet)

		NewPlayer(addr, recv, send)

		// TODO: disconnect logic
		world.onConnect <- send
		connected <- send
	}, packet.Type, MaxQueue)
}