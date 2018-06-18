import BackView from "./backView"

class ProjectDelete {
    constructor() {
        this.deleteProject = document.querySelector("#deleteProject");
        if(this.deleteProject === undefined || this.deleteProject == null)
            return;

        this.deleteProjectPopup = document.querySelector("#deleteProjectPopup");
        this.deleteProjectPopupCloseButton = document.querySelector("#deleteProjectPopup .close");
        this.deleteProjectNameInput = document.querySelector("#deleteProjectPopup input[type=text]");
        this.deleteProjectNameSubmitButton = document.querySelector("#deleteProjectPopup input[type=submit]");
        this.deleteProjectPopupError = this.deleteProjectPopup.querySelector(".error");
        this.projectDeleteBackView = new BackView();
        this.projectDeleteBackView.addWithHideElem(this.deleteProjectPopup);

        let that = this;

        this.deleteProject.addEventListener("click", function(e){
            that.deleteProjectPopupShow(e, that);
        }, true);
        this.deleteProjectPopupCloseButton.addEventListener("click", function(e){
            that.deleteProjectPopupClose(e, that);
        }, true);
        this.deleteProjectPopup.addEventListener("submit", function(e){
            that.deleteProjectPopupSend(e, that);
        }, true);
        this.deleteProjectNameInput.addEventListener("keyup", function(e){
            that.deleteProjectNameInputKeyUp(e, that);
        }, true);
        this.deleteProjectPopupError.addEventListener("click", function(e){
            that.deleteProjectPopupErrorClick(e, that);
        }, true);
    }

    deleteProjectPopupShow(e, that) {
        that.deleteProjectPopup.style.display = "block";
        that.projectDeleteBackView.show();
        if (that.taskBackView !== undefined && that.taskBackView != null)
            that.taskBackView.hide();
    }

    deleteProjectPopupClose(e, that) {
        that.projectDeleteBackView.hide();
    }

    deleteProjectPopupErrorClick(e, that) {
        that.deleteProjectPopupErrorHide();
    }

    deleteProjectPopupErrorHide(e, that) {
        that.deleteProjectPopupError.style.display = "none";
    }

    deleteProjectPopupErrorShow(text, that) {
        that.deleteProjectPopupError.style.display = "block";
        if (text !== undefined && text != null) {
            that.deleteProjectPopupError.innerText = text;
        }
    }

    deleteProjectNameInputKeyUp(e, that) {
        let typeText = that.deleteProjectNameInput.value;
        that.deleteProjectNameSubmitButton.disabled = (typeText !== projectName);
    }

    deleteProjectPopupSend(e, that) {
        e.preventDefault();
        fetch("/project/delete/" + projectId, {
            method: 'POST',
            body: "",
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        }).then(function (json) {
            if (json.success) location.href = "/home";
            else that.deleteProjectPopupErrorShow(json.message, that);
        });
    }
}

new ProjectDelete();