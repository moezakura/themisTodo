import Vue from 'vue';
import BackView from "./backView";
import TaskApi from "./taskApi";
import TaskBoard from "./taskBoard"
import ProjectUtils from "./projectUtils";

if (document.querySelector("#taskPopup")) {
    let backViewLayer = new BackView();

    let taskDeleteConfirmPopup = new Vue({
        delimiters: ['${', '}'],
        el: "#taskDeleteConfirmPopup",
        data: {
            taskDeleteConfirmPopupFlag: false,
            taskName: "",
            createDate: "",
        },
        methods: {
            clickCancel() {
                this.taskDeleteConfirmPopupFlag = false;
                backViewLayer.hide();
            },
            clickDelete() {
                let that = this,
                    createDate = this.createDate;
                TaskApi.Delete(createDate).then(function () {
                    that.taskDeleteConfirmPopupFlag = false;
                    backViewLayer.hide();
                    return createDate;
                }).then(function(createDate){
                    ProjectUtils.taskboadOnTaskDelete(createDate);
                });
            }
        }
    });

    let taskDeletePopup = new Vue({
        delimiters: ['${', '}'],
        el: '#taskDeletePopup',
        data: {
            taskDeletePopupFlag: false,
            taskName: "",
            createDate: "",
        },
        methods: {
            clickHide() {
                let task = TaskApi.NewTaskObject();
                let taskStatuses = TaskApi.GetTaskStatuses();
                task.status = taskStatuses.STATUS_HIDE;
                let that = this;

                TaskApi.Update(this.createDate, task).then(function (json) {
                    TaskBoard.loadTask(that.createDate).then(function () {
                        that.clickClose();
                    });
                });
            },
            clickDelete() {
                this.clickClose();
                taskDeleteConfirmPopup.taskDeleteConfirmPopupFlag = true;
                taskDeleteConfirmPopup.createDate = this.createDate;
                backViewLayer.show();
            },
            clickClose() {
                this.taskDeletePopupFlag = false;
                backViewLayer.hide();
            }
        }
    });

    new Vue({
        delimiters: ['${', '}'],
        el: '#taskPopupTrashButton',
        methods: {
            clickTaskPopupTrashButton() {
                backViewLayer.show();

                taskDeletePopup.taskDeletePopupFlag = true;
                taskDeletePopup.taskName = document.querySelector("#taskPopupTitle").value; // TODO: Rewrite Vue
                taskDeletePopup.createDate = document.querySelector("#taskPopup").dataset.taskCreatedDate; // TODO: Rewrite Vue
                document.querySelector("#taskPopupCloseButton").click();
            }
        },
    });

    backViewLayer.addHideEvent(function () {
        taskDeletePopup.taskDeletePopupFlag = false;
        taskDeleteConfirmPopup.taskDeleteConfirmPopupFlag = false;
    });
}