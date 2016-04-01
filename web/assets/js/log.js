(function () {
    if (window.WebSocket) {
        var conn = new WebSocket('ws://localhost:8000/api/wslog');
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

