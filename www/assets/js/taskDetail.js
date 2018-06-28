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

        let isHashEmpty = function() {
            if (document.location.hash === "") {
                backView.hide();
            }
        };
        backView.addHideEvent(function(){
            window.removeEventListener("hashchange", isHashEmpty);
        });

        window.addEventListener("hashchange", isHashEmpty);
    }

    static loadFromTaskId(taskId, projectId) {
        let backView = new BackView();
        backView.isDisporse = true;

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

        let errorMessage = taskPopup.querySelector(".error"),
            successMessage = taskPopup.querySelector(".success");
        errorMessage.style.display = "none";
        successMessage.style.display = "none";

        TaskDetail.load(taskId);
    }
    
    static load(taskId) {
        TaskDetail.setEmpty(true);
        let taskPopup = document.querySelector("#taskPopup");
        taskPopup.dataset.taskId = taskId;

        let loadView = new LoadingView();
        loadView.isDisporse = true;
        loadView.show();

        return TaskApi.GetSearch(taskId, projectId).then(function (json) {
            if (!json.success) {
                console.error("API ERROR");
                loadView.hide();
                return null;
            }

            TaskDetail.set(json.task);
            TaskDetail.replaceUrlHash(json.task.taskId);

            loadView.hide();
            return json.task;
        });
    }

    static loadAndShow(createDate) {
        // this.show();
        this.load(createDate);
    }

    static loadAndShowFromTaskId(taskId, projectId) {
        TaskDetail.show();
        TaskDetail.loadFromTaskId(taskId, projectId);
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
            taskPopupDeadlineChange = document.querySelector("#taskPopupDeadlineChange"),
            nowTime = new Date(),
            statusText = ["Todo", "Doing", "PullRequest", "Done"][taskObject.status];

        taskPopup.dataset.taskCreatedDate = taskObject.createDate;

        this.editable(false);
        taskDetailTitle.innerText = statusText + " Task Detail";

        taskPopupTaskId.innerText = "#" + taskObject.taskId;
        taskPopupTitle.value = taskObject.name;

        {
            taskPopupProgressTextSpans[0].innerText = taskObject.deadline;
            if (taskObject.status !== 3)
                taskPopupProgressTextSpans[1].innerText = "(あと" + taskObject.limitDate + "日)";
            else
                taskPopupProgressTextSpans[1].innerText = "(Already Completed!)";

            taskPopupDeadlineChange.value = taskObject.deadline;

            let progress = TaskDetail.deadLineProgress(taskObject.createDate, taskObject.deadline);

            let taskPopupProgressCurrent = document.querySelector("#taskPopupProgressCurrent");

            taskPopupProgressCurrent.classList.remove("over", "limit1", "limit2", "limit3");
            taskDetailTitle.classList.remove("over", "limit1", "limit2", "limit3");
            let addClassName = "";
            if (taskObject.limitDate <= 0) {
                addClassName = "over";
                progress = 100;
            }
            else if (taskObject.limitDate <= 1) addClassName = "limit1";
            else if (taskObject.limitDate <= 2) addClassName = "limit2";
            else if (taskObject.limitDate <= 3) addClassName = "limit3";

            if (addClassName.length > 0 && taskObject.status !== 3) {
                taskPopupProgressCurrent.classList.add(addClassName);
                taskDetailTitle.classList.add(addClassName);
            }

            taskPopupProgressCurrent.style.width = (progress > 100 ? 100 : progress) + "%";
        }

        {
            taskPopupAssignIcon.style.backgroundImage = "url(\"/assets/accountIcon/" + taskObject.assign + ".png?t=" + nowTime.getTime() + "\")";
            taskPopupAssign.value = taskObject.assignName;

            taskPopupCreatorIcon.style.backgroundImage = "url(\"/assets/accountIcon/" + taskObject.creator + ".png?t=" + nowTime.getTime() + "\")";
            taskPopupCreator.innerText = taskObject.creatorName;
        }

        taskPopupDescription.value = taskObject.description;

        TaskDetail.replaceUrlHash(taskObject.taskId);
    }

    static replaceUrlHash(taskId) {
        window.location.hash = taskId;
    }

    static refreshUrlHash() {
        window.location.hash = "";
    }

    static deadLineProgress(createDateNanon, deadLine) {
        let createDate = Math.round(createDateNanon / 1000000),
            limitDate = (new Date(deadLine)).getTime(),
            allDiff = Math.abs(limitDate - createDate),
            limit = Math.abs(limitDate - (new Date()).getTime());

        return limit / allDiff * 100;
    }

    static editable(isEdit) {
        let taskPopupInputs = taskPopup.querySelectorAll("input, textarea");
        taskPopupInputs.forEach(function (value) {
            value.readOnly = !isEdit;
        });

        let taskPopupDeadlineChange = document.querySelector("#taskPopupDeadlineChange"),
            taskPopupProgressText = document.querySelector("#taskPopupProgressText");
        taskPopupDeadlineChange.style.display = isEdit ? "block" : "none";
        taskPopupProgressText.style.display = isEdit ? "none" : "block";
    }
}