const meetController = require('../controllers/meetController')
const userController = require('../controllers/userController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {
    app.use('/api/meet', isUserVerified, meetController)
    app.use('/api/user', userController)
}

module.exports = router