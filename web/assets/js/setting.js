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
        smtpClicked: false
    },
    ready: function () {
        var self = this;
        $.get('/api/account', function (json) {
                console.info('account', json);
                if (json.success) {
                    self.emails = json.data;
                }
            })
            .fail(function (err) {
                console.error(err);
            });
        $.get('/api/contacts', function (json) {
                console.info('contacts', json);
                if (json.success) {
                    self.contacts = json.data;
                }
            })
            .fail(function (err) {
                console.error(err);
            });
        $.get('/api/smtpServer', function (json) {
                console.info('smtpServer', json);
                if (json.success) {
                    self.servers = json.data;
                }
            })
            .fail(function (err) {
                console.error(err);
            });
    },
    methods: {
        saveAccount: function () {
            var self = this;
            if (!checkParams('saveAccount', self)) {
                sweetAlert("Oops...", "missing info !", "error");
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
                        sweetAlert("Oops...", err, "error");
                    });
            }
        },
        manageEmail: function (email) {
            var self = this;
            console.info('manage: ', email);
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
            console.info('accountEmail ', self.accountEmail);
            if (!checkParams('deleteAccount', self)) {
                sweetAlert("Oops...", "missing info !", "error");
            } else {
                $.ajax({
                    url: '/api/account',
                    type: 'DELETE',
                    data: JSON.stringify({key: self.accountEmail}),
                    success: function (json) {
                        console.info('deleteAccount', json);
                        if (json.success) {
                            self.accountEmail = '';
                            self.accountPassword = '';
                            self.emails = json.data;
                            self.emailClicked = false;
                        }
                    },
                    error: function (err) {
                        console.error(err);
                        sweetAlert("Oops...", err, "error");
                    }
                });
            }
        },
        saveContacts: function () {
            var self = this;
            if (!checkParams('saveContacts', self)) {
                sweetAlert("Oops...", "missing info !", "error");
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
                    sweetAlert("Oops...", err, "error");
                })
            }
        },
        manageContacts: function (contacts) {
            var self = this;
            console.info('manage contacts ', contacts);
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
            console.info('delete contacts ', self.contactsEmail)
            if (!checkParams('deleteContacts', self)) {
                sweetAlert("Oops...", "missing info !", "error");
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
                        sweetAlert("Oops...", err, "error");
                    }
                });
            }
        },
        saveSmtp: function () {
            var self = this;
            if (!checkParams('saveSmtp', self)) {
                sweetAlert("Oops...", "missing info !", "error");
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
                    sweetAlert("Oops...", err, "error");
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
                sweetAlert("Oops...", "missing info !", "error");
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
                        sweetAlert("Oops...", err, "error");
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

