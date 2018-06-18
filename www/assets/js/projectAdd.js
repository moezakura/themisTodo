class ProjectAdd {
    constructor() {
        this.projectAddForm = document.querySelector("#projectAdd");
        if (this.projectAddForm === undefined || this.projectAddForm == null)
            return;
        this.errorElem = document.querySelector("#error");

        let that = this;

        this.projectAddForm.addEventListener("submit", function(e){
            that.postLogin(e, that);
        }, true);
        this.errorElem.addEventListener("click", function(e){
            that.clickError(e, that);
        });
    }

    postLogin(e, that) {
        e.preventDefault();

        let formData = new FormData(that.projectAddForm);
        let projectAddJson = {
            "name": formData.get("name"),
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
                that.errorElem.style.display = "block";
                that.errorElem.innerText = json.message;
            } else {
                that.errorElem.style.display = "none";
                location.href = "/project/view/" + json.id;
            }
        });
    }

    clickError(e, that) {
        that.errorElem.style.display = "none";
    }
}
new ProjectAdd();