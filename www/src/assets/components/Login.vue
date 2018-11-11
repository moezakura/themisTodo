<template>
    <div id="login">
        <form @submit.prevent="postLogin" class="basicForm">
            <h2>LOGIN</h2>
            <p id="error" @click="clearMessage">{{ errorMessage }}</p>
            <input type="text" name="id" v-model="form.userName" placeholder="id">
            <input type="password" name="pw" v-model="form.userPassword" placeholder="password">
            <input type="submit" value="LOGIN">
        </form>
    </div>
</template>

<script lang="ts">
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
                let loginJson = {
                    "id": this.formUsername,
                    "password": this.formPassword
                }

                fetch("/api/login", {
                    method: 'POST',
                    body: JSON.stringify(loginJson),
                    credentials: "same-origin"
                }).then((response) => {
                    return response.json()
                }).then((json) => {
                    if (!json.success) {
                        this.formError = json.message
                        this.hideStyle = 'block'
                    } else {
                        this.hideStyle = 'none'
                        location.href = 'home'
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