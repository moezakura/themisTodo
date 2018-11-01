import TaskDetail from "./taskDetail"
import TaskApi from "./taskApi";
import UserSearchDialog from "./userSearchDialog";
import LoadingView from "./loadingView";
import ProjectUtils from "./projectUtils";

class TaskEdit {
    constructor() {
        this.taskPopup = document.querySelector("#taskPopup");
        if (this.taskPopup === undefined || this.taskPopup == null)
            return;
        this.taskPopupEditButton = this.taskPopup.querySelector("#taskPopupEditButton");
        this.taskAssgin = this.taskPopup.querySelector("#taskPopupAssign");
        this.taskPopupAssignIcon = this.taskPopup.querySelector("#taskPopupAssignIcon");
        this.taskPopupEditCancelButton = this.taskPopup.querySelector("#taskPopupEditCancelButton");
        this.error = this.taskPopup.querySelector(".error");
        this.success = this.taskPopup.querySelector(".success");
        this.taskForm = {
            "name": this.taskPopup.querySelector("#taskPopupTitle"),
            "title": this.taskPopup.querySelector("#taskPopupTitle"),
            "assign": this.taskPopup.querySelector("#taskPopupAssign"),
            "deadline": this.taskPopup.querySelector("#taskPopupDeadlineChange"),
            "description": this.taskPopup.querySelector("#taskPopupDescription"),
        };

        this.assignUuid = -1;

        let that = this;

        this.userSearchDialogAssign = new UserSearchDialog(this.taskAssgin, {
            "singleEnter": true,
            "isIn": true,
            "forceSubmit": function (sendUuid) {
                if (that.assignUuid === sendUuid) {
                    that.editSubmit(null, that);
                    return;
                }
                let nowTime = new Date();
                that.taskPopupAssignIcon.style.backgroundImage = "url(\"/assets/accountIcon/" + sendUuid + ".png?t=" + nowTime.getTime() + "\")";
                that.assignUuid = sendUuid;
            }
        });

        this.error.addEventListener("click", function(){
            that.hideError();
        }, true);

        this.success.addEventListener("click", function(){
            that.hideSuccess();
        }, true);

        this.taskForm.deadline.addEventListener("change", function () {
            TaskEdit.taskDeadlineChange(that.taskPopup.dataset.taskCreatedDate, this.value);
        });

        this.taskPopup.addEventListener("submit", function (e) {
            that.editSubmit(e);
        }, true);

        this.taskPopupEditButton.addEventListener("click", function (e) {
            TaskEdit.editClick(e);
        }, true);

        this.taskPopupEditCancelButton.addEventListener("click", function (e){
            TaskEdit.editClick(e);
        }, true);
    }

    static editClick() {
        TaskDetail.toggleEditable();
    }

    static taskDeadlineChange(createDate, deadline) {
        let progress = TaskDetail.deadLineProgress(createDate, deadline);

        let taskPopupProgressCurrent = document.querySelector("#taskPopupProgressCurrent");
        if (progress >= 100) {
            progress = 100;
            taskPopupProgressCurrent.classList.add("over");
        } else taskPopupProgressCurrent.classList.remove("over");

        taskPopupProgressCurrent.style.width = progress + "%";
    }

    showError(message){
        this.error.style.display = "block";
        this.error.innerText = message;
    }

    hideError(){
        this.error.style.display = "none";
    }

    showSuccess(){
        this.success.style.display = "block";
    }

    hideSuccess(){
        this.success.style.display = "none";
    }

    getTaskObjectFromDetailForm() {
        let task = TaskApi.NewTaskObject();
        task.name = this.taskForm.name.value;
        task.assign = this.assignUuid;
        task.deadline = this.taskForm.deadline.value;
        task.description = this.taskForm.description.value;

        return task;
    }

    editSubmit(e) {
        if (e !== undefined && e != null)
            e.preventDefault();
        let loadView = new LoadingView();
        loadView.isDisporse = true;
        loadView.show();

        this.hideError();
        this.hideSuccess();
        let createDate = this.taskPopup.dataset.taskCreatedDate;
        let taskId = this.taskPopup.dataset.taskId;
        let updateTask = this.getTaskObjectFromDetailForm();

        let that = this;
        TaskApi.Update(createDate, updateTask).then(function(json){
            if(!json.success) {
                that.showError(json.message);
            }else{
                that.showSuccess();
                TaskDetail.setEditable(false);
            }
            TaskDetail.load(taskId).then(function(task){
                if(task == null) return;

                ProjectUtils.taskboadOnTaskUpdate(task.createDate, task);
                loadView.hide();
            });
        });
    }
}

new TaskEdit();