const meetController = require('../controllers/meetController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {
    app.use('/api/meet', isUserVerified, meetController)
}

module.exports = router