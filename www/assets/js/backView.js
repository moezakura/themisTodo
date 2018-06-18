export default class BackView {
    constructor() {
        this.view = document.createElement("div");
        this.view.classList.add("backView");
        this.withHideElem = [];
        this.hideEvents = [];
        this.isDisporse = false;

        let that = this;
        this.view.addEventListener("click", function(){
            that.hide();
        }, true);

        document.body.appendChild(this.view);
    }

    show() {
        this.display("block");
    }

    hide() {
        this.display("none");
        this.withHideElem.forEach(function(value){ value.style.display = "none"; });
        this.hideEvents.forEach(function(value){ value(); });
        if(this.isDisporse){
            this.view.remove();
        }
    }

    addWithHideElem(elem) {
        this.withHideElem.push(elem);
    }

    removeWithHideElem(elem) {
        let withElemTemp = this.withHideElem.slice(0, this.withHideElem.length);
        withElemTemp.forEach(function (v, i) {
            if(this.withHideElem[i] === v)
                this.withHideElem.slice(i, 1);
        });
    }

    addHideEvent(event){
        this.hideEvents.push(event);
    }

    display(displayText) {
        this.view.style.display = displayText;
    }
}