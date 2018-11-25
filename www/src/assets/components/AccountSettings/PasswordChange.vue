<template>
    <div>
        <form id="passwordChange" @submit.prevent="submitPasswordChange">
            <div class="success v-shown" v-show="changeSuccess" @click="clearMessages">Success changed</div>
            <div class="error v-shown" v-show="errorMessage.length > 0" @click="clearMessages">{{ errorMessage }}</div>
            <label for="accountSettingsPassword">New Password</label>
            <input type="password" id="accountSettingsPassword" name="accountSettingsPassword" v-model="newPassword">
            <label for="accountSettingsPasswordRe">New Password</label>
            <input type="password" id="accountSettingsPasswordRe" name="accountSettingsPasswordRe"
                   v-model="newPasswordRe">
            <input type="submit" value="CHANGE">
        </form>

        <transition name="current-password-dialog" appear>
            <div v-show="newPasswordChange">
                <form id="currentPasswordDialog" class="basicForm" @submit.prevent="submitPasswordConfirm">
                    <div class="loading-circle" v-if="loading"></div>
                    <div v-else>
                        <p>Current password</p>
                        <input type="password" id="currentPasswordInput" v-model="currentPassword">
                        <input type="submit" value="CONFIRM">
                    </div>
                </form>
                <div id="backViewLayer" class="backView" @click="hideCurrentPasswordDialog"></div>
            </div>
        </transition>
    </div>
</template>

<script lang="ts">
    import AccountApi from "../../scripts/api/AccountApi"
    import Account from "../../scripts/model/Account"

    export default {
        name: "PasswordChange",
        data: () => {
            return {
                newPassword: "",
                newPasswordRe: "",
                currentPassword: "",
                newPasswordChange: false,
                changeSuccess: false,
                loading: false,
                errorMessage: "",
            }
        },
        methods: {
            showCurrentPasswordDialog() {
                this.loading = false
                this.newPasswordChange = true
            },
            hideCurrentPasswordDialog() {
                this.newPasswordChange = false
            },
            clearMessages() {
                this.errorMessage = ''
                this.changeSuccess = false
            },
            submitPasswordChange() {
                this.clearMessages()

                if (this.newPassword !== this.newPasswordRe) {
                    this.errorMessage = 'password is not match'
                } else {
                    this.errorMessage = ''
                    this.showCurrentPasswordDialog()
                }
            },
            submitPasswordConfirm() {
                let changeObj = new Account()
                changeObj.password = this.newPassword
                changeObj.currentPassword = this.currentPassword

                this.clearMessages()
                this.loading = true
                this.$store.commit("incrementLoadingCount")

                AccountApi.change(changeObj).then(json => {
                    if (!json.success) {
                        this.errorMessage = json.message
                        this.changeSuccess = false
                    } else {
                        this.errorMessage = ""
                        this.changeSuccess = true
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                    this.hideCurrentPasswordDialog()
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .current-password-dialog-enter-active,
    .current-password-dialog-leave-active {
        transition: opacity .3s ease;
    }

    .current-password-dialog-enter,
    .current-password-dialog-leave-to {
        opacity: 0;
    }

    .loading-circle{
        position: relative;
        top: 50%;
        left: 50%;
        transform:translate(-50%, -50%);
    }
</style>