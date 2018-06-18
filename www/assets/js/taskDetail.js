import BackView from "./backView"

export default class TaskDetail {
    static show() {
        let taskPopup = document.querySelector("#taskPopup");
        taskPopup.style.display = "block";

        let backView = new BackView();
        backView.isDisporse = true;
        backView.show();

        backView.addWithHideElem(taskPopup);
    }

    static load(taskId) {

    }

    static loadAndShow(taskId) {
        this.show();
        this.load(taskId);
    }

    static set(taskObject) {

    }
}