<template>
    <div id="login">
        <form @submit.prevent="postLogin" class="basicForm">
            <h2>LOGIN</h2>
            <p id="error" v-show="errorMessage && errorMessage.length > 0" @click="clearMessage">{{ errorMessage }}</p>
            <input type="text" name="id" v-model="form.userName" placeholder="id">
            <input type="password" name="pw" v-model="form.userPassword" placeholder="password">
            <input type="submit" value="LOGIN">
        </form>
    </div>
</template>

<script lang="ts">
    import LoginRequest from "@scripts/model/api/LoginRequest"
    import AuthApi from "../scripts/api/AuthApi"

    export default {
        name: "Login",
        data: () => {
            return {
                errorMessage: "",
                form: {
                    userName: "",
                    userPassword: ""
                }
            }
        },
        methods: {
            clearMessage() {
                this.errorMessage = ""
            },
            postLogin() {
                this.$store.commit('incrementLoadingCount')
                const loginRequest = new LoginRequest()
                loginRequest.id = this.form.userName
                loginRequest.password = this.form.userPassword

                AuthApi.Login(loginRequest).then(res => {
                    if (!res.success) {
                        this.errorMessage = res.message
                    } else {
                        this.$router.push({name: 'dashboard'})
                    }
                }).finally(() => {
                    this.$store.commit('decrementLoadingCount')
                })
            }
        }
    }
</script>

<style scoped>

</style>