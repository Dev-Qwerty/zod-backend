const express = require('express')
const { customAlphabet } = require('nanoid')

// Import models
const List = require('../models/list')

const router = express.Router()
const parseJson = express.json({ extended: true })

// Create 12 digit alphanumeric code for boardId and chardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 12)

// @route GET /api/list/:projectid
// @desc Fetch list of all boards
router
    .route('/new')
    .post([parseJson], async (req, res) => {

    })


module.exports = router
