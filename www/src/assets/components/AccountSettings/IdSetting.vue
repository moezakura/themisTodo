<template>
    <form id="idChange" @submit.prevent="changeId">
        <div class="success v-shown" v-show="changeSuccess" @click="clearMessages">Success changed</div>
        <div class="error v-shown" v-show="errorMessage.length > 0" @click="clearMessages">{{ errorMessage }}</div>
        <label for="accountSettingsId">ID</label>
        <input type="text" id="accountSettingsId" name="accountSettingsId" v-model="profileId">
        <input type="submit" value="CHANGE">
    </form>
</template>

<script lang="ts">
    import AccountApi from "../../scripts/api/AccountApi"
    import Account from "../../scripts/model/Account"
    import User from "../../scripts/model/api/user/User";

    export default {
        name: "NameSetting",
        data: () => {
            return {
                changeSuccess: false,
                errorMessage: "",
            }
        },
        computed: {
            profileId: {
                get() {
                    if (this.$store.getters.getMyProfile == undefined) {
                        return ""
                    }
                    return this.$store.getters.getMyProfile.name
                },
                set(value){
                    let profile: User = Object.assign({}, this.$store.getters.getMyProfile)
                    profile.name = value
                    this.$store.commit("setMyProfile", profile)
                }
            }
        },
        methods: {
            changeId() {
                let changeObj = new Account()
                changeObj.name = this.profileId
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