var loginForm = document.querySelector("#login"),
    errorElem = document.querySelector("#error");

loginForm.addEventListener("submit", postLogin, true);
errorElem.addEventListener("click", clickError, true);


function postLogin(e) {
    e.preventDefault();

    let formData = new FormData(loginForm);
    let loginJson = {
        "id": formData.get("id"),
        "password": formData.get("pw")
    };

    fetch("", {
        method: 'POST',
        body: JSON.stringify(loginJson),
        credentials: "same-origin"
    }).then(function (response) {
        return response.json();
    }).then(function (json) {
        if(!json.success){
            errorElem.style.display = "block";
            errorElem.innerText = json.message;
        }else{
            errorElem.style.display = "none";
            location.href = "home";
        }
    });
}

function clickError() {
    errorElem.style.display = "none";
}