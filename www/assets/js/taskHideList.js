import Vue from "vue";


if(document.querySelector("#taskHideListPopup")) {
    new Vue({
        delimiters: ['${', '}'],
        el: '#taskHideListPopup',
        data: {
            taskHideListPopupFlag: false,
        },
        methods: {}
    });
}