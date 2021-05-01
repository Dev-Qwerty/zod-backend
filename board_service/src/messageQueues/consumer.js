const kafka = require('./kafka')
const userController = require('../controllers/userController')

const kafkaConsumer = async () => {
    try {
        const consumer = kafka.consumer({ groupId: 'group1' })

        await consumer.connect()
        await consumer.subscribe({ topic: 'project-js' })

        await consumer.run({
            eachMessage: async ({ topic, partition, message }) => {
                let value = JSON.parse(message.value)
                userController.newUser(value)
            }
        })
    } catch (error) {
        console.log(error)
    }
}


module.exports = kafkaConsumer()
