package main

import (
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/heatxsink/go-hue/src/groups"
	"github.com/heatxsink/go-hue/src/lights"
)

// Some config
var (
	topic       = "/lights/set"
	hueIP       = "192.168.1.5"
	hueUsername = "1207d2d110346af72dd4490522f5712b"

	l = lights.NewLights(hueIP, hueUsername)
	g = groups.NewGroup(hueIP, hueUsername)
)

// Handles incoming MQTT messages
var onMessage MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
	switch msg.Payload()[0] {
	default:
		fmt.Printf("Unknown payload: %s\n", msg.Payload())
	case '0':
		g.SetGroupState(0, lights.State{On: false})
	case '1':
		g.SetGroupState(0, lights.State{On: true, Bri: 255, Ct: 370})
	case '2':
		g.SetGroupState(0, lights.State{On: true, Bri: 144, Sat: 211, Hue: 13122})
	case '3':
		g.SetGroupState(0, lights.State{On: true, Bri: 40, Sat: 211, Hue: 13122})
	case '4':
		state := l.GetLight(1).State
		if !state.On {
			g.SetGroupState(0, lights.State{On: true, Bri: 255, Ct: 370})
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
