<template>
    <form id="accountAdd" class="basicForm" @submit.prevent="createAccount">
        <h2>New Account</h2>
        <div id="error" v-show="errorMessage !== undefined && errorMessage.length > 0" @click="clearMessages">{{
            errorMessage }}
        </div>
        <input type="text" placeholder="Add Account ID" name="name" v-model="newUserId">
        <input type="submit" value="ADD">
        <div id="accountAdded" v-if="added.userId.length > 0 && added.password.length > 0">
            <h3>Added Account</h3>
            <p>ID</p>
            <input type="text" readonly id="accountAddedId" v-model="added.userId">
            <p>Password</p>
            <input type="text" readonly id="accountAddedPassword" v-model="added.password">
        </div>
    </form>
</template>

<script lang="ts">
    import AccountApi from "../../../scripts/api/AccountApi"
    import AccountCreateRequest from "../../../scripts/model/api/AccountCreateRequest"

    export default {
        name: "NewAccount",
        data: () => {
            return {
                errorMessage: "",
                newUserId: "",
                added: {
                    userId: "",
                    password: "",
                }
            }
        },
        methods: {
            clearMessages() {
                this.errorMessage = ""
                this.$set(this.added, 'userId', '')
                this.$set(this.added, 'password', '')
            },
            createAccount() {
                this.clearMessages()
                let createRequest = new AccountCreateRequest()
                createRequest.name = this.newUserId
                this.$store.commit("incrementLoadingCount")

                AccountApi.create(createRequest).then(res => {
                    if (res.success) {
                        this.$set(this.added, 'userId', res.name)
                        this.$set(this.added, 'password', res.password)
                    } else {
                        this.errorMessage = res.message
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            }
        }
    }
</script>

<style lang="scss" scoped>

</style>