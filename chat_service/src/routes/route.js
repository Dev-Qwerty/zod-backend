const channelController = require('../controllers/channelController')

const initializeRoutes = (app) => {
    app.use('/channel', channelController)
}

module.exports = initializeRoutes