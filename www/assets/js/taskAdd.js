import ProjectUtils from "./projectUtils"
import UserSearchDialog from "./userSearchDialog"
import TaskApi from "./taskApi"

class TaskAdd {
    constructor() {
        this.taskAddForm = document.querySelector("#taskboardAddForm");
        if (this.taskAddForm === undefined || this.taskAddForm == null)
            return;
        this.taskAddFormErrorElem = this.taskAddForm.querySelector(".error");
        this.taskAddShow = document.querySelector("#taskboardAdd");
        this.closeFormElem = document.querySelector("#taskboardAddClose");
        this.assignInput = document.querySelector("#assign");

        let that = this;

        this.taskAddForm.addEventListener("submit", function (e) {
            that.postTaskAdd(e, that)
        }, true);
        this.taskAddFormErrorElem.addEventListener("click", ProjectUtils.clickObjectHide, true);
        this.taskAddShow.addEventListener("click", function (e) {
            that.taskAddShowClick(e, that)
        }, true);
        this.closeFormElem.addEventListener("click", function () {
            that.clickCloseForm(that);
        }, true);

        let nowTime = new Date(Date.now());
        this.assignUuid = 0;
        this.taskAddForm.querySelector("input[name=deadline]").value = ProjectUtils.dateFormat(nowTime);

        document.querySelector("body").addEventListener("keydown", function (e) {
            if (e.keyCode === 27 && that.taskAddForm.classList.contains("shown"))
                that.clickCloseForm(that);
        }, true);

        this.createTaskBoard();

        this.userSearchDialogAssign = new UserSearchDialog(this.assignInput, {
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
    }


    createTaskBoard() {
        let that = this;
        taskList.forEach(function (task) {
            let taskElem = ProjectUtils.createTaskItem(task.createDate, task.name, task.taskId, task.assignName,
                task.assign, task.deadlineMD, task.limitDate);

            if(document.taskBoard.taskBoardLists[task.status] !== undefined && document.taskBoard.taskBoardLists[task.status] != null)
                document.taskBoard.taskBoardLists[task.status].appendChild(taskElem);
        });
    }

    postTaskAdd(e, that) {
        try {
            if (e !== undefined && e != null) e.preventDefault();
        } catch (ex) {
        }

        let formData = new FormData(that.taskAddForm);
        let projectAddJson = {
            "name": formData.get("name"),
            "deadline": formData.get("deadline"),
            "description": formData.get("description"),
            "assign": that.assignUuid,
            "projectId": projectId
        };

        TaskApi.Create(projectAddJson).then(function (json) {
            if (!json.success) {
                that.taskAddFormErrorElem.style.display = "block";
                that.taskAddFormErrorElem.innerText = json.message;
            } else {
                that.taskAddFormErrorElem.style.display = "none";
                that.addFormClear(that);
                TaskApi.GetTaskFromCreateDate(json.createDate).then(function (json) {
                    if (!json.success) {
                        console.error("API ERROR");
                        return
                    }
                    let task = json.task;
                    let taskElem = ProjectUtils.createTaskItem(task.createDate, task.name, task.taskId, task.assignName,
                        task.assign, task.deadlineMD, task.limitDate);

                    document.taskBoard.taskBoardLists[0].appendChild(taskElem);
                });
            }
        });
    }

    taskAddShowClick(e, that) {
        e.preventDefault();
        that.taskAddForm.classList.add("shown");
    }

    clickCloseForm(that) {
        that.userSearchDialogAssign.hide();
        that.taskAddForm.classList.remove("shown");
    }

    addFormClear(that) {
        that.taskAddForm.reset();
        let nowTime = new Date(Date.now());
        that.taskAddForm.querySelector("input[name=deadline]").value = ProjectUtils.dateFormat(nowTime);
        that.userSearchDialogAssign.hide();
    }
}

new TaskAdd();