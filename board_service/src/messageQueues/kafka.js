const { Kafka } = require('kafkajs')

const kafka = new Kafka({
    clientId: 'board-service',
    brokers: ['localhost:9092']
})

module.exports = kafka
