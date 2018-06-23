import TaskDetail from "./taskDetail"
import TaskApi from "./taskApi";
import UserSearchDialog from "./userSearchDialog";

class TaskEdit {
    constructor() {
        this.taskPopup = document.querySelector("#taskPopup");
        if (this.taskPopup === undefined || this.taskPopup == null)
            return;
        this.taskPopupEditButton = this.taskPopup.querySelector("#taskPopupEditButton");
        this.taskAssgin = this.taskPopup.querySelector("#taskPopupAssign");

        let that = this;

        this.userSearchDialogAssign = new UserSearchDialog(this.taskAssgin, {
            "singleEnter": true,
            "isIn": true,
            "forceSubmit": function (sendUuid) {
                if (that.assignUuid === sendUuid) {
                    that.postTaskAdd(null, that);
                    return;
                }
                that.assignUuid = sendUuid;
            }
        });

        this.taskPopup.addEventListener("submit", function(e){
            that.editSubmit(e);
        }, true);
        this.taskPopupEditButton.addEventListener("click", function(e){
            TaskEdit.editClick(e);
        }, true);
    }

    static editClick() {
        TaskDetail.editable(true);
    }

    editSubmit(e) {
        e.preventDefault();
        let createDate = this.taskPopup.dataset.taskId;
        let task = TaskApi.NewTaskObject();

        TaskApi.Update();
    }
}

new TaskEdit();