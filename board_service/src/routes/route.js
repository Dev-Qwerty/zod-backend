const boardController = require('../controllers/boardController')
const cardController = require('../controllers/cardController')
const listController = require('../controllers/listController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {

    app.use('/api/board', isUserVerified, boardController)
    // app.use('/api/card', isUserVerified, cardController)
    // app.use('/api/list', isUserVerified, listController)
    
}

module.exports = router
