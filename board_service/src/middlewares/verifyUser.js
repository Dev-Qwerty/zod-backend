const auth = require('../config/firebase')

const isUserVerified = (req, res, next) => {
    auth.verifyIdToken(req.headers.authorization).then((decodedToken) => {
        req.decodedToken = decodedToken
        next()
    }).catch( (error) => {
        res.status(401).send(error)
    })
}

module.exports = isUserVerified
