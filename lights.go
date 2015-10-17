package main

import (
	"encoding/json"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"io/ioutil"
	"net/http"
	"strings"
)

// Some config
var (
	hubUrl = "http://192.168.1.5/api/1207d2d110346af72dd4490522f5712b"
	topic  = "/lights/set"
)

// Represents a Hue hub response
type hubResponse struct {
	State lightState `json:"state"`
}

// Represents a Hue light state
type lightState struct {
	On        bool    `json:"on"`
	Bri       int     `json:"bri,omitempty"`
	Hue       int     `json:"hue,omitempty"`
	Sat       int     `json:"sat,omitempty"`
	Colortemp float32 `json:"ct,omitempty"`
}

// Encodes a lightState to a json string
func (g *lightState) encode() string {
	output, err := json.Marshal(*g)
	if err != nil {
		panic(err)
	}
	return string(output)
}

// Returns a light's state
func lightGet(i int) lightState {
	url := fmt.Sprintf("%s/lights/%d", hubUrl, i)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var hubresp hubResponse
	err = json.Unmarshal(body, &hubresp)
	if err != nil {
		panic(err)
	}

	return hubresp.State
}

// Sends a lightState to the hub
func lightsSet(action lightState) {
	asJson := action.encode()

	url := fmt.Sprintf("%s/groups/0/action", hubUrl)
	req, err := http.NewRequest("PUT", url, strings.NewReader(asJson))
	if err != nil {
		panic(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
}

// Handles incoming MQTT messages
var onMessage MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
	switch msg.Payload()[0] {
	default:
		fmt.Printf("Unknown payload: %s\n", msg.Payload())
	case '0':
		lightsSet(lightState{On: false})
	case '1':
		lightsSet(lightState{On: true, Bri: 255, Colortemp: 370})
	case '2':
		lightsSet(lightState{On: true, Bri: 144, Sat: 211, Hue: 13122})
	case '3':
		lightsSet(lightState{On: true, Bri: 40, Sat: 211, Hue: 13122})
	case '4':
		state := lightGet(1)
		if !state.On {
			lightsSet(lightState{On: true, Bri: 255, Colortemp: 370})
		}
	}
}

func main() {
	// Create a ClientOptions struct setting the broker address, clientid, turn
	// off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("golights")
	opts.SetDefaultPublishHandler(onMessage)

	// Create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to topic
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Printf("Ready.\n")
	// Wait forever
	select {}
}
