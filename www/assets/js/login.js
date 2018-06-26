import Vue from 'vue';

new Vue({
    delimiters: ['${', '}'],
    el: '#login',
    data: {
        formError: null,
        formUsername: '',
        formPassword: '',
        hideStyle: 'none'
    },
    methods: {
        postLogin() {
            let loginJson = {
                "id": this.formUsername,
                "password": this.formPassword
            }
            fetch("", {
                method: 'POST',
                body: JSON.stringify(loginJson),
                credentials: "same-origin"
            }).then((response) => {
                return response.json()
            }).then((json) => {
                if(!json.success) {
                    this.formError = json.message
                    this.hideStyle = 'block'
                } else {
                    this.hideStyle = 'none'
                    location.href = 'home'
                }
            })
        },
        clickError() {
            this.hideStyle = 'none'
        }
    },
})