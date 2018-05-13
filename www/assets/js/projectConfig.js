var taskConfigShow = document.querySelector("#taskboardConfig"),
    backViewLayerElem = document.querySelector("#backViewLayer"),
    projectConfigPopup = document.querySelector("#projectConfigPopup"),
    postTaskConfigForm = document.querySelector("#projectConfigForm"),
    taskboardTitleElem = document.querySelector("#taskboardTitle>span"),
    projectConfigPopupErrorElem = projectConfigPopup.querySelector(".error");

postTaskConfigForm.addEventListener("submit", postTaskConfig, true);
taskConfigShow.addEventListener("click", taskConfigShowClick, true);
backViewLayerElem.addEventListener("click", backViewLayerElemClick, true);
projectConfigPopupErrorElem.addEventListener("click", clickError, true);


var userSelectInput = document.querySelector("#userSelect"),
    userSelectUserList = document.querySelector("#usersSearchedList"),
    projectMemberAddForm = projectConfigPopup.querySelector("#projectMemberAddForm");
var selectUserIndex = -1, selectMode = false, oldSearchInput = "", searchResult = [];

userSelectInput.addEventListener("keyup", userSelectInputKeyUp, true);
userSelectInput.addEventListener("keydown", userSelectInputKeyDown, true);
projectMemberAddForm.addEventListener("submit", projectMemberAddSubmit, true);

document.querySelector("body").addEventListener("keydown", function (e) {
    if (e.keyCode === 27 && projectConfigPopup.style.display === "block")
        backViewLayerElemClick();
}, true);

function taskConfigShowClick(e) {
    e.preventDefault();
    projectConfigPopup.style.display = "block";
    backViewLayerElem.style.display = "block";
}

function backViewLayerElemClick() {
    projectConfigPopup.style.display = "none";
    backViewLayerElem.style.display = "none";
    userSelectUserList.style.display = "none";
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


function userSelectInputKeyUp(e) {
    let inputText = userSelectInput.value;
    if (inputText.length < 1) {
        userSelectUserList.style.display = "none";
        selectUserIndex = -1;
        return;
    }

    if (e.keyCode === 38 || e.keyCode === 40 || e.keyCode === 13) return;

    if (inputText !== oldSearchInput) {
        selectUserIndex = -1;
    }

    let queryObject = {
        "name": inputText,
        "displayName": inputText,
        "project": projectId,
        "max": 20
    };
    let queryString = "?";
    for (key in queryObject) queryString += key + "=" + encodeURIComponent(queryObject[key]) + "&";
    queryString = queryString.slice(0, -1);

    fetch("/account/search" + queryString, {
        method: 'GET'
    }).then(function (response) {
        return response.json();
    }).then(function (json) {
        searchResult = json;
        let targetPos = userSelectInput.getBoundingClientRect();
        userSelectUserList.style.display = "block";
        userSelectUserList.style.left = (targetPos.left + window.pageXOffset) + "px";
        userSelectUserList.style.top = (targetPos.top + targetPos.height + window.pageYOffset) + "px";
        userSelectUserList.style.width = targetPos.width + "px";
        userSelectUserList.scrollTop = 0;
        selectMode = false;

        userSelectUserList.innerHTML = "";
        json.forEach(function (value) {
            let elem = ProjectUtils.createUserListLine(value.uuid, value.name, value.displayName);
            userSelectUserList.appendChild(elem);
        });
    });
}

function userSelectInputKeyDown(e) {
    if (e.keyCode === 38) {
        //up key
        selectUserIndex--;
        if (selectUserIndex < 0) selectUserIndex = searchResult.length - 1;
    } else if (e.keyCode === 40) {
        //down key
        selectUserIndex++;
        if (selectUserIndex > searchResult.length - 1) selectUserIndex = 0;
    } else return;

    if (selectUserIndex < 0 || selectUserIndex < 0 || selectUserIndex > searchResult.length - 1) return;
    selectMode = true;

    let all = userSelectUserList.querySelectorAll("li");
    all.forEach(function (value) {
        value.classList.remove("select");
    });

    let targetId = "#searchResult" + searchResult[selectUserIndex].uuid;
    let target = userSelectUserList.querySelector(targetId);
    let targetPos = target.getBoundingClientRect();
    let userSelectUserListPos = userSelectUserList.getBoundingClientRect();
    userSelectUserList.scrollTop = targetPos.y + userSelectUserList.scrollTop - userSelectUserListPos.height + targetPos.height - userSelectUserListPos.y;
    target.classList.add("select");
}

function projectMemberAddSubmit(e) {
    e.preventDefault();

    if (selectMode) {
        if (!(selectUserIndex < 0 || selectUserIndex < 0 || selectUserIndex > searchResult.length - 1))
            userSelectInput.value = searchResult[selectUserIndex].name;
        userSelectUserList.style.display = "none";
        selectMode = false;
    } else {
        let sendUuid = -1;
        if (selectUserIndex < 0 || selectUserIndex < 0 || selectUserIndex > searchResult.length - 1) {

        } else sendUuid = searchResult[selectUserIndex].uuid;

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
}
