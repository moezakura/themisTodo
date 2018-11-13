<template>
    <form id="projectMemberAddForm" class="basicForm" autocomplete="off">
        <p>Member Invite</p>
        <user-select id="userSelect" v-model="selectUser" is-show="true" :is-in-project="false"></user-select>
        <p>Member Manage</p>
        <ul class="joined-users usersList">
            <li v-for="member in members">
                <div class="icon"
                     :style="{'background-image': `url('/api/assets/accountIcon/${member.uuid}.png')`}"></div>
                <div class="name">
                    <div class="nameId">{{ member.name }}</div>
                    <div class="displayName">{{ member.displayName }}</div>
                </div>
                <div class="delete">
                    <i class="fas fa-minus-circle"></i>
                </div>
            </li>
        </ul>
    </form>
</template>

<script lang="ts">
    import UserSelect from "../../Common/UserSelect"
    import User from "../../../scripts/model/api/user/User"
    import ProjectApi from "../../../scripts/api/ProjectApi"
    import Project from "../../../scripts/model/api/project/Project"
    import AddMemberRequest from "../../../scripts/model/api/AddMemberRequest";

    export default {
        name: "Member",
        components: {UserSelect},
        data: () => {
            const selectUser = new User()
            selectUser.displayName = ""
            selectUser.name = ""
            selectUser.uuid = -1
            const members: Array<User> = []

            return {
                selectUser: selectUser,
                members: members
            }
        },
        computed: {
            project(): Project {
                return this.$store.getters.getCurrentProject
            }
        },
        watch: {
            selectUser(value) {
                if (value.uuid <= 0) {
                    return
                }

                this.$store.commit("incrementLoadingCount")
                const updateRequest = new AddMemberRequest()
                updateRequest.uuid = value.uuid
                ProjectApi.addMemberToProject(this.project.uuid, updateRequest).then(res => {
                    if (res.success) {
                        this.members.push(value)

                        const selectUser = new User()
                        selectUser.displayName = ""
                        selectUser.name = ""
                        selectUser.uuid = -1

                        this.selectUser = selectUser
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            }
        },
        created() {
            this.$store.commit("incrementLoadingCount")
            ProjectApi.getMembers(this.project.uuid).then(res => {
                if (res.success) {
                    this.members = res.members
                }
            }).finally(() => {
                this.$store.commit("decrementLoadingCount")
            })
        }
    }
</script>

<style scoped>

</style>