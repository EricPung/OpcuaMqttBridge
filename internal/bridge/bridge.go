package bridge

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
	opcua "github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"gorm.io/gorm"

	"github.com/example/opcuamqttbridge/internal/models"
	mqttclient "github.com/example/opcuamqttbridge/internal/mqtt"
	opcuaclient "github.com/example/opcuamqttbridge/internal/opcua"
)

// Start loads configuration from the database and starts bridging.
func Start(ctx context.Context, db *gorm.DB) error {
	var brokers []models.MQTTBroker
	if err := db.Find(&brokers).Error; err != nil {
		return err
	}
	brokerClients := make(map[uint]mqttlib.Client)
	for _, b := range brokers {
		c, err := mqttclient.Connect(b.URL, b.Username, b.Password)
		if err != nil {
			log.Printf("mqtt connect failed for %s: %v", b.URL, err)
			continue
		}
		brokerClients[b.ID] = c
	}

	var servers []models.OPCUAServer
	if err := db.Find(&servers).Error; err != nil {
		return err
	}

	var points []models.Point
	if err := db.Find(&points).Error; err != nil {
		return err
	}
	pointsByServer := make(map[uint][]models.Point)
	for _, p := range points {
		pointsByServer[p.OPCUAServerID] = append(pointsByServer[p.OPCUAServerID], p)
	}

	var wg sync.WaitGroup
	for _, s := range servers {
		pts := pointsByServer[s.ID]
		if len(pts) == 0 {
			continue
		}
		wg.Add(1)
		go func(srv models.OPCUAServer, pts []models.Point) {
			defer wg.Done()
			runServer(ctx, srv, pts, brokerClients)
		}(s, pts)
	}
	wg.Wait()
	return nil
}

func runServer(ctx context.Context, srv models.OPCUAServer, points []models.Point, brokers map[uint]mqttlib.Client) {
	client, err := opcuaclient.Connect(ctx, srv.URL, srv.SecurityMode, srv.SecurityPolicy, srv.Username, srv.Password)
	if err != nil {
		log.Printf("opcua connect failed for %s: %v", srv.URL, err)
		return
	}
	defer client.Close()

	sub, err := client.Subscribe(&opcua.SubscriptionParameters{Interval: time.Second}, nil)
	if err != nil {
		log.Printf("subscription failed: %v", err)
		return
	}
	defer sub.Cancel()

	handleToPoint := make(map[uint32]models.Point)
	for i, p := range points {
		id, err := ua.ParseNodeID(p.NodeID)
		if err != nil {
			log.Printf("invalid node id %s: %v", p.NodeID, err)
			continue
		}
		h := uint32(i + 1)
		req := &ua.MonitoredItemCreateRequest{
			ItemToMonitor:       &ua.ReadValueID{NodeID: id, AttributeID: ua.AttributeIDValue},
			MonitoringMode:      ua.MonitoringModeReporting,
			RequestedParameters: &ua.MonitoringParameters{ClientHandle: h, SamplingInterval: 1000},
		}
		_, err = sub.Monitor(ua.TimestampsToReturnBoth, req)
		if err != nil {
			log.Printf("monitor failed: %v", err)
			continue
		}
		handleToPoint[h] = p
	}

	ch := sub.Run(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-ch:
			if res == nil {
				continue
			}
			switch x := res.(type) {
			case *ua.DataChangeNotification:
				for _, item := range x.MonitoredItems {
					p := handleToPoint[item.ClientHandle]
					m := brokers[p.MQTTBrokerID]
					if m == nil {
						continue
					}
					b, _ := json.Marshal(item.Value.Value())
					m.Publish(p.Topic, 0, false, b)
				}
			}
		}
	}
}
