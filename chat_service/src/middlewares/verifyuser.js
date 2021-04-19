const firebaseAuth = require('../config/firebase')

const VerifyUser = (req, res, next) => {
    const idToken = req.get('Authorization')
    firebaseAuth.verifyIdToken(idToken)
        .then((decodedToken) => {
            req.decodedToken = decodedToken
            next()
        })
        .catch((error) => {
            console.log('VerifyUser: ', error)
            res.status(401).send(error)
        })
}

module.exports = VerifyUser