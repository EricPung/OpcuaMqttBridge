package models

import "gorm.io/gorm"

// OPCUAServer holds connection info for an OPC UA server.
type OPCUAServer struct {
	gorm.Model
	URL            string
	SecurityMode   string
	SecurityPolicy string
	AuthType       string
	Username       string
	Password       string
	Points         []Point
}

// MQTTBroker holds connection info for an MQTT broker.
type MQTTBroker struct {
	gorm.Model
	URL      string
	Username string
	Password string
	Points   []Point
}

// Point links an OPC UA node to an MQTT topic.
type Point struct {
	gorm.Model
	OPCUAServerID uint
	MQTTBrokerID  uint
	NodeID        string
	Topic         string
}
