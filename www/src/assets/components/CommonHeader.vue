<template>
    <header>
        <div class="left">
            <router-link to="/"><img src="/assets/images/logo.png"></router-link>
            <router-link to="/home">
                <div>Home</div>
            </router-link>
            <router-link to="/project/new">
                <div>New Project</div>
            </router-link>
            <router-link to="/account/add">
                <div>New Account</div>
            </router-link>
        </div>
        <div class="right">
            <div class="profile-icon" @click="profileMenuToggle"
                 :style="{ 'background-image': `url('${userIconPath}')`}"></div>
            <transition>
                <div v-show="isShowUserMenu" @click="profileMenuHide">
                    <ul class="profile-menu">
                        <li>
                            <router-link to="/settings">
                                <div>Your Profile</div>
                            </router-link>
                        </li>
                        <li @click="logout">
                            <router-link to="/">
                                <div>Logout</div>
                            </router-link>
                        </li>
                    </ul>
                    <div class="profile-menu-back" @click="profileMenuHide"></div>
                </div>
            </transition>
        </div>
    </header>
</template>

<script lang="ts">
    import User from "../scripts/model/api/user/User"
    import AccountApi from "../scripts/api/AccountApi"

    export default {
        name: "CommonHeader",
        data: () => {
            return {
                isShowUserMenu: false
            }
        },
        computed: {
            myProfile(): User | undefined {
                return this.$store.getters.getMyProfile
            },
            userIconPath(): string {
                if (this.myProfile == undefined || this.myProfile.iconPath == undefined ||
                    this.myProfile.iconPath.length <= 0) {
                    return ""
                }
                return `/api/account/icon/${this.myProfile.iconPath}`
            }
        },
        created() {
            this.$store.commit("incrementLoadingCount")
            AccountApi.GetProfile().then(res => {
                if (res.success) {
                    this.$store.commit("setMyProfile", res.user)
                }
            }).finally(() => {
                this.$store.commit("decrementLoadingCount")
            })
        },
        methods: {
            profileMenuToggle() {
                this.isShowUserMenu = !this.isShowUserMenu
            },
            profileMenuHide() {
                this.isShowUserMenu = false
            },
            logout() {

            }
        }
    }
</script>

<style lang="scss" scoped>
    .v-enter,
    .v-leave-to {
        opacity: 0;
        div, ul {
            transform: translateY(-5px);
        }
    }

    .v-enter-active,
    .v-leave-active {
        transition: ease all .2s;
        div, ul {
            transition: ease all .2s;
        }
    }
</style>