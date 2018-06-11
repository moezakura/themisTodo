var taskConfigShow = document.querySelector("#taskboardConfig"),
    taskBackView = new BackView(),
    projectConfigPopup = document.querySelector("#projectConfigPopup"),
    postTaskConfigForm = document.querySelector("#projectConfigForm"),
    taskboardTitleElem = document.querySelector("#taskboardTitle>span"),
    projectConfigPopupErrorElem = projectConfigPopup.querySelector(".error");

postTaskConfigForm.addEventListener("submit", postTaskConfig, true);
taskConfigShow.addEventListener("click", taskConfigShowClick, true);
projectConfigPopupErrorElem.addEventListener("click", clickError, true);
taskBackView.addWithHideElem(projectConfigPopup);
taskBackView.addHideEvent(function(){
    userSearchDialog.hide();
});

var userSelectInput = document.querySelector("#userSelect");

var userSearchDialog = new UserSearchDialog(userSelectInput, {
    "forceSubmit": function(sendUuid){
        let projectAddMemberJson = {
            "uuid": sendUuid
        };

        fetch("/project/addUser/" + projectId, {
            method: "POST",
            body: JSON.stringify(projectAddMemberJson),
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        }).then(function (json) {
            if (!json.success) {
                projectConfigPopupErrorElem.style.display = "block";
                projectConfigPopupErrorElem.innerText = json.message;
            } else {
                projectConfigPopupErrorElem.style.display = "none";
                userSelectInput.value = "";
                memberList.push(json.addedAccount);
                memberUpdate();
            }
        });
    }
});

document.querySelector("body").addEventListener("keydown", function (e) {
    if (e.keyCode === 27 && projectConfigPopup.style.display === "block")
        taskBackView.hide();
}, true);

function taskConfigShowClick(e) {
    e.preventDefault();
    projectConfigPopup.style.display = "block";
    taskBackView.show();
}

function postTaskConfig(e) {
    e.preventDefault();

    let formData = new FormData(postTaskConfigForm);
    let projectNewName = formData.get("name");
    let postTaskConfigJson = {
        "name": projectNewName,
        "description": formData.get("description"),
    };

    fetch("/project/update/" + projectId, {
        method: 'POST',
        body: JSON.stringify(postTaskConfigJson),
        credentials: "same-origin"
    }).then(function (response) {
        return response.json();
    }).then(function (json) {
        if (!json.success) {
            projectConfigPopupErrorElem.style.display = "block";
            projectConfigPopupErrorElem.innerText = json.message;
        } else {
            projectConfigPopupErrorElem.style.display = "none";
            taskboardTitleElem.innerText = projectNewName;
        }
    });
}




