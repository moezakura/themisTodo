<template>
    <div>
        <form id="projectAdd" class="basicForm" @submit.prevent="newProject">
            <h2>New Project</h2>
            <div id="error" v-show="errorMessage!==undefined && errorMessage.length > 0" @click="clearMessage">{{ errorMessage }}</div>
            <input type="text" placeholder="Add Project Name" name="name" v-model="projectName">
            <textarea placeholder="Project Description" name="description" v-model="projectDescription"></textarea>
            <input type="submit" value="ADD">
        </form>
    </div>
</template>

<script lang="ts">
    import ProjectApi from "../scripts/api/ProjectApi"
    import ProjectAddRequest from "../scripts/model/api/ProjectAddRequest"

    export default {
        name: "NewProject",
        data: () => {
            return {
                errorMessage: "",
                projectName: "",
                projectDescription: "",
            }
        },
        methods: {
            clearMessage() {
                this.errorMessage = ""
            },
            newProject() {
                let addRequest = new ProjectAddRequest()
                addRequest.name = this.projectName
                addRequest.description = this.projectDescription
                this.$store.commit("incrementLoadingCount")
                this.clearMessage()

                ProjectApi.create(addRequest).then(res => {
                    if (res.success) {
                        this.$router.push({name: "taskBoard", params: {projectId: res.id}})
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