var deleteProject = document.querySelector("#deleteProject"),
    deleteProjectPopup = document.querySelector("#deleteProjectPopup"),
    deleteProjectPopupCloseButton = document.querySelector("#deleteProjectPopup .close"),
    deleteProjectNameInput = document.querySelector("#deleteProjectPopup input[type=text]"),
    deleteProjectNameSubmitButton = document.querySelector("#deleteProjectPopup input[type=submit]"),
    deleteProjectPopupError = deleteProjectPopup.querySelector(".error"),
    projectDeleteBackView = new BackView();
projectDeleteBackView.addWithHideElem(deleteProjectPopup);

deleteProject.addEventListener("click", deleteProjectPopupShow, true);
deleteProjectPopupCloseButton.addEventListener("click", deleteProjectPopupClose, true);
deleteProjectPopup.addEventListener("submit", deleteProjectPopupSend, true);
deleteProjectNameInput.addEventListener("keyup", deleteProjectNameInputKeyUp, true);
deleteProjectPopupError.addEventListener("click", deleteProjectPopupErrorClick, true);

function deleteProjectPopupShow() {
    deleteProjectPopup.style.display = "block";
    projectDeleteBackView.show();
    if (taskBackView !== undefined && taskBackView != null)
        taskBackView.hide();
}

function deleteProjectPopupClose() {
    projectDeleteBackView.hide();
}

function deleteProjectPopupErrorClick() {
    deleteProjectPopupErrorHide();
}

function deleteProjectPopupErrorHide() {
    deleteProjectPopupError.style.display = "none";
}

function deleteProjectPopupErrorShow(text) {
    deleteProjectPopupError.style.display = "block";
    if (text !== undefined && text != null) {
        deleteProjectPopupError.innerText = text;
    }
}

function deleteProjectNameInputKeyUp() {
    let typeText = deleteProjectNameInput.value;
    deleteProjectNameSubmitButton.disabled = (typeText !== projectName);
}

function deleteProjectPopupSend(e) {
    e.preventDefault();
    fetch("/project/delete/" + projectId, {
        method: 'POST',
        body: "",
        credentials: "same-origin"
    }).then(function (response) {
        return response.json();
    }).then(function (json) {
        if (json.success) location.href = "/home";
        else deleteProjectPopupErrorShow(json.message);
    });
}