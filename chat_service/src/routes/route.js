const channelController = require('../controllers/channelController')

const initializeRoutes = (app) => {
    app.use('/api/channel', channelController)
}

module.exports = initializeRoutes