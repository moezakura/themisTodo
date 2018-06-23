import BackView from "./backView"
import LoadingView from "./loadingView";
import TaskApi from "./taskApi";

export default class TaskDetail {
    static show() {
        let taskPopup = document.querySelector("#taskPopup");
        taskPopup.style.display = "block";

        let backView = new BackView();
        backView.isDisporse = true;
        backView.show();

        backView.addWithHideElem(taskPopup);

        let taskPopupCloseButton = document.querySelector("#taskPopupCloseButton");
        let closeEventFunc = function () {
            backView.hide();
            TaskDetail.setEmpty(true);
            taskPopupCloseButton.removeEventListener("click", closeEventFunc, true);
        };
        backView.addHideEvent(function () {
            TaskDetail.setEmpty(true);
            taskPopupCloseButton.removeEventListener("click", closeEventFunc, true);
        });

        taskPopupCloseButton.addEventListener("click", closeEventFunc, true);
    }

    static load(taskId) {
        TaskDetail.setEmpty(true);
        let taskPopup = document.querySelector("#taskPopup");
        taskPopup.dataset.taskId = taskId;

        let loadView = new LoadingView();
        loadView.isDisporse = true;
        loadView.show();

        TaskApi.GetTaskFromCreateDate(taskId).then(function (json) {
            if (!json.success) {
                console.error("API ERROR");
                loadView.hide();
                return
            }

            TaskDetail.set(json.task);

            loadView.hide();
        });
    }

    static loadAndShow(taskId) {
        this.show();
        this.load(taskId);
    }

    static setEmpty(readOnly) {
        if (readOnly === undefined || readOnly == null)
            readOnly = false;

        let taskPopup = document.querySelector("#taskPopup"),
            taskDetailTitle = taskPopup.querySelector("h2"),
            taskPopupTaskId = taskPopup.querySelector("#taskPopupTaskId"),
            taskPopupTitle = taskPopup.querySelector("#taskPopupTitle"),
            taskPopupProgressText = taskPopup.querySelector("#taskPopupProgressText"),
            taskPopupDescription = taskPopup.querySelector("#taskPopupDescription"),
            taskPopupProgressTextSpans = taskPopupProgressText.querySelectorAll("span"),
            taskPopupAssignIcon = taskPopup.querySelector("#taskPopupAssignIcon"),
            taskPopupAssign = taskPopup.querySelector("#taskPopupAssign"),
            taskPopupCreatorIcon = taskPopup.querySelector("#taskPopupCreatorIcon"),
            taskPopupCreator = taskPopup.querySelector("#taskPopupCreator"),
            taskPopupProgressCurrent = document.querySelector("#taskPopupProgressCurrent");

        this.editable(!readOnly);
        taskPopup.dataset.taskId = "";
        taskDetailTitle.innerText = "";
        taskPopupTaskId.innerText = "#";
        taskPopupTitle.value = "";
        taskPopupProgressTextSpans[0].innerText = "";
        taskPopupProgressTextSpans[1].innerText = "";
        taskPopupProgressCurrent.style.width = "0%";
        taskPopupAssignIcon.style.backgroundImage = "";
        taskPopupAssign.value = "";
        taskPopupCreatorIcon.style.backgroundImage = "";
        taskPopupCreator.innerText = "";
        taskPopupDescription.value = "";
    }

    static set(taskObject) {
        let taskPopup = document.querySelector("#taskPopup"),
            taskDetailTitle = taskPopup.querySelector("h2"),
            taskPopupTaskId = taskPopup.querySelector("#taskPopupTaskId"),
            taskPopupTitle = taskPopup.querySelector("#taskPopupTitle"),
            taskPopupProgressText = taskPopup.querySelector("#taskPopupProgressText"),
            taskPopupDescription = taskPopup.querySelector("#taskPopupDescription"),
            taskPopupProgressTextSpans = taskPopupProgressText.querySelectorAll("span"),
            taskPopupAssignIcon = taskPopup.querySelector("#taskPopupAssignIcon"),
            taskPopupAssign = taskPopup.querySelector("#taskPopupAssign"),
            taskPopupCreatorIcon = taskPopup.querySelector("#taskPopupCreatorIcon"),
            taskPopupCreator = taskPopup.querySelector("#taskPopupCreator"),
            nowTime = new Date(),
            statusText = ["Todo", "Doing", "PullRequest", "Done"][taskObject.status];

        this.editable(false);
        taskDetailTitle.innerText = statusText + " Task Detail";

        taskPopupTaskId.innerText = "#" + taskObject.taskId;
        taskPopupTitle.value = taskObject.name;

        {
            taskPopupProgressTextSpans[0].innerText = taskObject.deadline;
            taskPopupProgressTextSpans[1].innerText = "(あと" + taskObject.limitDate + "日)";

            let createDate = Math.round(taskObject.createDate / 1000000),
                limitDate = (new Date(taskObject.deadline)).getTime(),
                allDiff = limitDate - createDate,
                limit = limitDate - (new Date()).getTime(),
                progress = 100 - limit / allDiff * 100;

            let taskPopupProgressCurrent = document.querySelector("#taskPopupProgressCurrent");
            if (progress >= 100) {
                progress = 100;
                taskPopupProgressCurrent.classList.add("over");
            } else taskPopupProgressCurrent.classList.remove("over");
            taskPopupProgressCurrent.style.width = progress + "%";
        }

        {
            taskPopupAssignIcon.style.backgroundImage = "url(\"/assets/accountIcon/" + taskObject.assign + ".png?t=" + nowTime.getTime() + "\")";
            taskPopupAssign.value = taskObject.assignName;

            taskPopupCreatorIcon.style.backgroundImage = "url(\"/assets/accountIcon/" + taskObject.creator + ".png?t=" + nowTime.getTime() + "\")";
            taskPopupCreator.innerText = taskObject.creatorName;
        }

        taskPopupDescription.value = taskObject.description;
    }

    static editable(isEdit) {
        let taskPopupInputs = taskPopup.querySelectorAll("input, textarea");
        taskPopupInputs.forEach(function (value) {
            value.readOnly = !isEdit;
        });
    }
}