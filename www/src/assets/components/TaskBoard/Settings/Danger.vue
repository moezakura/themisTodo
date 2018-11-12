<template>
    <div class="basicForm">
        <p>This action cannot be undone. This will permanently delete the <strong>{{ projectName }}</strong>,
            project, todo, settings and other all, and remove all project associations.</p>
        <p>Please type in the name of the project to confirm.</p>
        <input type="text" v-model="typeProjectName" @keydown.enter="deleteProject">
        <div id="deleteProject" class="deleteButton" :disabled="isDisable" @click="deleteProject">DELETE PROJECT</div>
    </div>
</template>

<script lang="ts">
    import ProjectApi from "../../../scripts/api/ProjectApi"

    export default {
        name: "Danger",
        data: () => {
            return {
                typeProjectName: ""
            }
        },
        computed: {
            projectName(): string {
                return this.$store.getters.getCurrentProject.name
            },
            isDisable(): boolean {
                return this.$store.getters.getCurrentProject.name != this.typeProjectName
            }
        },
        methods: {
            deleteProject() {
                if (this.isDisable) {
                    return
                }
                this.$store.commit("incrementLoadingCount")
                const projectId = this.$store.getters.getCurrentProject.uuid
                ProjectApi.deleteProject(projectId).then(res => {
                    if (res.success) {
                        this.$router.push({name: "dashboard"})
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    strong {
        font-weight: 900;
        font-size: 17px;
    }

    .basicForm {
        input {
            margin-bottom: 35px;
            &:focus {
                border-color: rgb(239, 83, 80);
            }
        }
    }

    #deleteProject {
        &[disabled] {
            opacity: .7;
            border-color: rgb(180, 180, 180);
            background-color: rgba(255, 255, 2550, 0.1);

            &:hover {
                border-color: rgb(180, 180, 180);
                background-color: rgba(255, 255, 2550, 0.1);
            }
        }
    }
</style>