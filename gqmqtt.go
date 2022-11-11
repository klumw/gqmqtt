/*
 * Copyright 2022 Winfried Klum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.bug.st/serial"
)

func main() {

	baudRates := []int{2400, 4800, 9600, 14400, 19200, 28800, 38400, 57600, 115200}
	models := []string{"GMC-500+", "GMC-320", "GMC-280"}
	cSize := 4
	foundPort := false

	var serPort = flag.String("s", "/dev/ttyUSB0", "serial port name")
	var verbose = flag.Bool("v", false, "verbose mode")
	var baudRate = flag.Uint("b", 115200, "baud rate")
	var model = flag.String("m", "GMC-500+", "GQ Geiger counter model")
	var topic = flag.String("t", "tele/geiger/cpm", "mqtt topic")
	var host = flag.String("h", "tcp://localhost:1883", "host url")
	var interval = flag.Uint("i", 60, "update interval for mqtt topic in seconds")
	var user = flag.String("u", "", "mqtt user")
	var pwd = flag.String("p", "", "mqtt password")

	flag.Parse()

	if !isElementInArray(int(*baudRate), &baudRates) {
		exitWithMsg("invalid baud rate")
	}

	if !isElementInArray(*model, &models) {
		exitWithMsg("invalid Geiger Counter model")
	}

	ports, err := serial.GetPortsList()
	if err != nil {
		exitWithError(err)
	}
	if len(ports) == 0 {
		exitWithMsg("no serial ports found!")
	}
	for _, port := range ports {
		if *serPort == port {
			foundPort = true
			if *verbose {
				fmt.Printf("found port: %v\n", port)
			}
		}
	}

	if !foundPort {
		exitWithMsg("serial port not found!")
	}

	mode := &serial.Mode{
		BaudRate: int(*baudRate),
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(*serPort, mode)
	if err != nil {
		exitWithError(err)
	}

	err = port.SetMode(mode)
	if err != nil {
		exitWithError(err)
	}

	client := createMqttClient(host, user, pwd)
	if *verbose {
		fmt.Print("mqtt client created\n")
	}
	if *model != "GMC-500+" {
		cSize = 2
		if *baudRate > 57600 {
			exitWithMsg("max. baud rate for selected model is 57600, but using " + strconv.FormatUint(uint64(*baudRate), 10))
		}
	}

	buff := make([]byte, cSize)

	for {
		_, err := port.Write([]byte("<GETCPM>>"))
		if err != nil {
			log.Fatal(err)
			break
		}
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n < (cSize - 1) {
			fmt.Println("\neof")
			break
		}
		cpm, _ := bytesToCpmValue(&buff)
		if *verbose {
			fmt.Printf("%d,", cpm)
		}
		sendCPM(client, cpm, topic)
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

func createMqttClient(host *string, user *string, pwd *string) MQTT.Client {

	opts := MQTT.NewClientOptions().AddBroker(*host)

	if *pwd != "" {
		opts.SetPassword(*pwd)
	}
	if *user != "" {
		opts.SetUsername(*user)
	}
	opts.SetClientID("Geiger Counter")

	c := MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		exitWithError(token.Error())
	}
	return c
}

func bytesToCpmValue(buff *[]byte) (uint32, error) {
	var err error
	switch len(*buff) {
	case 4:
		return binary.BigEndian.Uint32(*buff), err
	case 2:
		return uint32(binary.BigEndian.Uint16(*buff)), err
	default:
		err = errors.New("invalid byte length")
		return 0, err
	}
}

func sendCPM(client MQTT.Client, value uint32, topic *string) {
	if !client.IsConnected() {
		log.Fatal("mqtt not connected!")
	} else {
		token := client.Publish(*topic, byte(0), false, strconv.FormatUint(uint64(value), 10))
		token.Wait()
	}
}

func exitWithMsg(msg string) {
	log.Fatal(msg)
	os.Exit(1)
}

func exitWithError(err error) {
	log.Fatal(err)
	os.Exit(1)
}

func isElementInArray[T comparable](element T, elements *[]T) bool {
	for _, v := range *elements {
		if v == element {
			return true
		}
	}
	return false
}
