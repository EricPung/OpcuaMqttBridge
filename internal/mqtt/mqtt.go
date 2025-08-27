package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Connect creates an MQTT client and connects to broker.
func Connect(url, username, password string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker(url)
	if username != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, err
	}
	return client, nil
}
