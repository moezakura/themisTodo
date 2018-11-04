import AccountApi from "./api/AccountApi"
import BackView from "./backView"

class AccountConfig {
    constructor() {
        this.idChangeForm = document.querySelector("#idChange");
        if (this.idChangeForm === undefined || this.idChangeForm == null) return;
        return;

        this.displayNameChangeForm = document.querySelector("#displayNameChange");
        this.passwordChangeForm = document.querySelector("#passwordChange");
        this.idChangeFormError = this.idChangeForm.querySelector(".error");
        this.displayNameChangeFormError = this.displayNameChangeForm.querySelector(".error");
        this.passwordChangeFormError = this.passwordChangeForm.querySelector(".error");
        this.changeSuccess = document.querySelector(".success");
        this.idChangeSubmitFake = this.idChangeForm.querySelector(".fas.fa-check");
        this.idChangeSubmitButton = this.idChangeForm.querySelector("input[type=submit]");
        this.displayNameChangeSubmitFake = this.displayNameChangeForm.querySelector(".fas.fa-check");
        this.displayNameChangeSubmitButton = this.displayNameChangeForm.querySelector("input[type=submit]");
        this.currentPasswordDialog = document.querySelector("#currentPasswordDialog");

        this.backView = new BackView();
        this.backView.addWithHideElem(this.currentPasswordDialog);

        this.currentPasswordInput = this.currentPasswordDialog.querySelector("#currentPasswordInput");

        let dialogList = document.querySelectorAll(".error, .success");
        dialogList.forEach(function (value) {
            value.addEventListener("click", function (e) {
                this.style.display = "none";
            }, true);
        });

        document.AccountConfig = this;
        let that = this;

        this.idChangeForm.addEventListener("submit", function (e) {
            that.idChangeSubmit(e, this, that);
        }, true);
        this.displayNameChangeForm.addEventListener("submit", function (e) {
            that.displayNameChangeSubmit(e, this, that);
        }, true);
        this.passwordChangeForm.addEventListener("submit", function (e) {
            that.passwordChangeSubmit(e, this, that);
        }, true);
        this.idChangeSubmitFake.addEventListener("click", function () {
            that.idChangeSubmitButton.click();
        }, true);
        this.displayNameChangeSubmitFake.addEventListener("click", function () {
            that.displayNameChangeSubmitButton.click();
        }, true);
        this.currentPasswordDialog.addEventListener("submit", function (e) {
            that.passwordChangeConfirmSubmit(e, this, that);
        }, true);


        var selectPassword = "";
    }

    idChangeSubmit(e, _this, that) {
        e.preventDefault();
        let targetForm = new FormData(_this);
        let changeObj = AccountApi.NewAccountObject(accountUuid);
        changeObj.name = targetForm.get("accountSettingsId");
        AccountApi.Change(changeObj).then(function (json) {
            if (!json.success) {
                that.idChangeFormError.style.display = "block";
                that.idChangeFormError.innerText = json.message;
                that.changeSuccess.style.display = "none";
            } else {
                that.idChangeFormError.style.display = "none";
                that.changeSuccess.style.display = "block";
            }
        })
    }

    displayNameChangeSubmit(e, _this, that) {
        e.preventDefault();
        let targetForm = new FormData(_this);
        let changeObj = AccountApi.NewAccountObject(accountUuid);
        changeObj.displayName = targetForm.get("accountSettingsDisplayName");
        AccountApi.Change(changeObj).then(function (json) {
            if (!json.success) {
                that.displayNameChangeFormError.style.display = "block";
                that.displayNameChangeFormError.innerText = json.message;
                that.changeSuccess.style.display = "none";
            } else {
                that.displayNameChangeFormError.style.display = "none";
                that.changeSuccess.style.display = "block";
            }
        })
    }

    passwordChangeSubmit(e, _this, that) {
        e.preventDefault();
        let targetForm = new FormData(_this);
        let password = targetForm.get("accountSettingsPassword"),
            passwordRe = targetForm.get("accountSettingsPasswordRe");

        if (password !== passwordRe) {
            that.passwordChangeFormError.style.display = "block";
            that.passwordChangeFormError.innerText = "password is not match";
            that.changeSuccess.style.display = "none";
            return;
        }
        that.selectPassword = password;
        that.passwordChangeFormError.style.display = "none";
        that.changeSuccess.style.display = "none";

        that.currentPasswordDialog.style.display = "block";
        that.backView.show();

        that.currentPasswordInput.focus();
    }

    passwordChangeConfirmSubmit(e, _this, that) {
        e.preventDefault();
        if (that.currentPasswordDialog.style.display === "block") {
            let changeObj = AccountApi.NewAccountObject(accountUuid);
            changeObj.password = that.selectPassword;
            changeObj.currentPassword = that.currentPasswordInput.value;
            that.currentPasswordInput.value = "";
            AccountApi.Change(changeObj).then(function (json) {
                that.selectPassword = "";
                if (!json.success) {
                    that.passwordChangeFormError.style.display = "block";
                    that.passwordChangeFormError.innerText = json.message;
                    that.changeSuccess.style.display = "none";
                } else {
                    that.passwordChangeFormError.style.display = "none";
                    that.changeSuccess.style.display = "block";
                }

                that.backView.hide();
            })
        }
    }
}

import Vue from 'vue';
import IdSetting from "./components/AccountSettings/IdSetting.vue"

if(document.querySelector("#accountSettings")) {
    new Vue({
        delimiters: ['${', '}'],
        el: '#accountSettings',
        data: {
        },
        components: {
            IdSetting
        },
        created () {

        },
        methods : { },
    })
}


new AccountConfig();