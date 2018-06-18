class Login {
    constructor() {
        this.loginForm = document.querySelector("#login");
        this.errorElem = document.querySelector("#error");

        if (this.loginForm === undefined || this.loginForm == null) return;

        let that = this;
        this.loginForm.addEventListener("submit", function (e) {
            that.postLogin(e, that)
        }, true);

        this.errorElem.addEventListener("click", function (e) {
            that.clickError(e, that);
        }, true);
    }

    postLogin(e, that) {
        e.preventDefault();

        let formData = new FormData(that.loginForm);
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
            if (!json.success) {
                that.errorElem.style.display = "block";
                that.errorElem.innerText = json.message;
            } else {
                that.errorElem.style.display = "none";
                location.href = "home";
            }
        });
    }

    clickError(e, that) {
        that.errorElem.style.display = "none";
    }
}

new Login();
