<template>
    <form id="displayNameChange" @submit.prevent="changeDisplayName">
        <div class="success v-shown" v-show="changeSuccess" @click="clearMessages">Success changed</div>
        <div class="error v-shown" v-show="errorMessage.length > 0" @click="clearMessages">{{ errorMessage }}</div>
        <label for="accountSettingsDisplayName">Display name</label>
        <input type="text" id="accountSettingsDisplayName" name="accountSettingsDisplayName" v-model="userName">
        <input type="submit" value="CHANGE">
    </form>
</template>

<script lang="ts">
    import AccountApi from "../../scripts/api/AccountApi"
    import Account from "../../scripts/model/Account"
    import User from "../../scripts/model/api/user/User"

    export default {
        name: "DisplayNameSetting",
        data: () => {
            return {
                changeSuccess: false,
                errorMessage: "",
            }
        },
        computed: {
            userName: {
                get() {
                    if (this.$store.getters.getMyProfile == undefined) {
                        return ""
                    }
                    return this.$store.getters.getMyProfile.displayName
                },
                set(value) {
                    let profile: User = Object.assign({}, this.$store.getters.getMyProfile)
                    profile.displayName = value
                    this.$store.commit("setMyProfile", profile)
                }
            }
        },
        methods: {
            changeDisplayName() {
                let changeObj = new Account()
                changeObj.displayName = this.userName
                this.$store.commit("incrementLoadingCount")
                AccountApi.Change(changeObj).then(json => {
                    if (!json.success) {
                        this.errorMessage = json.message
                        this.changeSuccess = false
                    } else {
                        this.errorMessage = ""
                        this.changeSuccess = true
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
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