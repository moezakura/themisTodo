import Vue from "vue";
import BackView from "./backView";


if(document.querySelector("#taskHideListPopup")) {
    let backView = new BackView();

    let taskHideListPopup = new Vue({
        delimiters: ['${', '}'],
        el: '#taskHideListPopup',
        data: {
            taskHideListPopupFlag: false,
        },
        methods: {}
    });
    backView.addHideEvent(function(){
        taskHideListPopup.taskHideListPopupFlag = false;
    });

    new Vue({
        el: "#taskboardHideTaskShown",
        methods:{
            click(){
                taskHideListPopup.taskHideListPopupFlag = true;
                backView.show();
            }
        }
    })
}