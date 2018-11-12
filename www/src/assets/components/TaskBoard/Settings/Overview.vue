<template>
    <form id="projectConfigForm" class="basicForm" autocomplete="off" @submit.prevent="changeProject">
        <p>Name</p>
        <input type="text" name="name" v-model="projectName">
        <p>Description</p>
        <textarea name="description" v-model="projectDescription"></textarea>
        <input type="submit" value="SAVE">
    </form>
</template>

<script lang="ts">
    import Project from "../../../scripts/model/api/project/Project"
    import ProjectApi from "../../../scripts/api/ProjectApi"
    import ProjectUpdateRequest from "../../../scripts/model/ProjectUpdateRequest"

    export default {
        name: "Overview",
        computed: {
            project(): Project {
                return this.$store.getters.getCurrentProject
            }
        },
        data: () => {
            return {
                projectName: "",
                projectDescription: "",
            }
        },
        created() {
            this.projectName = this.project.name
            this.projectDescription = this.project.description
        },
        methods: {
            changeProject() {
                this.$store.commit("incrementLoadingCount")
                this.$store.commit("setProjectSettingsProps", {
                    key: "errorMessage",
                    value: "",
                })

                let updateRequest = new ProjectUpdateRequest()

                updateRequest.name = this.projectName
                updateRequest.description = this.projectDescription

                ProjectApi.updateProject(this.project.uuid, updateRequest).then(res => {
                    if (res.success) {
                        const storeProject = this.$store.getters.getCurrentProject
                        let project = Object.assign({}, storeProject)
                        project.name = updateRequest.name
                        project.description = updateRequest.description
                        this.$store.commit("setCurrentProject", project)
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

<style lang="scss" scoped>
    #projectConfigForm {
        textarea {
            height: 250px;
            resize: none;
        }
    }
</style>