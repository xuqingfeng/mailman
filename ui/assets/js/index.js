var myDropzone = {},
    toContactsLoaded = false,
    ccContactsLoaded = false;
Dropzone.autoDiscover = false;
var vue = new Vue({
    el: '#index',
    data: {
        subject: '',
        to: [],
        cc: [],
        from: '',
        priority: false,
        body: '',
        emails: [],
        sendClicked: false,
        previewLinkIsHidden: true,
        lang: 'en',
        i18n: {},
        token: Date.now().toString()
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

                // dorpzone
                myDropzone = new Dropzone("form#attachment", {
                    url: "/api/file",
                    addRemoveLinks: true,
                    filesizeBase: 1024,
                    uploadMultiple: true,
                    maxFiles: 5,
                    //myDropzone.processQueue();
                    //autoProcessQueue: false,
                    dictDefaultMessage: self.i18n.add_attachment,
                    dictRemoveFile: self.i18n.remove_file
                });
                // todo save attachment fail
                myDropzone.on('sending', function (file, xhr, formData) {
                    formData.append('token', self.token);
                });
            }
        });
        $.get('/api/account', function (json) {
                if (json.success && json.data) {
                    self.emails = json.data;
                    self.from = self.emails[0];
                }
            })
            .fail(function (err) {
                swal(self.i18n.oops, err, "error");
            });
    },
    methods: {
        send: function () {
            var self = this;

            self.to = $('#to').val();
            self.cc = $('#cc').val();

            if (!self.subject || !self.from || !self.body || self.to.length < 1) {
                swal(self.i18n.oops, self.i18n.missing_info, "error");
            } else {
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
                            var data = JSON.stringify({
                                    subject: self.subject,
                                    to: self.to,
                                    cc: self.cc,
                                    from: self.from,
                                    priority: self.priority,
                                    body: self.body,
                                    token: self.token,
                                    attachmentFileNames: getAcceptedFileNames()
                                });
                            $.post('/api/mail', data, function (json) {
                                    if (json.success) {
                                        self.previewLinkIsHidden = true;
                                        swal(self.i18n.email_delivered, json.msg, "success");
                                    } else {
                                        swal(self.i18n.oops, self.i18n.send_email_fail, "error");
                                    }
                                })
                                .fail(function (err) {
                                    swal(self.i18n.oops, self.i18n.request_fail, "error");
                                });
                        }
                    });
                } else {
                    self.sendClicked = true;
                    var data = JSON.stringify({
                            subject: self.subject,
                            to: self.to,
                            cc: self.cc,
                            from: self.from,
                            priority: self.priority,
                            body: self.body,
                            token: self.token,
                            attachmentFileNames: getAcceptedFileNames()
                        });
                    $.post('/api/mail', data, function (json) {
                            if (json.success) {
                                self.previewLinkIsHidden = true;
                                swal(self.i18n.email_delivered, json.msg, "success");
                            } else {
                                swal(self.i18n.oops, self.i18n.send_email_fail, "error");
                            }
                        })
                        .fail(function (err) {
                            swal(self.i18n.oops, self.i18n.request_fail, "error");
                        });
                }

            }
        },
        preview: function () {
            var self = this;
            if (!self.body) {
                swal(self.i18n.oops, self.i18n.missing_info, "error");
            } else {
                $.post('/api/preview', JSON.stringify({
                            body: self.body
                        }), function (json) {
                        if (json.success) {
                            self.previewLinkIsHidden = false;
                        }
                    })
                    .fail(function (err) {
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
        add_attachment: i18next.t('add_attachment'),
        remove_file: i18next.t('remove_file'),

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

function getAcceptedFileNames() {
    var acceptedFiles = myDropzone.getAcceptedFiles();
    var acceptedFileNames = [];
    acceptedFiles.filter(function (f) {
        acceptedFileNames.push(f.name);
    });

    return acceptedFileNames;
}

$(function () {
    var REGEX_EMAIL = '([a-z0-9!#$%&\'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&\'*+/=?^_`{|}~-]+)*@' +
        '(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)';

    $('#to').selectize({
        persist: false,
        maxItems: null,
        valueField: 'email',
        labelField: 'name',
        searchField: ['name', 'email'],
        options: [],
        render: {
            item: function (item, escape) {
                return '<div>' +
                    (item.name ? '<span class="name">' + escape(item.name) + '</span> ' : '') +
                    (item.email ? '<span class="email">' + escape(item.email) + '</span>' : '') +
                    '</div>';
            },
            option: function (item, escape) {
                var label = item.name || item.email;
                var caption = item.name ? item.email : null;
                return '<div>' +
                    '<span class="prefix">' + escape(label) + '</span> ' +
                    (caption ? '<span class="caption">' + escape(caption) + '</span>' : '') +
                    '</div>';
            }
        },
        load: function (query, callback) {
            $.ajax({
                url: '/api/contacts',
                type: 'GET',
                error: function () {
                    callback();
                },
                success: function (json) {
                    callback(json.data);
                    toContactsLoaded = true;
                }
            });
        },
        createFilter: function (input) {
            var match, regex;

            // email@address.com
            regex = new RegExp('^' + REGEX_EMAIL + '$', 'i');
            match = input.match(regex);
            if (match) return !this.options.hasOwnProperty(match[0]);

            // name <email@address.com>
            regex = new RegExp('^([^<]*)\<' + REGEX_EMAIL + '\>$', 'i');
            match = input.match(regex);
            if (match) return !this.options.hasOwnProperty(match[2]);

            return false;
        },
        create: function (input) {
            if ((new RegExp('^' + REGEX_EMAIL + '$', 'i')).test(input)) {
                return {email: input};
            }
            var match = input.match(new RegExp('^([^<]*)\<' + REGEX_EMAIL + '\>$', 'i'));
            if (match) {
                return {
                    email: match[2],
                    name: $.trim(match[1])
                };
            }
            //alert('Invalid email address.');
            return false;
        }
    });

    $('#cc').selectize({
        persist: false,
        maxItems: null,
        valueField: 'email',
        labelField: 'name',
        searchField: ['name', 'email'],
        options: [],
        render: {
            item: function (item, escape) {
                return '<div>' +
                    (item.name ? '<span class="name">' + escape(item.name) + '</span> ' : '') +
                    (item.email ? '<span class="email">' + escape(item.email) + '</span>' : '') +
                    '</div>';
            },
            option: function (item, escape) {
                var label = item.name || item.email;
                var caption = item.name ? item.email : null;
                return '<div>' +
                    '<span class="prefix">' + escape(label) + '</span> ' +
                    (caption ? '<span class="caption">' + escape(caption) + '</span>' : '') +
                    '</div>';
            }
        },
        load: function (query, callback) {
            $.ajax({
                url: '/api/contacts',
                type: 'GET',
                error: function () {
                    callback();
                },
                success: function (json) {
                    callback(json.data);
                    ccContactsLoaded = true;
                }
            });
        },
        createFilter: function (input) {
            var match, regex;

            // email@address.com
            regex = new RegExp('^' + REGEX_EMAIL + '$', 'i');
            match = input.match(regex);
            if (match) return !this.options.hasOwnProperty(match[0]);

            // name <email@address.com>
            regex = new RegExp('^([^<]*)\<' + REGEX_EMAIL + '\>$', 'i');
            match = input.match(regex);
            if (match) return !this.options.hasOwnProperty(match[2]);

            return false;
        },
        create: function (input) {
            if ((new RegExp('^' + REGEX_EMAIL + '$', 'i')).test(input)) {
                return {email: input};
            }
            var match = input.match(new RegExp('^([^<]*)\<' + REGEX_EMAIL + '\>$', 'i'));
            if (match) {
                return {
                    email: match[2],
                    name: $.trim(match[1])
                };
            }
            //alert('Invalid email address.');
            return false;
        }
    });
});