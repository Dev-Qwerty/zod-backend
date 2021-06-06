const meetController = require('../controllers/meetController')

const router = (app) => {
    app.use('/api/meet', meetController)
}

module.exports = router