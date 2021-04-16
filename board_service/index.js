const express = require('express')
require('dotenv').config()
require('./src/config/db')

const app = express()

app.listen(3000, () => {
    console.log('server running...')
});
