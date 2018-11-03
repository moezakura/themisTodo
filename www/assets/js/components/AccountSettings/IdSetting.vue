<template>
    <form id="idChange" @submit.prevent="changeId">
        <div class="success v-shown" v-show="changeSuccess">Success changed</div>
        <div class="error v-shown" v-show="errorMessage.length > 0">{{ errorMessage }}</div>
        <label for="accountSettingsId">ID</label>
        <input type="text" id="accountSettingsId" name="accountSettingsId" v-model="userId">
        <i class="fas fa-check"></i>
        <input type="submit">
    </form>
</template>


<script lang="ts">
    import AccountApi from "../../accountApi"

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
                let changeObj = AccountApi.NewAccountObject()
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
            }
        }
    }
</script>

<style scoped>

</style>