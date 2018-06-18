export default class LoadingView {
    constructor() {
        this.view = document.createElement("div");
        this.view.classList.add("loading");
        this.isDisporse = false;

        let loadingBody = document.createElement("div");
        loadingBody.classList.add("loadingBody");

        this.view.appendChild(loadingBody);
        document.body.appendChild(this.view);
    }

    show() {
        this.view.style.display = "block";
    }

    hide() {
        this.view.style.display = "none";
        if(this.isDisporse){
            this.view.remove();
        }
    }
}