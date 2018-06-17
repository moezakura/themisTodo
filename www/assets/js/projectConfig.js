import BackView from "./backView"
import ProjectUtils from "./projectUtils";
import UserSearchDialog from "./userSearchDialog"
import ProjectMemberList from "./projectMemberList"

class ProjectConfig {
    constructor() {
        this.taskConfigShow = document.querySelector("#taskboardConfig");
        if (this.taskConfigShow === undefined || this.taskConfigShow == null)
            return;

        this.taskBackView = new BackView();
        this.projectConfigPopup = document.querySelector("#projectConfigPopup");
        this.postTaskConfigForm = document.querySelector("#projectConfigForm");
        this.taskboardTitleElem = document.querySelector("#taskboardTitle>span");
        this.projectConfigPopupErrorElem = this.projectConfigPopup.querySelector(".error");

        ProjectMemberList.update();
        let that = this;

        this.postTaskConfigForm.addEventListener("submit", function(e){
            that.postTaskConfig(e, that);
        }, true);
        this.taskConfigShow.addEventListener("click", function (e){
            that.taskConfigShowClick(e, that)
        }, true);
        this.projectConfigPopupErrorElem.addEventListener("click", ProjectUtils.clickObjectHide, true);
        this.taskBackView.addWithHideElem(this.projectConfigPopup);
        this.taskBackView.addHideEvent(function () {
            that.userSearchDialog.hide();
        });

        let userSelectInput = document.querySelector("#userSelect");

        this.userSearchDialog = new UserSearchDialog(userSelectInput, {
            "forceSubmit": function (sendUuid) {
                let projectAddMemberJson = {
                    "uuid": sendUuid
                };

                fetch("/project/addUser/" + projectId, {
                    method: "POST",
                    body: JSON.stringify(projectAddMemberJson),
                    credentials: "same-origin"
                }).then(function (response) {
                    return response.json();
                }).then(function (json) {
                    if (!json.success) {
                        that.projectConfigPopupErrorElem.style.display = "block";
                        that.projectConfigPopupErrorElem.innerText = json.message;
                    } else {
                        that.projectConfigPopupErrorElem.style.display = "none";
                        userSelectInput.value = "";
                        memberList.push(json.addedAccount);
                        ProjectMemberList.update();
                    }
                });
            }
        });

        document.querySelector("body").addEventListener("keydown", function (e) {
            if (e.keyCode === 27 && that.projectConfigPopup.style.display === "block")
                that.taskBackView.hide();
        }, true);
    }

    taskConfigShowClick(e, that) {
        e.preventDefault();
        that.projectConfigPopup.style.display = "block";
        that.taskBackView.show();
    }

    postTaskConfig(e, that) {
        e.preventDefault();

        let formData = new FormData(that.postTaskConfigForm);
        let projectNewName = formData.get("name");
        let postTaskConfigJson = {
            "name": projectNewName,
            "description": formData.get("description"),
        };

        fetch("/project/update/" + projectId, {
            method: 'POST',
            body: JSON.stringify(postTaskConfigJson),
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        }).then(function (json) {
            if (!json.success) {
                that.projectConfigPopupErrorElem.style.display = "block";
                that.projectConfigPopupErrorElem.innerText = json.message;
            } else {
                that.projectConfigPopupErrorElem.style.display = "none";
                that.taskboardTitleElem.innerText = projectNewName;
            }
        });
    }
}

new ProjectConfig();
