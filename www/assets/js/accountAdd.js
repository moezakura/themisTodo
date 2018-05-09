var accountAddForm = document.querySelector("#accountAdd"),
    errorElem = document.querySelector("#error");

accountAddForm.addEventListener("submit", postLogin, true);
errorElem.addEventListener("click", clickError, true);


function postLogin(e) {
    e.preventDefault();

    let formData = new FormData(accountAddForm);
    let accountAddJson = {
        "name": formData.get("name")
    };

    fetch("", {
        method: 'POST',
        body: JSON.stringify(accountAddJson),
        credentials: "same-origin"
    }).then(function (response) {
        return response.json();
    }).then(function (json) {
        if (!json.success) {
            errorElem.style.display = "block";
            errorElem.innerText = json.message;
        } else {
            errorElem.style.display = "none";
        }
    });
}

function clickError() {
    errorElem.style.display = "none";
}