const kafka = require('./kafkaConfig')
const { CreateChannelEveryone,
    AcceptProjectInvite
} = require('./service')

const kafkaConsumer = async () => {
    try {
        const consumer = kafka.consumer({ groupId: 'chat-group' })
        await consumer.connect()
        await consumer.subscribe({ topic: 'Project', fromBeginning: true })
        await consumer.run({
            eachMessage: async ({ topic, partition, message }) => {
                if (message.key.toString() == 'Create Project') {
                    CreateChannelEveryone(JSON.parse(message.value))
                } else if (message.key.toString() == 'Accept Project Invite') {
                    AcceptProjectInvite(JSON.parse(message.value))
                }
            }
        })
    } catch (error) {
        console.log(`kafkaConsumer: ${error}`)
    }
}

module.exports = kafkaConsumer