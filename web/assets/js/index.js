new Vue({
    el: '#index',
    data: {
        subject: '',
        to: '',
        cc: '',
        from: '',
        body: '',
        emails: [],
        previewLinkIsHidden: true
    },
    ready: function () {
        var self = this;
        $.get('/api/account', function (json) {
                if (json.success && json.data) {
                    self.emails = json.data;
                    self.from = self.emails[0];
                }
            })
            .fail(function (err) {
                console.error(err);
                sweetAlert("Oops...", err, "error");
            });
    },
    methods: {
        send: function () {
            var self = this;
            var to = self.to.split(',').filter(function (n) {
                    return n;
                }),
                cc = self.cc.split(',').filter(function (n) {
                    return n;
                });
            if (!self.subject || !self.from || !self.body || to.length < 1) {
                // sweetAlert
                console.error('empty');
                sweetAlert("Oops...", "missing info !", "error");
            } else {
                var data = {
                    data: JSON.stringify({
                        subject: self.subject,
                        to: to,
                        cc: cc,
                        from: self.from,
                        body: self.body
                    })
                };
                $.post('/api/mail', data, function (json) {
                        console.info(json);
                        if (json.success) {
                            self.previewLinkIsHidden = true;
                            swal("Awesome", json.msg, "success");
                        } else {
                            sweetAlert("Oops...", "send mail fail !", "error");
                        }
                    })
                    .fail(function (err) {
                        console.error(err);
                        sweetAlert("Oops...", "request fail !", "error");
                    });
            }
        },
        preview: function () {
            var self = this;
            if (!self.body) {
                console.error('empty');
                sweetAlert("Oops...", "missing info !", "error");
            } else {
                $.post('/api/preview', {
                        data: JSON.stringify({
                            body: self.body
                        })
                    }, function (json) {
                        console.info(json);
                        if (json.success) {
                            self.previewLinkIsHidden = false;
                        }
                    })
                    .fail(function (err) {
                        console.error(err);
                        sweetAlert("Oops...", "request fail !", "error");
                    });
            }
        }
    }
});