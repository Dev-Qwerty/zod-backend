const boardController = require('../controllers/boardController')
const cardController = require('../controllers/cardController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {

    app.use('/api/board', isUserVerified, boardController)
    app.use('/api/card', isUserVerified, cardController)
    
}

module.exports = router
