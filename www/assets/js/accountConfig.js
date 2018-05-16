var idChangeForm = document.querySelector("#idChange"),
    displayNameChangeForm = document.querySelector("#displayNameChange"),
     passwordChangeForm = document.querySelector("#passwordChange"),
    idChangeFormError = idChangeForm.querySelector(".error");

idChangeForm.addEventListener("submit", idChangeSubmit, true);
displayNameChangeForm.addEventListener("submit", displayNameChangeSubmit, true);
passwordChangeForm.addEventListener("submit", passwordChangeSubmit, true);


function idChangeSubmit(e) {
    e.preventDefault();
    let targetForm = new FormData(idChangeForm);
    let changeObj = AccountApi.NewAccountObject(accountUuid);
    changeObj.name = targetForm.get("accountSettingsId")
    AccountApi.Change(changeObj).then(function (json) {
        if(!json.success){
            idChangeFormError.style.display = "block";
            idChangeFormError.style.innerText = json.message;
        }else{
            idChangeFormError.style.display = "none";
        }
    })
}

function displayNameChangeSubmit(e){
    e.preventDefault();

}

function passwordChangeSubmit(e){
    e.preventDefault();

}