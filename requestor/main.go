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
	"sync/atomic"
	"time"

	"github.com/nats-io/go-nats"
)

// Variables
var (
	debug                 = false
	defaultListenPort     = 6655
	defaultListenAddress  = "0.0.0.0"
	defaultRequestorCount = 1
	defaultSubject        = "demo.requests"
)

// Debugf prints to output if we are debugging
func Debugf(format string, v ...interface{}) {
	if debug {
		log.Printf(format, v...)
	}
}

// NOTE: Use tls scheme for TLS, e.g. nats-req -s tls://demo.nats.io:4443 foo hello
func usage() {
	log.Fatalf("Usage: requestor [-s server (%s)] [-nr requestor count] <subject> <msg> \n", nats.DefaultURL)
}

func main() {
	var (
		urls    string
		subject string
		nr      int
	)

	flag.StringVar(&urls, "s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&subject, "subj", defaultSubject, "The nats server URLs (separated by comma)")

	log.SetFlags(0)
	flag.Parse()

	args := flag.Args()
	if len(args) == 1 {
		subject = args[0]
	}

	log.Printf("Server URLs:     %s\n", urls)
	log.Printf("Requestor Count: %d\n", nr)
	log.Printf("Subject: %s\n", subject)

	nc, err := nats.Connect(urls)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	count := int64(0)

	// Start a group of requestors, simulating load.  Each requestor will
	// send a request, then store the duration into Prometheus.
	for true {
		c := atomic.AddInt64(&count, 1)

		// each request has a sequence for tracing
		payload := fmt.Sprintf("request-%d", c)

		// make the request and save the duration
		msg, err := nc.Request(subject, []byte(payload), 10*time.Second)
		if err != nil {
			Debugf("Request error: %v", err)
		}

		if msg != nil {
			log.Printf("%s", msg.Data)
		}
		time.Sleep(time.Millisecond * 25)
	}

	// Exit via the interrupt handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	nc.Close()

	fmt.Printf("Exiting...\n")
}
