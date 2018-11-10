class AccountAdd {
    constructor() {
        this.accountAddForm = document.querySelector("#accountAdd");
        if (this.accountAddForm === undefined || this.accountAddForm == null)
            return;
        this.errorElem = document.querySelector("#error");
        this.accountAddedElem = this.accountAddForm.querySelector("#accountAdded");
        this.accountAddedIdInput = this.accountAddedElem.querySelector("#accountAddedId");
        this.accountAddedPwInput = this.accountAddedElem.querySelector("#accountAddedPassword");

        let that = this;

        this.accountAddForm.addEventListener("submit", function(e){
            that.postLogin(e, that);
        }, true);
        this.errorElem.addEventListener("click", function (){
            that.clickError(that);
        }, true);
    }

    postLogin(e, that) {
        e.preventDefault();

        let formData = new FormData(that.accountAddForm);
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
                that.errorElem.style.display = "block";
                that.errorElem.innerText = json.message;
            } else {
                that.errorElem.style.display = "none";
                that.accountAddedElem.style.display = "block";
                that.setAdded(json.name, json.password, that)
            }
        });
    }

    setAdded(id, password, that) {
        that.accountAddedIdInput.value = id;
        that.accountAddedPwInput.value = password;
    }

    clickError(that) {
        that.errorElem.style.display = "none";
    }
}

new AccountAdd();