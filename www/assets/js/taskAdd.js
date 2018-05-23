var taskAddForm = document.querySelector("#taskboardAddForm"),
    taskAddFormErrorElem = taskAddForm.querySelector(".error"),
    taskAddShow = document.querySelector("#taskboardAdd"),
    closeFormElem = document.querySelector("#taskboardAddClose");

taskAddForm.addEventListener("submit", postTaskAdd, true);
taskAddFormErrorElem.addEventListener("click", clickError, true);
taskAddShow.addEventListener("click", taskAddShowClick, true);
closeFormElem.addEventListener("click", clickCloseForm, true);

var nowTime = new Date(Date.now());
taskAddForm.querySelector("input[name=deadline]").value = dateFormat(nowTime);

document.querySelector("body").addEventListener("keydown", function (e) {
    if (e.keyCode === 27 && taskAddForm.classList.contains("shown"))
        clickCloseForm();
}, true);

createTaskBoard();
function createTaskBoard(){
    taskList.forEach(function(task){
        let taskElem = ProjectUtils.createTaskItem(task.createDate, task.name, task.taskId, task.assignName,
            task.assign, task.deadlineMD, task.limitDate);

        taskBoardLists[task.status].appendChild(taskElem);
    });
}

function postTaskAdd(e) {
    e.preventDefault();

    let formData = new FormData(taskAddForm);
    let projectAddJson = {
        "name": formData.get("name"),
        "deadline": formData.get("deadline"),
        "description": formData.get("description"),
        "assign": Number(formData.get("assign")),
        "projectId": projectId
    };

    TaskApi.Create(projectAddJson).then(function (json) {
        if (!json.success) {
            taskAddFormErrorElem.style.display = "block";
            taskAddFormErrorElem.innerText = json.message;
        } else {
            taskAddFormErrorElem.style.display = "none";
            addFormClear();
            TaskApi.GetTaskFromCreateDate(json.createDate).then(function (json) {
                if(!json.success) {
                    console.error("API ERROR");
                    return
                }
                let task = json.task;
                let taskElem = ProjectUtils.createTaskItem(task.createDate, task.name, task.taskId, task.assignName,
                    task.assign, task.deadlineMD, task.limitDate);

                taskBoardLists[0].appendChild(taskElem);
            });
        }
    });
}

function clickError(e) {
    e.target.style.display = "none";
}

function taskAddShowClick(e) {
    e.preventDefault();
    taskAddForm.classList.add("shown");
}

function clickCloseForm() {
    taskAddForm.classList.remove("shown");
}

function addFormClear() {
    taskAddForm.reset();
    let nowTime = new Date(Date.now());
    taskAddForm.querySelector("input[name=deadline]").value = dateFormat(nowTime);
}

function dateFormat(date) {
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