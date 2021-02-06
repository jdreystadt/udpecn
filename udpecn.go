/*
Copyright 2021 John Dreystadt

Licensed under the MIT License.

This software is so simple you should probably just read it rather
than try to use any of the actual code. But feel free to do whatever.
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var modeStr = "both"

type valueStruct struct {
	modeValue *string
}

func (mode valueStruct) String() string {
	if mode.modeValue == nil {
		return ""
	}
	return *mode.modeValue
}

func (mode valueStruct) Set(s string) error {
	switch s {
	case "both", "in", "out":
		*mode.modeValue = s
		return nil
	default:
		*mode.modeValue = ""
		return errors.New("mode must be both, in, or out")
	}
}
func main() {
	// Global variables
	var fport = flag.String("f", "8001", "Outgoing port")                // First port
	var sport = flag.String("s", "8002", "Incoming port")                // Second port
	var dest = flag.String("d", "127.0.0.1", "Destination IPv4 address") // Default destination
	// mode.modeValue = "both"
	value := valueStruct{&modeStr}
	flag.Var(&value, "m", "both, in or out")
	mode := value.modeValue

	// Channel for catching control-c
	cCatch := make(chan os.Signal, 1)
	signal.Notify(cCatch, os.Interrupt)

	// Channels for 1 second timer and 5 second timer
	c1Ticker := time.NewTicker(1 * time.Second)
	c1Timer := c1Ticker.C
	c5Ticker := time.NewTicker(5 * time.Second)
	c5Timer := c5Ticker.C

	// Parse command line for all entries
	flag.Parse()
	// Open network ports
	outFirstPort, inFirstPort, error := openPorts(*fport, *dest)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	outSecondPort, inSecondPort, error := openPorts(*sport, *dest)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	fmt.Println(*outFirstPort)
	fmt.Println(*inFirstPort)
	fmt.Println(*outSecondPort)
	fmt.Println(*inSecondPort)
	// Outgoing first port

	// Main Loop
MainLoop:
	for {
		select {
		case <-c1Timer:
			fmt.Println("1 Second Timer")
		case <-c5Timer:
			fmt.Println("5 Second Timer")
		case <-cCatch:
			fmt.Println("Got Control-C")
			c1Ticker.Stop()
			c5Ticker.Stop()
			break MainLoop
		}
	}
	// Shutdown
	fmt.Println(*fport)
	fmt.Println(*sport)
	fmt.Println(*dest)
	fmt.Println(*mode)
	fmt.Println("Goodby")
}

func openPorts(port string, addr string) (*net.UDPConn, *net.UDPConn, error) {
	lstr := net.JoinHostPort("", port)
	rstr := net.JoinHostPort(addr, port)
	laddr, error := net.ResolveUDPAddr("udp", lstr)
	if error != nil {
		fmt.Println(error)
		return nil, nil, error
	}
	raddr, error := net.ResolveUDPAddr("udp", rstr)
	if error != nil {
		fmt.Println(error)
		return nil, nil, error
	}
	outPort, error := net.DialUDP("udp", laddr, raddr)
	if error != nil {
		fmt.Println(error)
		return nil, nil, error
	}
	inPort, error := net.ListenUDP("udp", laddr)
	if error != nil {
		fmt.Println(error)
		return nil, nil, error
	}
	return outPort, inPort, nil
}
