const express = require('express')
require('dotenv').config()

require('./src/config/db')

const initializeRoutes = require('./src/routes/route')

const app = express()

initializeRoutes(app)

app.listen(process.env.PORT, () => {
    console.log(`Server running on port ${process.env.PORT}`)
})