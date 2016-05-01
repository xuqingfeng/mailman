(function () {
    if (window.WebSocket) {
        // get random port ?
        var loc = window.location, wsUri;
        if ('https' === loc.scheme) {
            wsUri = 'wss:';
        } else {
            wsUri = 'ws:';
        }
        wsUri += '//' + loc.host + '/api/wslog';
        console.info('wsUri', wsUri);
        var conn = new WebSocket(wsUri);
        conn.onopen = function () {
            console.info('ws opened');

            conn.send("Hi There !");
        };
        conn.onclose = function () {
            console.warn('ws closed');
        };

        conn.onmessage = function (evt) {
            console.log('receive data ', evt.data);
            $('.log').prepend(evt.data + '<br>')
        }
    }
})();

