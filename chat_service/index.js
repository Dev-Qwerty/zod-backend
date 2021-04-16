const express = require('express')
require('dotenv').config()
require('./src/config/db')

const app = express()

app.listen("8080", () => {
    console.log("Server running in port 8080")
})