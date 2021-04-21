const mongoose = require('mongoose')

mongoose.connect(process.env.DBURI, { useNewUrlParser: true, useUnifiedTopology: true, useFindAndModify: false })
    .catch((error) => {
        console.log(error);
    });

db = mongoose.connection

db.on('open', () => {
    console.log('Connected to db');
})

db.on('error', () => {
    console.log('Cant connect to db');
})