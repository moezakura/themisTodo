class LoadingView {
    constructor() {
        this.view = document.createElement("div");
        this.view.classList.add("loading");

        let loadingBody = document.createElement("div");
        loadingBody.classList.add("loadingBody");

        this.view.appendChild(loadingBody);
    }

    show() {
        this.view.style.display = "block";
    }

    hide() {
        this.view.style.display = "none";
    }
}