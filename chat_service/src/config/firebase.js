const admin = require("firebase-admin");

try {
    admin.initializeApp({
        credential: admin.credential.cert(process.env.FIREBASE_SERVICE_ACCOUNT)
    })
} catch (error) {
    console.log("Initiaize firebase: ", error)
}

const firebaseAuth = admin.auth()

module.exports = firebaseAuth;