<template>
    <div>
        <h2 id="taskboardTitle"><i class="fas fa-tasks"></i><span>{{ project.name }}</span></h2>
        <a id="taskboardConfig" href="#"><i class="fas fa-cog"></i>Config</a>
        <a id="taskboardAdd" href="#"><i class="fas fa-plus-circle"></i>ADD</a>
        <div id="taskboard">
            <div id="taskBoardMinSized">
                <section id="todo">
                    <div class="statusName">Todo</div>
                    <ul class="taskList" data-status="todo"></ul>
                </section>
                <section id="doing">
                    <div class="statusName">Doing</div>
                    <ul class="taskList" data-status="doing"></ul>
                </section>
                <section id="pr">
                    <div class="statusName">PullRequest</div>
                    <ul class="taskList" data-status="pr"></ul>
                </section>
                <section id="done">
                    <div class="statusName">Done</div>
                    <ul class="taskList" data-status="done"></ul>
                </section>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
    import Project from "../scripts/model/api/project/Project"
    import ProjectApi from "../scripts/api/ProjectApi"

    export default {
        name: "TaskBoard",
        data() {
            const project = new Project()

            return {
                project: project,
            }
        },
        computed: {
            projectId(): number {
                return this.$route.params["projectId"]
            }
        },
        methods: {
            loadProjectInfo() {
                ProjectApi.getProject(this.projectId).then(res => {
                    this.project = res.project
                    this.$store.commit("setCurrentProject", res.project)
                })
            }
        },
        created() {
            this.loadProjectInfo()
        }
    }
</script>

<style scoped>

</style>