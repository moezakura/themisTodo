import TaskDetail from "./taskDetail"

export default class ProjectUtils {
    static createUserListLine(uuid, name, displayName) {
        let parentLi = document.createElement("li");
        parentLi.dataset.uuid = uuid;
        parentLi.dataset.name = name;
        parentLi.dataset.displayName = displayName;
        parentLi.id = "searchResult" + uuid;

        let iconDiv = document.createElement("div");
        let nowTime = new Date();
        iconDiv.classList.add("icon");
        iconDiv.style.backgroundImage = "url(\"/assets/accountIcon/" + uuid + ".png?t=" + nowTime.getTime() + "\")";
        parentLi.appendChild(iconDiv);

        let nameLineDiv = document.createElement("div");
        nameLineDiv.classList.add("name");

        {
            let nameDiv = document.createElement("div");
            nameDiv.classList.add("nameId");
            nameDiv.innerText = name;

            let displayNameDiv = document.createElement("div");
            displayNameDiv.classList.add("displayName");
            displayNameDiv.innerText = displayName;

            nameLineDiv.appendChild(nameDiv);
            nameLineDiv.appendChild(displayNameDiv);
        }

        parentLi.appendChild(nameLineDiv);

        return parentLi;
    }

    static createTaskItem(createTime, name, taskId, assignName, assignIcon, deadLine, limitDate, status) {
        let parentLi = document.createElement("li");
        parentLi.id = "singleTask-" + createTime;
        parentLi.dataset.id = createTime;
        parentLi.dataset.taskId = taskId;
        parentLi.addEventListener("click", function(){
            TaskDetail.replaceUrlHash(taskId);
        });

        if (status === undefined || status == null || status !== 3) {
            if (limitDate <= 0) parentLi.classList.add("over");
            else if (limitDate <= 1) parentLi.classList.add("limit1");
            else if (limitDate <= 2) parentLi.classList.add("limit2");
            else if (limitDate <= 3) parentLi.classList.add("limit3");
        }


        let nowTime = new Date();

        {
            let taskTitle = document.createElement("div");
            taskTitle.innerText = name;
            taskTitle.classList.add("taskTitle");

            let taskIdDiv = document.createElement("div");
            taskIdDiv.innerText = "#" + taskId;
            taskIdDiv.classList.add("taskId");

            let taskAssign = document.createElement("div");
            taskAssign.classList.add("taskAssign");
            {
                let taskAssignIcon = document.createElement("div");
                taskAssignIcon.style.backgroundImage = "url(\"/assets/accountIcon/" + assignIcon + ".png?t=" + nowTime.getTime() + "\")";
                taskAssignIcon.classList.add("taskAssignIcon");

                let taskAssignName = document.createElement("div");
                taskAssignName.innerText = assignName;
                taskAssignName.classList.add("taskAssignName");

                taskAssign.appendChild(taskAssignIcon);
                taskAssign.appendChild(taskAssignName);
            }

            let taskLimit = document.createElement("div");
            taskLimit.classList.add("taskLimit");
            {
                let taskLimitIcon = document.createElement("i");
                taskLimitIcon.classList.add("fas");
                taskLimitIcon.classList.add("fa-calendar-alt");

                let taskLimitDead = document.createElement("span");
                taskLimitDead.classList.add("deadlineDate");
                taskLimitDead.innerText = deadLine;

                let taskLimitDeadRel = document.createElement("span");
                taskLimitDeadRel.classList.add("deadlineDate");

                if (status === undefined || status == null || status !== 3)
                    taskLimitDeadRel.innerText = "あと" + limitDate + "日";
                else
                    taskLimitDeadRel.innerText = "Completed!";

                taskLimit.appendChild(taskLimitIcon);
                taskLimit.appendChild(taskLimitDead);
                taskLimit.appendChild(taskLimitDeadRel);
            }

            parentLi.appendChild(taskTitle);
            parentLi.appendChild(taskIdDiv);
            parentLi.appendChild(taskAssign);
            parentLi.appendChild(taskLimit);
        }

        return parentLi
    }

    /***
     *
     * @param createDate {Number}
     * @param taskObject {Array}
     */
    static taskboadOnTaskUpdate(createDate, taskObject) {
        let taskLi = document.querySelector("#singleTask-" + createDate);
        if (taskLi === undefined || taskLi == null)
            return;

        let name = taskLi.querySelector(".taskTitle"),
            assignName = taskLi.querySelector(".taskAssignName"),
            assignIcon = taskLi.querySelector(".taskAssignIcon"),
            deadlines = taskLi.querySelectorAll(".deadlineDate"),
            deadlineDate = deadlines[0],
            deadlineLimit = deadlines[1],
            nowTime = new Date();

        // task
        {
            taskLi.classList.remove("over", "limit1", "limit2", "limit3");
            if (taskObject.status === undefined || taskObject.status == null || taskObject.status !== 3) {
                if (taskObject.limitDate <= 0) taskLi.classList.add("over");
                else if (taskObject.limitDate <= 1) taskLi.classList.add("limit1");
                else if (taskObject.limitDate <= 2) taskLi.classList.add("limit2");
                else if (taskObject.limitDate <= 3) taskLi.classList.add("limit3");
            }
        }

        // name
        {
            name.innerText = taskObject.name;
        }

        // assign
        {
            assignName.innerText = taskObject.assignName;
            assignIcon.style.backgroundImage = "url(\"/assets/accountIcon/" + taskObject.assign + ".png?t=" + nowTime.getTime() + "\")";
        }

        // deadline
        {
            deadlineDate.innerText = taskObject.deadlineMD;

            if (taskObject.status === undefined || taskObject.status == null || taskObject.status !== 3)
                deadlineLimit.innerText = "あと" + taskObject.limitDate + "日";
            else
                deadlineLimit.innerText = "Completed!";
        }
    }

    static clickObjectHide(e) {
        e.target.style.display = "none";
    }

    static dateFormat(date) {
        let y = date.getFullYear();
        let m = date.getMonth() + 1;
        let d = date.getDate();

        if (m < 10) {
            m = '0' + m;
        }
        if (d < 10) {
            d = '0' + d;
        }

        return y + '-' + m + '-' + d;
    }
}