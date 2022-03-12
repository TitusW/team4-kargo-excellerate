const express = require('express');
const compression = require('compression');
const helmet = require('helmet');
const router = require('./api');

const notFoundMiddleware = require('../middlewares/notFoundHandler');


const app = express();
app.use(helmet());
app.use(compression());
app.use(express.json());
app.use('/v1', router);
app.use(notFoundMiddleware);

module.exports = app;