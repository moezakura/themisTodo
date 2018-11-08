<template>
    <form id="idChange" @submit.prevent="changeId">
        <div class="success v-shown" v-show="changeSuccess" @click="clearMessages">Success changed</div>
        <div class="error v-shown" v-show="errorMessage.length > 0" @click="clearMessages">{{ errorMessage }}</div>
        <label for="accountSettingsId">ID</label>
        <input type="text" id="accountSettingsId" name="accountSettingsId" v-model="userId">
        <input type="submit" value="CHANGE">
    </form>
</template>

<script lang="ts">
    import AccountApi from "../../api/AccountApi"
    import Account from "../../model/Account"

    export default {
        name: "NameSetting",
        data: () => {
            return {
                userId: "",
                changeSuccess: false,
                errorMessage: "",
            }
        },
        methods: {
            changeId() {
                let changeObj = new Account()
                changeObj.name = this.userId
                AccountApi.Change(changeObj).then(json => {
                    if (!json.success) {
                        this.errorMessage = json.message
                        this.changeSuccess = false
                    } else {
                        this.errorMessage = ""
                        this.changeSuccess = true
                    }
                })
            },
            clearMessages() {
                this.changeSuccess = false
                this.errorMessage = ''
            }
        }
    }
</script>

<style scoped>

</style>