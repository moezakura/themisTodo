import BackView from "./backView"
import LoadingView from "./loadingView";
import TaskApi from "./taskApi";
import ProjectUtils from "./projectUtils";

export default class TaskDetail {
    static show() {
        let taskPopup = document.querySelector("#taskPopup");
        taskPopup.style.display = "block";

        let backView = new BackView();
        backView.isDisporse = true;
        backView.show();

        backView.addWithHideElem(taskPopup);
        window.addEventListener("hashchange", function() {
                if (document.location.hash === "") {
                    backView.hide();
                }
            }, false
        );
    }

    static load(createDate) {
        TaskApi.GetTaskFromCreateDate(createDate).then(function (json) {
            if (!json.success) {
                console.error("API ERROR");
                return
            }

            TaskDetail.replaceUrlHash(json.task.taskId);
        });
    }

    static loadFromTaskId(taskId, projectId) {
        let loadView = new LoadingView();
        loadView.isDisporse = true;
        loadView.show();

        TaskApi.GetSearch(taskId, projectId).then(function (json) {
            if (!json.success) {
                console.error("API ERROR");
                loadView.hide();
                return
            }

            TaskDetail.set(json.task);

            loadView.hide();
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
            taskPopupInputs = taskPopup.querySelectorAll("input, textarea"),
            nowTime = new Date(),
            statusText = ["Todo", "Doing", "PullRequest", "Done"][taskObject.status];

        taskPopupInputs.forEach(function(value){
           value.readOnly = true;
        });

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

        TaskDetail.replaceUrlHash(taskObject.taskId);
    }

    static replaceUrlHash(taskId) {
        window.location.hash = taskId;
    }

    static refreshUrlHash() {
        window.location.hash = "";
    }
}