// Copyright 2012-2018 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/go-nats"
)

// Variables
var (
	debug                = false
	defaultListenAddress = "0.0.0.0"
	defaultDelay         = "50ms"
	defaultSubject       = "demo.requests"
	defaultQueueGroup    = "demo"
)

// Debugf prints to output if we are debugging
func Debugf(format string, v ...interface{}) {
	if debug {
		log.Printf(format, v...)
	}
}

func main() {
	var (
		urls     string
		subject  string
		qgroup   string
		delayStr string
	)

	flag.StringVar(&urls, "s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&subject, "subj", defaultSubject, "The subject to listen to")
	flag.StringVar(&qgroup, "qg", defaultQueueGroup, "The name of the queue group")
	flag.StringVar(&delayStr, "delay", defaultDelay, "Duration to delay the response")
	flag.BoolVar(&debug, "debug", false, "Enable debugging")

	log.SetFlags(0)
	flag.Parse()

	nc, err := nats.Connect(urls)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}
	defer nc.Drain()

	delay, err := time.ParseDuration(delayStr)
	if err != nil {
		log.Fatalf("Couldn't parse delay: %v", err)
	}

	_, err = nc.QueueSubscribe(subject, qgroup, func(msg *nats.Msg) {
		time.Sleep(delay)
		fmt.Printf("Received: %s\n", msg.Data)
		nc.Publish(msg.Reply, msg.Data)
	})
	if err != nil {
		log.Fatalf("couldn't subscribe: %v", err)
	}

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Printf("Draining...\n")
	nc.Drain()
	fmt.Printf("Exiting.\n")
}
