var taskAddForm = document.querySelector("#taskboardAddForm"),
    errorElem = document.querySelector("#error"),
    taskAddShow = document.querySelector("#taskboardAdd"),
    closeFormElem = document.querySelector("#taskboardAddClose");

taskAddForm.addEventListener("submit", postTaskAdd, true);
errorElem.addEventListener("click", clickError, true);
taskAddShow.addEventListener("click", taskAddShowClick, true);
closeFormElem.addEventListener("click", clickCloseForm, true);

var nowTime = new Date(Date.now());
taskAddForm.querySelector("input[name=deadline]").value = dateFormat(nowTime);


function postTaskAdd(e) {
    e.preventDefault();

    let formData = new FormData(taskAddForm);
    let projectAddJson = {
        "name": formData.get("name"),
        "deadline": formData.get("deadline"),
        "description": formData.get("description")
    };

    fetch("", {
        method: 'POST',
        body: JSON.stringify(projectAddJson),
        credentials: "same-origin"
    }).then(function (response) {
        return response.json();
    }).then(function (json) {
        if (!json.success) {
            errorElem.style.display = "block";
            errorElem.innerText = json.message;
        } else {
            errorElem.style.display = "none";
            TaskApi.GetTaskFromCreateDate(json.createDate).then(function (json) {
                alert();
            });
        }
    });
}

function clickError() {
    errorElem.style.display = "none";
}

function taskAddShowClick() {
    taskAddForm.style.right = "0";
}

function clickCloseForm() {
    taskAddForm.style.right = "-" + taskAddForm.clientWidth + "px";
}

function dateFormat(date) {
    var y = date.getFullYear();
    var m = date.getMonth() + 1;
    var d = date.getDate();

    if (m < 10) {
        m = '0' + m;
    }
    if (d < 10) {
        d = '0' + d;
    }

    return y + '-' + m + '-' + d;
}