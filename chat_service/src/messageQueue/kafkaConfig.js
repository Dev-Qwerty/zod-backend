const { Kafka } = require('kafkajs')

const kafka = new Kafka({
    clientId: 'chat-client',
    brokers: ['localhost:9092']
})

module.exports = kafka