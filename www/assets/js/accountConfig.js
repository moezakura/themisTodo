var idChangeForm = document.querySelector("#idChange"),
    displayNameChangeForm = document.querySelector("#displayNameChange"),
    passwordChangeForm = document.querySelector("#passwordChange"),
    idChangeFormError = idChangeForm.querySelector(".error"),
    displayNameChangeFormError = displayNameChangeForm.querySelector(".error"),
    passwordChangeFormError = passwordChangeForm.querySelector(".error"),
    changeSuccess = document.querySelector(".success"),
    idChangeSubmitFake = idChangeForm.querySelector(".fas.fa-check"),
    idChangeSubmitButton = idChangeForm.querySelector("input[type=submit]"),
    displayNameChangeSubmitFake = displayNameChangeForm.querySelector(".fas.fa-check"),
    displayNameChangeSubmitButton = displayNameChangeForm.querySelector("input[type=submit]"),
    currentPasswordDialog = document.querySelector("#currentPasswordDialog"),
    backViewLayer = document.querySelector("#backViewLayer"),
    currentPasswordInput = currentPasswordDialog.querySelector("#currentPasswordInput");

var dialogList = document.querySelectorAll(".error, .success");
dialogList.forEach(function (value) {
    value.addEventListener("click", function (e) {
        this.style.display = "none";
    }, true);
});

idChangeForm.addEventListener("submit", idChangeSubmit, true);
displayNameChangeForm.addEventListener("submit", displayNameChangeSubmit, true);
passwordChangeForm.addEventListener("submit", passwordChangeSubmit, true);
idChangeSubmitFake.addEventListener("click", function () {
    idChangeSubmitButton.click();
}, true);
displayNameChangeSubmitFake.addEventListener("click", function () {
    displayNameChangeSubmitButton.click();
}, true);
backViewLayer.addEventListener("click", backViewLayerClick, true);
currentPasswordDialog.addEventListener("submit", passwordChangeConfirmSubmit, true);


var selectPassword = "";

function idChangeSubmit(e) {
    e.preventDefault();
    let targetForm = new FormData(this);
    let changeObj = AccountApi.NewAccountObject(accountUuid);
    changeObj.name = targetForm.get("accountSettingsId");
    AccountApi.Change(changeObj).then(function (json) {
        if (!json.success) {
            idChangeFormError.style.display = "block";
            idChangeFormError.innerText = json.message;
            changeSuccess.style.display = "none";
        } else {
            idChangeFormError.style.display = "none";
            changeSuccess.style.display = "block";
        }
    })
}

function displayNameChangeSubmit(e) {
    e.preventDefault();
    let targetForm = new FormData(this);
    let changeObj = AccountApi.NewAccountObject(accountUuid);
    changeObj.displayName = targetForm.get("accountSettingsDisplayName");
    AccountApi.Change(changeObj).then(function (json) {
        if (!json.success) {
            displayNameChangeFormError.style.display = "block";
            displayNameChangeFormError.innerText = json.message;
            changeSuccess.style.display = "none";
        } else {
            displayNameChangeFormError.style.display = "none";
            changeSuccess.style.display = "block";
        }
    })
}

function backViewLayerClick() {
    currentPasswordDialog.style.display = "none";
    backViewLayer.style.display = "none";
}

function passwordChangeSubmit(e) {
    e.preventDefault();
    let targetForm = new FormData(this);
    let password = targetForm.get("accountSettingsPassword"),
        passwordRe = targetForm.get("accountSettingsPasswordRe");

    if (password !== passwordRe) {
        passwordChangeFormError.style.display = "block";
        passwordChangeFormError.innerText = "password is not match";
        changeSuccess.style.display = "none";
        return;
    }
    selectPassword = password;
    passwordChangeFormError.style.display = "none";
    changeSuccess.style.display = "none";

    currentPasswordDialog.style.display = "block";
    backViewLayer.style.display = "block";

    currentPasswordInput.focus();
}

function passwordChangeConfirmSubmit(e) {
    e.preventDefault();
    if (currentPasswordDialog.style.display == "block") {
        let changeObj = AccountApi.NewAccountObject(accountUuid);
        changeObj.password = selectPassword;
        changeObj.currentPassword = currentPasswordInput.value;
        currentPasswordInput.value = "";
        AccountApi.Change(changeObj).then(function (json) {
            selectPassword = "";
            if (!json.success) {
                passwordChangeFormError.style.display = "block";
                passwordChangeFormError.innerText = json.message;
                changeSuccess.style.display = "none";
            } else {
                passwordChangeFormError.style.display = "none";
                changeSuccess.style.display = "block";
            }

            backViewLayerClick();
        })
    }
}