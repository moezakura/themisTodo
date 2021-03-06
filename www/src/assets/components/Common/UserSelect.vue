<template>
    <div class="user-select">
        <input v-model="name" @focus="" :readonly="!isEditing" @keydown.up.prevent="upKeyDown"
               @keydown.down.prevent="downKeyDown" @keydown.enter.prevent="enterKeyDown">
        <ul class="usersList userSearchDialog user-select-list" v-if="isEditing">
            <li v-for="(user, key) in userList" :class="{'select': key === selectedIndex}" @click="selectMouse(key)"
                @mouseenter="hoverMouse(key)">
                <div class="icon" :style="{'background-image': `url('/api/account/icon/${user.iconPath}')`}"></div>
                <div class="name">
                    <div class="nameId">{{ user.name }}</div>
                    <div class="displayName">{{ user.displayName }}</div>
                </div>
            </li>
        </ul>
    </div>
</template>

<script lang="ts">
    import AccountApi from "../../scripts/api/AccountApi"
    import AccountSearchRequest from "../../scripts/model/api/AccountSearchRequest"
    import User from "../../scripts/model/api/user/User"

    export default {
        name: "UserSelect",
        props: ['isShow', 'value', 'isInProject', 'readonly'],
        data: () => {
            const userList: Array<User> = []

            return {
                userList: userList,
                selectedIndex: 0,
            }
        },
        computed: {
            name: {
                get(): string {
                    if (this.value == undefined || this.value.displayName == undefined) {
                        return ""
                    }
                    return this.value.displayName
                },
                set(value) {
                    let changedValue = Object.assign(this.value, {})
                    changedValue.displayName = value
                    this.$emit('input', changedValue)
                    this.$emit('change', changedValue)

                    const searchRequest = new AccountSearchRequest()
                    searchRequest.name = value
                    searchRequest.displayName = value
                    searchRequest.project = this.projectId
                    searchRequest.isInProject = this.searchIsInProject
                    searchRequest.max = 20

                    AccountApi.search(searchRequest).then(res => {
                        this.userList = res
                    })
                }
            },
            isEditing(): boolean {
                return !(this.readonly)
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
            upKeyDown() {
                this.selectedIndex--
                if (this.selectedIndex < 0) this.selectedIndex = this.userList.length - 1
            },
            downKeyDown() {
                this.selectedIndex++
                if (this.selectedIndex > this.userList.length - 1) this.selectedIndex = 0
            },
            enterKeyDown() {
                let selectedUser = this.userList[this.selectedIndex]
                if (selectedUser === undefined) {
                    for (const user of this.userList) {
                        if (user.displayName == name) {
                            selectedUser = user
                        }
                    }
                    return
                }
                this.$emit('input', selectedUser)
                this.$emit('change', selectedUser)
                this.userList = []
            },
            hoverMouse(index: number) {
                this.selectedIndex = index
            },
            selectMouse(index: number) {
                const selectedUser = this.userList[index]
                if (selectedUser === undefined) {
                    return
                }
                this.$emit('input', selectedUser)
                this.$emit('change', selectedUser)
                this.userList = []
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
        float: left;
        width: calc(100% - 35px - 5px);

        input {
            font-size: 12px;
        }

        .user-select-list {
            position: relative;
            width: 100%;
            box-sizing: border-box;
            z-index: 120;
        }
    }
</style>