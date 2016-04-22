new Vue({
    el: '#index',
    data: {
        subject: '',
        to: '',
        cc: '',
        from: '',
        priority: false,
        body: '',
        emails: [],
        sendClicked: false,
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
                console.info('priority', self.priority);
                if (self.sendClicked) {
                    swal({
                        title: "Send Again ?",
                        type: "warning",
                        showCancelButton: true,
                        confirmButtonText: "Yes, do it !",
                        cancelButtonText: "No, cancel it !",
                        closeOnConfirm: true,
                        closeOnCancel: true
                    }, function (isConfirm) {
                        if (isConfirm) {
                            var data = {
                                data: JSON.stringify({
                                    subject: self.subject,
                                    to: to,
                                    cc: cc,
                                    from: self.from,
                                    priority: self.priority,
                                    body: self.body
                                })
                            };
                            console.log('data', data);
                            $.post('/api/mail', data, function (json) {
                                    console.info(json);
                                    if (json.success) {
                                        self.previewLinkIsHidden = true;
                                        swal("Email Delivered !", json.msg, "success");
                                    } else {
                                        sweetAlert("Oops...", "send mail fail !", "error");
                                    }
                                })
                                .fail(function (err) {
                                    console.error(err);
                                    sweetAlert("Oops...", "request fail !", "error");
                                });
                        }
                    });
                } else {
                    self.sendClicked = true;
                    var data = {
                        data: JSON.stringify({
                            subject: self.subject,
                            to: to,
                            cc: cc,
                            from: self.from,
                            priority: self.priority,
                            body: self.body
                        })
                    };
                    console.log('data', data);
                    $.post('/api/mail', data, function (json) {
                            console.info(json);
                            if (json.success) {
                                self.previewLinkIsHidden = true;
                                swal("Email Delivered !", json.msg, "success");
                            } else {
                                sweetAlert("Oops...", "send mail fail !", "error");
                            }
                        })
                        .fail(function (err) {
                            console.error(err);
                            sweetAlert("Oops...", "request fail !", "error");
                        });
                }

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