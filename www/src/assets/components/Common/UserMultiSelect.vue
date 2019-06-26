<template>
    <div class="user-select">
        <input class="fake-submit" type="button" @click="showDialog" v-model="selectButtonLabel">

        <transition>
            <div class="user-multi-select-dialog-container" v-if="isDialogShow">
                <div class="user-multi-select-dialog">
                    <input v-model="searchValue" @keydown.up.prevent="upKeyDown" @keydown.down.prevent="downKeyDown"
                           @keydown.enter.prevent="enterKeyDown">
                    <ul class="user-select-list usersList">
                        <li v-for="(user, key) in displayUsers" :class="{'select': key === selectedIndex}"
                            @click="selectMouse(key)" @mouseenter="hoverMouse(key)" class="users-list-row">
                            <div class="checkbox" :class="{'checked': user.check}">
                                <i class="fas fa-check"></i>
                            </div>
                            <div class="icon"
                                 :style="{'background-image': `url('/api/account/icon/${user.iconPath}')`}"></div>
                            <div class="name">
                                <div class="nameId">{{ user.name }}</div>
                                <div class="displayName">{{ user.displayName }}</div>
                            </div>
                        </li>
                    </ul>
                </div>
                <div class="user-multi-select-dialog-background" @click="hideDialog"></div>
            </div>
        </transition>
    </div>
</template>

<script lang="ts">
    import User from "../../scripts/model/api/user/User"
    import ProjectApi from "@scripts/api/ProjectApi";
    import UserWithCheck from "@scripts/model/UserWithCheck";

    interface UserMultiSelectData {
        selectedIndex: number,
        isDialogShow: boolean,
        selectButtonLabel: string,
        selectedUserList: Array<User>,
        userList: Array<UserWithCheck>,
        searchValue: string,
    }

    export default {
        name: "UserMultiSelect",
        props: ['value', 'isInProject'],
        data(): UserMultiSelectData {
            return {
                selectedIndex: 0,
                isDialogShow: false,
                selectButtonLabel: "Select users",
                selectedUserList: [],
                userList: [],
                searchValue: "",
            }
        },
        computed: {
            displayUsers(): Array<UserWithCheck> {
                const returnList = new Array<UserWithCheck>()
                const searchValue = this.searchValue.toLowerCase()
                if (searchValue === "") {
                    return this.userList
                }

                this.userList.forEach((user: UserWithCheck) => {
                    const name = user.name.toLowerCase()
                    const displayName = user.displayName.toLowerCase()
                    if (displayName.includes(searchValue) || name.includes(searchValue)) {
                        returnList.push(user)
                    }
                })
                return returnList
            },
            searchIsInProject(): boolean {
                if (this.isInProject === undefined) {
                    return true
                }
                return this.isInProject
            },
            projectId(): number | undefined {
                if (this.$store.getters.getCurrentProject === undefined) {
                    return
                }
                const project = this.$store.getters.getCurrentProject
                return project.uuid
            }
        },
        methods: {
            async load(): Promise<void> {
                const res = await ProjectApi.getMembers(this.projectId)
                if (res.success) {
                    this.userList = new Array<UserWithCheck>()
                    res.members.forEach((user: User) => {
                        let userWithCheck: UserWithCheck = new UserWithCheck(user)
                        userWithCheck.check = false
                        if (Array.isArray(this.value)) {
                            const valueSelected = <Array<UserWithCheck>>this.value
                            valueSelected.some((valueUser: UserWithCheck) => {
                                if (valueUser.check && valueUser.uuid === user.uuid) {
                                    userWithCheck.check = true
                                    return true
                                }
                                return false
                            })
                        }
                        this.userList.push(userWithCheck)
                    })
                }
            },
            upKeyDown(): void {
                this.selectedIndex--
                if (this.selectedIndex < 0) this.selectedIndex = this.displayUsers.length - 1
            },
            downKeyDown(): void {
                this.selectedIndex++
                if (this.selectedIndex > this.displayUsers.length - 1) this.selectedIndex = 0
            },
            enterKeyDown(): void {
                let selectedUser: UserWithCheck = this.displayUsers[this.selectedIndex]
                if (selectedUser === undefined) {
                    for (const user of this.displayUsers) {
                        if (user.displayName == name) {
                            selectedUser = user
                        }
                    }
                    return
                }
                this.selectUser(selectedUser.uuid)
            },
            hoverMouse(index: number): void {
                this.selectedIndex = index
            },
            selectMouse(index: number): void {
                const selectedUser: UserWithCheck = this.displayUsers[index]
                if (selectedUser === undefined) {
                    return
                }
                this.selectUser(selectedUser.uuid)
            },
            selectUser(uuid: number): void {
                let selectedUsers = new Array<UserWithCheck>()
                this.userList.forEach((user: UserWithCheck, index: number) => {
                    if (user.uuid === uuid) {
                        this.$set(this.userList[index], "check", !user.check)
                    }

                    if (user.check) {
                        selectedUsers.push(user)
                    }
                })

                this.selectButtonLabel = selectedUsers.length > 0 ? selectedUsers.length + " users" : "Select users"

                this.$emit('input', selectedUsers)
                this.$emit('change', selectedUsers)
            },
            showDialog(): void {
                this.load()
                this.isDialogShow = true
            },
            hideDialog(): void {
                this.isDialogShow = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .v-enter,
    .v-leave-to {
        opacity: 0;
    }

    .v-enter-active,
    .v-leave-active {
        transition: opacity .2s;
    }

    .user-select {
        input {
            font-size: 12px;
        }
    }

    .user-multi-select-dialog-container {
        .user-multi-select-dialog {
            position: fixed;
            top: 50%;
            left: 50%;
            width: 400px;
            height: 300px;
            transform: translate(-50%, -50%);
            background: $backgroundColor;
            z-index: 50;

            input {
                margin: 10px auto;
                width: 95%;
            }

            .user-select-list {
                margin: 10px auto;
                width: 95%;
                height: calc(100% - #{40px + 10px * 2});
                box-sizing: border-box;
                overflow-y: auto;

                .users-list-row {
                    display: flex;
                    -webkit-user-select: none;
                    -moz-user-select: none;
                    -ms-user-select: none;
                    user-select: none;

                    .checkbox {
                        width: 20px;
                        height: 20px;
                        font-size: 16px;
                        line-height: 20px;
                        text-align: center;
                        border: solid 1px white;
                        margin: 7.5px 5px;
                        box-sizing: border-box;

                        i {
                            opacity: 0;
                            transition: ease opacity .1s;
                        }

                        &.checked {
                            i {
                                opacity: 1;
                            }
                        }
                    }

                    .name {
                        width: calc(100% - #{20px + 5px * 2 + 35px});
                    }
                }
            }
        }

        .user-multi-select-dialog-background {
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(black, .5);
            z-index: 49;
        }
    }
</style>