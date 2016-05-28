new Vue({
    el: '#setting',
    data: {
        accountEmail: '',
        accountPassword: '',
        contactsName: '',
        contactsEmail: '',
        smtpAddress: '',
        smtpServer: '',
        emails: [],
        contacts: [],
        servers: [],
        emailClicked: false,
        contactsClicked: false,
        smtpClicked: false,
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
                if (json.success) {
                    self.emails = json.data;
                }
            })
            .fail(function (err) {
                console.error(err);
            });
        $.get('/api/contacts', function (json) {
                if (json.success) {
                    self.contacts = json.data;
                }
            })
            .fail(function (err) {
                console.error(err);
            });
        $.get('/api/smtpServer', function (json) {
                if (json.success) {
                    self.servers = json.data;
                }
            })
            .fail(function (err) {
                console.error(err);
            });
    },
    watch: {
        'lang': function (val, oldVal) {
            var self = this;
            $.post('/api/lang', {
                data: JSON.stringify({type: val})
            }, function (json) {
                if (json.success) {
                    self.lang = val;
                    i18next.changeLanguage(self.lang);
                    applyLangSet(self);
                }
            });
        }
    },
    methods: {
        saveAccount: function () {
            var self = this;
            if (!checkParams('saveAccount', self)) {
                swal(self.i18n.oops, self.i18n.missing_info, "error");
            } else {
                $.post('/api/account', {
                        data: JSON.stringify({
                            email: self.accountEmail,
                            password: self.accountPassword
                        })
                    }, function (json) {
                        if (json.success && json.data) {
                            self.accountEmail = '';
                            self.accountPassword = '';
                            self.emails = json.data;
                        }
                    })
                    .fail(function (err) {
                        console.error(err);
                        swal(self.i18n.oops, err, "error");
                    });
            }
        },
        manageEmail: function (email) {
            var self = this;
            if (!self.emailClicked) {
                self.accountEmail = email;
                self.emailClicked = !self.emailClicked;
            } else {
                self.accountEmail = '';
                self.emailClicked = !self.emailClicked;
            }
        },
        deleteAccount: function () {
            var self = this;
            if (!checkParams('deleteAccount', self)) {
                swal(self.i18n.oops, self.i18n.missing_info, "error");
            } else {
                $.ajax({
                    url: '/api/account',
                    type: 'DELETE',
                    data: JSON.stringify({key: self.accountEmail}),
                    success: function (json) {
                        if (json.success) {
                            self.accountEmail = '';
                            self.accountPassword = '';
                            self.emails = json.data;
                            self.emailClicked = false;
                        }
                    },
                    error: function (err) {
                        console.error(err);
                        swal(self.i18n.oops, err, "error");
                    }
                });
            }
        },
        saveContacts: function () {
            var self = this;
            if (!checkParams('saveContacts', self)) {
                swal(self.i18n.oops, self.i18n.missing_info, "error");
            } else {
                $.post('/api/contacts', {
                    data: JSON.stringify({
                        name: self.contactsName,
                        email: self.contactsEmail
                    })
                }, function (json) {
                    if (json.success && json.data) {
                        self.contactsName = '';
                        self.contactsEmail = '';
                        self.contacts = json.data;
                    }
                }).fail(function (err) {
                    console.error(err);
                    swal(self.i18n.oops, err, "error");
                })
            }
        },
        manageContacts: function (contacts) {
            var self = this;
            if (!self.contactsClicked) {
                self.contactsEmail = contacts.email;
                self.contactsName = contacts.name;
                self.contactsClicked = !self.contactsClicked;
            } else {
                self.contactsEmail = '';
                self.contactsName = '';
                self.contactsClicked = !self.contactsClicked;
            }
        },
        deleteContacts: function () {
            var self = this;
            if (!checkParams('deleteContacts', self)) {
                swal(self.i18n.oops, self.i18n.missing_info, "error");
            } else {
                $.ajax({
                    url: '/api/contacts',
                    type: 'DELETE',
                    data: JSON.stringify({key: self.contactsEmail}),
                    success: function (json) {
                        if (json.success) {
                            self.contactsName = '';
                            self.contactsEmail = '';
                            self.contacts = json.data;
                            self.contactsClicked = false;
                        }
                    },
                    error: function (err) {
                        swal(self.i18n.oops, err, "error");
                    }
                });
            }
        },
        saveSmtp: function () {
            var self = this;
            if (!checkParams('saveSmtp', self)) {
                swal(self.i18n.oops, self.i18n.missing_info, "error");
            } else {
                $.post('/api/smtpServer', {
                    data: JSON.stringify({
                        address: self.smtpAddress,
                        server: self.smtpServer
                    })
                }, function (json) {
                    if (json.success && json.data) {
                        self.smtpAddress = '';
                        self.smtpServer = '';
                        self.servers = json.data;
                    }
                }).fail(function (err) {
                    console.error(err);
                    swal(self.i18n.oops, err, "error");
                })
            }
        },
        manageSmtp: function (smtp) {
            var self = this;
            if (!self.smtpClicked) {
                self.smtpAddress = smtp.address;
                self.smtpServer = smtp.server;
                self.smtpClicked = !self.smtpClicked;
            } else {
                self.smtpAddress = '';
                self.smtpServer = '';
                self.smtpClicked = !self.smtpClicked;
            }
        },
        deleteSmtp: function () {
            var self = this;
            if (!checkParams('deleteSmtp', self)) {
                swal(self.i18n.oops, self.i18n.missing_info, "error");
            } else {
                $.ajax({
                    url: '/api/smtpServer',
                    type: 'DELETE',
                    data: JSON.stringify({key: self.smtpAddress}),
                    success: function (json) {
                        if (json.success) {
                            self.smtpAddress = '';
                            self.smtpServer = '';
                            self.servers = json.data;
                            self.smtpClicked = !self.smtpClicked;
                        }
                    },
                    error: function (err) {
                        swal(self.i18n.oops, err, "error");
                    }
                });
            }
        }
    }
});

function checkParams(type, self) {

    switch (type) {
        case 'saveAccount':

            // simplify
            return !(!self.accountEmail || !self.accountPassword);
            break;
        case 'saveContacts':

            return !(!self.contactsName || !self.contactsEmail);
            break;
        case 'saveSmtp':

            return !(!self.smtpAddress || !self.smtpServer);
            break;
        case 'deleteAccount':

            return !(!self.accountEmail);
            break;
        case 'deleteContacts':

            return !(!self.contactsEmail);
            break;
        case 'deleteSmtp':

            return !(!self.smtpAddress);
            break;
        default:
            return false;
    }
}

function applyLangSet(self) {
    // i18n
    self.i18n = {
        account: i18next.t('account'),
        email: i18next.t('email'),
        password: i18next.t('password'),
        contacts: i18next.t('contacts'),
        name: i18next.t('name'),
        advanced_settings: i18next.t('advanced_settings'),
        lang: i18next.t('lang'),
        custom_smtp_server: i18next.t('custom_smtp_server'),
        mail_address: i18next.t('mail_address'),
        smtp_server: i18next.t('smtp_server'),
        save: i18next.t('save'),
        delete: i18next.t('delete'),
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

