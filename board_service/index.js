const express = require('express')
require('dotenv').config()

require('./src/config/db')
const router = require('./src/routes/route')
// require('./src/messageQueues/consumer').kafkaConsumer

const app = express()

router(app)

app.listen(3000, () => {
    console.log('server running...')
});
