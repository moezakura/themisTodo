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

<script>
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

            }
        }
    }
</script>

<style scoped>

</style>