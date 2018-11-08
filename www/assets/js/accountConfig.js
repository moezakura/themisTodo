import Vue from 'vue';
import IdSetting from "./components/AccountSettings/IdSetting.vue"
import DisplayNameSetting from "./components/AccountSettings/DisplayNameSetting"
import PasswordChange from "./components/AccountSettings/PasswordChange"

if(document.querySelector("#accountSettings")) {
    new Vue({
        delimiters: ['${', '}'],
        el: '#accountSettings',
        data: {
        },
        components: {
            IdSetting,
            DisplayNameSetting,
            PasswordChange
        },
        created () {

        },
        methods : { },
    })
}