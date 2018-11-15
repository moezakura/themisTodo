<template>
    <form id="projectMemberAddForm" class="basicForm" autocomplete="off">
        <p>Member Invite</p>
        <user-select id="userSelect" v-model="selectUser" is-show="true" :is-in-project="false"></user-select>
        <p>Member Manage</p>
        <ul class="joined-users usersList">
            <li v-for="member in members">
                <div class="icon"
                     :style="{'background-image': `url('/api/account/icon/${member.iconPath}')`}"></div>
                <div class="name">
                    <div class="nameId">{{ member.name }}</div>
                    <div class="displayName">{{ member.displayName }}</div>
                </div>
                <div class="delete" @click="removeUser(member.uuid)">
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
    import AddMemberRequest from "../../../scripts/model/api/AddMemberRequest"
    import DeleteMemberRequest from "../../../scripts/model/api/DeleteMemberRequest"

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
                this.$store.commit("setProjectSettingsProps", {
                    key: "errorMessage",
                    value: ""
                })

                const addRequest = new AddMemberRequest()
                addRequest.uuid = value.uuid
                ProjectApi.addMemberToProject(this.project.uuid, addRequest).then(res => {
                    if (res.success) {
                        this.members.push(value)

                        const selectUser = new User()
                        selectUser.displayName = ""
                        selectUser.name = ""
                        selectUser.uuid = -1

                        this.selectUser = selectUser
                    } else {
                        this.$store.commit("setProjectSettingsProps", {
                            key: "errorMessage",
                            value: res.message
                        })
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
        },
        methods: {
            removeUser(userId) {
                this.$store.commit("incrementLoadingCount")
                this.$store.commit("setProjectSettingsProps", {
                    key: "errorMessage",
                    value: ""
                })

                const deleteRequest = new DeleteMemberRequest()
                deleteRequest.uuid = userId

                ProjectApi.removeMemberFromProject(this.project.uuid, deleteRequest).then(res => {
                    if (res.success) {
                        this.members = this.members.filter(member => {
                            return member.uuid != userId
                        })
                    } else {
                        this.$store.commit("setProjectSettingsProps", {
                            key: "errorMessage",
                            value: res.message
                        })
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            }
        }
    }
</script>

<style scoped>

</style>