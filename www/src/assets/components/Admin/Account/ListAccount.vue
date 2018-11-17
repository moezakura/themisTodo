<template>
    <div class="basicForm">
        <h2>Account List</h2>
        <ul class="usersList user-select-list">
            <li v-for="(user, key) in userList">
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
    import AccountApi from "../../../scripts/api/AccountApi"

    export default {
        name: "ListAccount",
        data: () => {
            const userList: Array<Account> = []
            return {
                userList: userList
            }
        },
        created() {
            this.$store.commit("incrementLoadingCount")
            AccountApi.getList().then(res => {
                if (res.success) {
                    this.userList = res.users
                }
            }).finally(() => {
                this.$store.commit("decrementLoadingCount")
            })
        }
    }
</script>

<style lang="scss" scoped>
    h2 {
        text-align: center;
        letter-spacing: 1.4px;
        font-size: 22px;
    }

    .user-select-list {
        margin: 0;

        li {
            margin: 5px 0
        }
    }
</style>