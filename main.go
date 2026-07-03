package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
)

var version string

const UNKNOWN = 3
const CRITICAL = 2
const WARNING = 1
const OK = 0

func printVersion() {
	fmt.Printf(`%s Compiler: %s %s`,
		version,
		runtime.Compiler,
		runtime.Version())
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opt := Opt{}
	psr := flags.NewParser(&opt, flags.Default)
	_, err := psr.Parse()
	if err != nil {
		os.Exit(UNKNOWN)
	}

	if opt.Version {
		printVersion()
		return OK
	}

	opt.bufferSize = uint64(opt.MaxBufferSize)

	if opt.WaitFor && opt.WaitForMax == 0 {
		fmt.Printf("wait-for-max is required when wait-for is enabled\n")
		return UNKNOWN
	}

	if opt.ExpectContent != "" && opt.Base64ExpectContent != "" {
		fmt.Printf("Both string and base64-string are specified\n")
		return UNKNOWN
	}

	if opt.ExpectContent != "" {
		opt.expectByte = []byte(opt.ExpectContent)
	}
	if opt.Base64ExpectContent != "" {
		data, err := base64.StdEncoding.DecodeString(opt.Base64ExpectContent)
		if err != nil {
			fmt.Printf("Failed decode base64-string: %v\n", err)
			return UNKNOWN
		}
		opt.expectByte = data
	}

	if opt.TCP4 && opt.TCP6 {
		fmt.Printf("Both tcp4 and tcp6 are specified\n")
		return UNKNOWN
	}

	if opt.SNI && opt.Hostname == "" {
		fmt.Printf("hostname is required when use sni\n")
		return UNKNOWN
	}

	if opt.Hostname == "" && opt.IPAddress == "" {
		fmt.Printf("Specify either hostname or ipaddress\n")
		return UNKNOWN
	}

	if opt.Hostname == "" {
		opt.Hostname = opt.IPAddress
	}

	if opt.IPAddress == "" {
		host, _, err := net.SplitHostPort(opt.Hostname)
		if err != nil {
			opt.IPAddress = opt.Hostname
		} else {
			opt.IPAddress = host
		}
	}

	if opt.Port == 0 {
		_, port, err := net.SplitHostPort(opt.Hostname)
		if err == nil {
			p, _ := strconv.Atoi(port)
			// skip error check OK
			opt.Port = p
		}
	}

	if opt.Port == 0 {
		if opt.SSL {
			opt.Port = 443
		} else {
			opt.Port = 80
		}
	}

	if opt.URI == "" {
		opt.URI = "/"
	}

	transport := opt.MakeTransport()
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: opt.Timeout,
	}

	ctx := context.Background()
	timeout := opt.Timeout + 3*time.Second
	if opt.WaitForMax > 0 {
		timeout = opt.WaitForMax
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	requestNum := 0
	if opt.WaitFor {
		consecutive := opt.Consecutive - 1
		for ctx.Err() == nil {
			requestNum++
			okMsg, reqErr := opt.Request(ctx, client)
			interval := opt.Interim
			if reqErr == nil && consecutive <= 0 {
				log.Printf("request[%d]: %s", requestNum, okMsg)
				fmt.Println(okMsg)
				return OK
			} else if reqErr == nil {
				consecutive--
				log.Printf("request[%d]: %s", requestNum, okMsg)
			} else {
				interval = opt.WaitForInterval
				consecutive = opt.Consecutive - 1
				log.Printf("request[%d]: %s", requestNum, reqErr.Error())
			}
			select {
			case <-ctx.Done():
			case <-time.After(interval):
			}
		}
		fmt.Printf("Give up waiting for success\n")
		return UNKNOWN
	}

	consecutive := opt.Consecutive - 1
	var rErr *RequestError
	for ctx.Err() == nil {
		var okMsg string
		requestNum++
		okMsg, rErr = opt.Request(ctx, client)
		if rErr == nil && consecutive <= 0 {
			log.Printf("request[%d]: %s", requestNum, okMsg)
			fmt.Println(okMsg)
			return OK
		} else if rErr == nil {
			consecutive--
			log.Printf("request[%d]: %s", requestNum, okMsg)
		} else {
			break
		}
		select {
		case <-ctx.Done():
		case <-time.After(opt.Interim):
		}
	}
	fmt.Println(rErr.Error())
	return rErr.Code()
}
