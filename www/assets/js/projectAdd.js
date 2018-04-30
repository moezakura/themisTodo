var projectAddForm = document.querySelector("#projectAdd"),
    errorElem = document.querySelector("#error");

projectAddForm.addEventListener("submit", postLogin, true);
errorElem.addEventListener("click", clickError, true);


function postLogin(e) {
    e.preventDefault();

    let formData = new FormData(projectAddForm);
    let projectAddJson = {
        "name": formData.get("name"),
        "description": formData.get("description")
    };

    fetch("", {
        method: 'POST',
        body: JSON.stringify(projectAddJson)
    }).then(function (response) {
        return response.json();
    }).then(function (json) {
        if (!json.success) {
            errorElem.style.display = "block";
            errorElem.innerText = json.message;
        } else {
            errorElem.style.display = "none";
            location.href = "/project/" + json.id;
        }
    });
}

function clickError() {
    errorElem.style.display = "none";
}