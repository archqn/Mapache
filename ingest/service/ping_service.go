package service

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"ingest/database"
	"ingest/rabbitmq"
	"ingest/utils"
	"time"
)

func Ping() map[string]interface{} {
	queryStatus, queryLatency := PingQuery()
	singlestoreStatus, singlestoreLatency := PingSingleStore()
	rabbitStatus, rabbitLatency := PingRabbitMQ()
	return map[string]interface{}{
		"query": map[string]interface{}{
			"status":  BoolToStatusString(queryStatus),
			"latency": queryLatency,
		},
		"singlestore": map[string]interface{}{
			"status":  BoolToStatusString(singlestoreStatus),
			"latency": singlestoreLatency,
		},
		"rabbit": map[string]interface{}{
			"status":  BoolToStatusString(rabbitStatus),
			"latency": rabbitLatency,
		},
	}
}

func PingQuery() (bool, int) {
	return false, 0
}

func PingSingleStore() (bool, int) {
	start := time.Now()
	err := database.PingDB()
	if err != nil {
		utils.SugarLogger.Errorln("Failed to ping singlestore: ", err)
		return false, int(time.Since(start).Milliseconds())
	}
	utils.SugarLogger.Infoln("Pinged singlestore in ", time.Since(start).Milliseconds(), "ms")
	return true, int(time.Since(start).Milliseconds())
}

func PingRabbitMQ() (bool, int) {
	start := time.Now()
	ping, _ := uuid.NewUUID()
	topic := "meta/ping"
	token := rabbitmq.Client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		if string(msg.Payload()) == ping.String() {
			go rabbitmq.Client.Unsubscribe(topic)
		}
	})
	token.Wait()
	token = rabbitmq.Client.Publish(topic, 1, false, ping.String())
	sent := token.WaitTimeout(5 * time.Second)
	if sent {
		utils.SugarLogger.Infoln("Pinged rabbitmq in ", time.Since(start).Milliseconds(), "ms")
	} else {
		utils.SugarLogger.Errorln("Failed to ping rabbitmq: ", token.Error())
	}
	return sent, int(time.Since(start).Milliseconds())
}

func BoolToStatusString(b bool) string {
	if b {
		return "ONLINE"
	}
	return "OFFLINE"
}
