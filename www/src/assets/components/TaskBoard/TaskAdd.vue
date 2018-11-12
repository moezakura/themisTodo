<template>
    <form id="taskboardAddForm" class="basicForm" @submit.prevent="submitAdd">
        <h3>New Task</h3>
        <div id="taskboardAddClose" @click="hideTaskAdd"><i class="fas fa-angle-double-right"></i></div>
        <div class="error" v-show="errorMessage !== undefined && errorMessage.length > 0" @click="clearMessage">{{
            errorMessage }}
        </div>
        <p>Name</p>
        <input type="text" name="name" v-model="form.name">
        <p class="label">Creator:<span>{{ displayName }}</span></p>
        <p class="label">ID:<span id="createTaskId"></span></p>
        <p>Assign</p>
        <user-select :is-show="true" :is-in-project="true" v-model="form.selectUser" id="assign"></user-select>
        <p>Deadline</p>
        <input type="date" value="" name="deadline" v-model="form.deadline">
        <p>Description</p>
        <textarea name="description" v-model="form.description"></textarea>
        <input type="submit" value="ADD">
    </form>
</template>

<script lang="ts">
    import TaskApi from "../../scripts/api/TaskApi"
    import TaskAddRequest from "../../scripts/model/api/TaskAddRequest"
    import Project from "../../scripts/model/api/project/Project"
    import User from "../../scripts/model/api/user/User"
    import UserSelect from "@components/Common/UserSelect"

    export default {
        name: "TaskAdd",
        props: ['value'],
        components: {
            UserSelect
        },
        data: () => {
            const selectUser = new User()
            selectUser.displayName = ""
            selectUser.name = ""
            selectUser.uuid = -1
            return {
                errorMessage: "",
                form: {
                    name: "",
                    deadline: "",
                    description: "",
                    selectUser: selectUser,
                }
            }
        },
        computed: {
            displayName(): string {
                return ""
            },
            project(): Project | undefined {
                if (this.$store.getters.getCurrentProject === undefined) {
                    return
                }
                return this.$store.getters.getCurrentProject
            }
        },
        methods: {
            hideTaskAdd() {
                this.$emit("input", false)
                this.$emit("change", false)
            },
            clearMessage() {
                this.errorMessage = ""
            },
            clearForm() {
                const selectUser = new User()
                selectUser.displayName = ""
                selectUser.name = ""
                selectUser.uuid = -1
                return {
                    name: "",
                    deadline: "",
                    description: "",
                    selectUser: selectUser,
                }
            },
            submitAdd() {
                this.$store.commit("incrementLoadingCount")
                let addRequest = new TaskAddRequest()
                addRequest.name = this.form.name
                addRequest.assign = this.form.selectUser.uuid
                addRequest.description = this.form.description
                addRequest.deadline = this.form.deadline
                addRequest.projectId = this.project.uuid
                TaskApi.Create(addRequest).then(res => {
                    if (res.success) {
                        this.form = this.clearForm()
                    } else {
                        this.errorMessage = res.message
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
</style>