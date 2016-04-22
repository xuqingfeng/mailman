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
        previewLinkIsHidden: true,
        lang: 'en',
        i18n: {}
    },
    ready: function () {
        var self = this;

        // default
        applyLangSet(self);
        $.get('/api/lang', function (json) {
            if (json.success) {
                self.lang = json.data;
                i18next.changeLanguage(self.lang);
                applyLangSet(self);
            }
        });
        $.get('/api/account', function (json) {
                if (json.success && json.data) {
                    self.emails = json.data;
                    self.from = self.emails[0];
                }
            })
            .fail(function (err) {
                console.error(err);
                swal(self.i18n.oops, err, "error");
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
                // swal
                swal(self.i18n.oops, self.i18n.missing_info, "error");
            } else {
                console.info('priority', self.priority);
                if (self.sendClicked) {
                    swal({
                        title: self.i18n.send_again,
                        type: "warning",
                        showCancelButton: true,
                        confirmButtonText: self.i18n.yes_do_it,
                        cancelButtonText: self.i18n.no_cancel_it,
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
                                        swal(self.i18n.email_delivered, json.msg, "success");
                                    } else {
                                        swal(self.i18n.oops, self.i18n.send_email_fail, "error");
                                    }
                                })
                                .fail(function (err) {
                                    console.error(err);
                                    swal(self.i18n.oops, self.i18n.request_fail, "error");
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
                                swal(self.i18n.email_delivered, json.msg, "success");
                            } else {
                                swal(self.i18n.oops, self.i18n.send_email_fail, "error");
                            }
                        })
                        .fail(function (err) {
                            console.error(err);
                            swal(self.i18n.oops, self.i18n.request_fail, "error");
                        });
                }

            }
        },
        preview: function () {
            var self = this;
            if (!self.body) {
                console.error('empty');
                swal(self.i18n.oops, self.i18n.missing_info, "error");
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
                        swal(self.i18n.oops, self.i18n.request_fail, "error");
                    });
            }
        }
    }
});

function applyLangSet(self) {

    // i18n
    self.i18n = {
        subject: i18next.t('subject'),
        to: i18next.t('to'),
        cc: i18next.t('cc'),
        from: i18next.t('from'),
        body: i18next.t('body'),
        support_markdown_syntax: i18next.t('support_markdown_syntax'),
        high_priority: i18next.t('high_priority'),
        preview: i18next.t('preview'),
        send: i18next.t('send'),
        go_to_preview: i18next.t('go_to_preview'),
        index: i18next.t('index'),
        setting: i18next.t('setting'),
        log: i18next.t('log'),

        oops: i18next.t('oops'),
        missing_info: i18next.t('missing_info'),
        send_again: i18next.t('send_again'),
        yes_do_it: i18next.t('yes_do_it'),
        no_cancel_it: i18next.t('no_cancel_it'),
        email_delivered: i18next.t('email_delivered'),
        request: i18next.t('request'),
        send_email_fail: i18next.t('send_email_fail')
    };
}