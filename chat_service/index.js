const express = require('express')
require('dotenv').config()

const VerifyUser = require('./src/middlewares/verifyuser')
require('./src/config/db')
require('./src/config/firebase')

const initializeRoutes = require('./src/routes/route')

const app = express()

app.use(VerifyUser)

initializeRoutes(app)

app.listen(process.env.PORT, () => {
    console.log(`Server running on port ${process.env.PORT}`)
})