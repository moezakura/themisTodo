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
            <router-link to="/admin">
                <div>Admin Area</div>
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
                            <a>
                                <div>Logout</div>
                            </a>
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
            AccountApi.getProfile().then(res => {
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
                localStorage.removeItem("accessToken")
                this.$store.commit("setToken", "")
                this.$router.push({name: "welcome"})
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

    header {
        background-color: rgba(black, 0.4);
        background-repeat: no-repeat;
        height: $headerHeight;
        line-height: $headerHeight;
        overflow: hidden;
        margin-bottom: 10px;
        display: flex;

        .right {
            margin: 0 15px 0 auto;

            .profile-icon {
                width: 45px;
                height: 45px;
                box-sizing: border-box;
                margin-top: 10px;
                border-radius: 50%;
                border: solid 1px rgba(white, .2);
                background-position: center;
                background-size: cover;
                background-repeat: no-repeat;
            }
        }

        .left {
            a {
                display: block;
                float: left;
                height: $headerHeight - 10px;
                line-height: $headerHeight - 10px;
                margin: 5px 10px;
                text-decoration: none;
                color: white;

                img,
                div {
                    height: 100%;
                }

                div {
                    box-sizing: border-box;
                    padding: 0 20px;
                    letter-spacing: 1.5px;
                    font-size: 16px;
                    text-align: center;
                }

                &::before,
                &::after {
                    display: block;
                    float: left;
                    position: relative;
                    content: " ";
                }

                &::after {
                    position: relative;
                    width: 100%;
                    height: 50%;
                    top: -50%;
                    opacity: 0;
                    background: -moz-linear-gradient(top, accentColor(0) 0%, accentColor(1) 100%);
                    background: -webkit-linear-gradient(top, accentColor(0) 0%, accentColor(1) 100%);
                    background: linear-gradient(to bottom, accentColor(0) 0%, accentColor(1) 100%);
                    z-index: 0;
                    transition: ease opacity .3s;
                }

                &:hover::after {
                    opacity: 1;
                }

                &::before {
                    width: 1px;
                    height: 100%;
                    left: -10px;
                    margin-right: 10px;
                    background: -moz-linear-gradient(top, accentColor(0) 0%, accentColor(1) 20%, accentColor(1) 80%, accentColor(0) 100%);
                    background: -webkit-linear-gradient(top, accentColor(0) 0%, accentColor(1) 20%, accentColor(1) 80%, accentColor(0) 100%);
                    background: linear-gradient(to bottom, accentColor(0) 0%, accentColor(1) 20%, accentColor(1) 80%, accentColor(0) 100%);
                }

                &:first-child::before,
                &:first-child::after,
                &:nth-child(2)::before {
                    display: none;
                }
            }
        }

        .profile-menu {
            position: absolute;
            right: 15px;
            top: $headerHeight;
            width: 170px;
            background-color: $very-dark-gray;
            box-shadow: -2px -2px 5px rgba(black, .6);
            z-index: 100;

            margin: 0;
            padding: 0;

            li {
                height: 45px;
                line-height: 45px;
                padding: 0 12px;
                transition: ease background-color .3s;

                &:hover {
                    background-color: rgba(black, .3);
                }

                a {
                    text-decoration: none;
                    letter-spacing: 2px;
                    font-size: 14px;
                    color: white;
                }
            }
        }

        .profile-menu-back {
            position: fixed;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            z-index: 99;
        }

    }
</style>