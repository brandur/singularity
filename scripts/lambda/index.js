console.log('Loading event.');

//
// index.js
//
// This is a tiny script designed for use with AWS Lambda which rebuilds the
// master branch of a target Travis repository.
//
// Please note that it must be configured with the values below before it can
// be used.
//

// Configuration
var repository = 'brandur/singularity';
var travisToken = '';

var https = require('https');

exports.handler = function (event, context) {
    console.log('Running event.');

    body='{"request": {"branch": "master"}}';

    var options = {
        host: 'api.travis-ci.org',
        port: 443,
        path: '/repo/' + encodeURIComponent(repository) + '/requests',
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Authorization': 'token ' + travisToken,
            'Content-Length': Buffer.byteLength(body),
            'Content-Type': 'application/json',
            'Travis-API-Version': '3',
        }
    }

    var req = https.request(options, function (res) {
        res.setEncoding('utf-8');

        // Collect response data as it comes back.
        var responseString = '';
        res.on('data', function (data) {
            responseString += data;
        });

        res.on('end', function () {
            console.log('Travis response status: ' + res.statusCode)
            console.log('Travis response: ' + responseString);

            if (res.statusCode >= 200 && res.statusCode < 300) {
                context.succeed('Rebuild executed successfully.');
            } else {
                context.fail('Rebuild failed with HTTP error.');
            }
        });
    });

    req.on('error', function (e) {
        console.error('Network error: ' + e.message);
        context.fail('Rebuild failed with network error.');
    });

    console.log('Making API request.')
    req.write(body);
    req.end();

    console.log('Finished running event.');
}
