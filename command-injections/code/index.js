'use strict';

const bodyParser = require('body-parser');
const _ = require('lodash');

const express = require('express');
const app = express();
app.use(bodyParser.json());

app.post('/echo', (req, res) => {
    const out = {
        value: 'echo',
        formatter: () => {
            return {'echo': out.value};
        }
    };
    
    _.merge(out, req.body);
    const outStr = out.formatter();
    res.json(outStr);
});

app.listen(8000);