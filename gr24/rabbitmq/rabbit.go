package rabbitmq

import (
	"fmt"
	"gr24/config"
	"gr24/service"
	"gr24/utils"
	"math/rand"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Client mqtt.Client

var clientID string

func InitializeRabbit() {
	id := rand.Intn(100)
	clientID = "ingest_mqtt_" + strconv.Itoa(id)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", config.MQTTHost, config.MQTTPort))
	opts.SetClientID(clientID)
	opts.SetUsername(config.MQTTUser)
	opts.SetPassword(config.MQTTPassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		utils.SugarLogger.Fatalln(token.Error())
	}
	Client = client
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	utils.SugarLogger.Infoln("[MQ] Received message: " + string(msg.Payload()) + " from topic: " + msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	utils.SugarLogger.Infoln("[MQ] Connected to RabbitMQ as: " + clientID)
	InitializeIngest()
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	utils.SugarLogger.Infoln("[MQ] Connection lost: ", err)
}

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 0, nil)
	token.Wait()
	utils.SugarLogger.Infoln("[MQ] Subscribed to topic: ", topic)
}

func InitializeIngest() {
	sub(Client, "meta")
	subscribePedal(Client)
	subscribeACU(Client)
	subscribeBCM(Client)
	subscribeWheel(Client)
	subscribeSteeringWheel(Client)
	subscribeVDM(Client)
}

func subscribePedal(client mqtt.Client) {
	client.Subscribe("gr24/+/pedal", 0, service.PedalIngestCallback)
	utils.SugarLogger.Infoln("[MQ] Subscribed to topic: gr24/+/pedal")
}

func subscribeACU(client mqtt.Client) {
	client.Subscribe("gr24/+/acu", 0, service.ACUIngestCallback)
	utils.SugarLogger.Infoln("[MQ] Subscribed to topic: gr24/+/acu")
}

func subscribeBCM(client mqtt.Client) {
	client.Subscribe("gr24/+/bcm", 0, service.BCMIngestCallback)
	utils.SugarLogger.Infoln("[MQ] Subscribed to topic: gr24/+/bcm")
}

func subscribeWheel(client mqtt.Client) {
	client.Subscribe("gr24/+/wheel/+", 0, service.WheelIngestCallback)
	utils.SugarLogger.Infoln("[MQ] Subscribed to topic: gr24/+/wheel/+")
}

func subscribeSteeringWheel(client mqtt.Client) {
	client.Subscribe("gr24/+/steering_wheel", 0, service.SteeringWheelIngestCallback)
	utils.SugarLogger.Infoln("[MQ] Subscribed to topic: gr24/+/steering_wheel")
}

func subscribeVDM(client mqtt.Client) {
	client.Subscribe("gr24/+/vdm", 0, service.VDMIngestCallback)
	utils.SugarLogger.Infoln("[MQ] Subscribed to topic: gr24/+/vdm")
}
