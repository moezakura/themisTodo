<template>
    <div id="home">
        <div id="home-container">
            <section>
                <div class="sectionTitle">Your Projects</div>
                <ul class="joinProject">
                    <project-line v-for="project in projects" :key="project.uuid" :project="project"></project-line>
                </ul>
            </section>
            <section>
                <div class="sectionTitle taskListTitle">My Todo Tasks</div>
                <ul class="taskList" id="todoList">
                    <task-line v-for="task in todoList" :key="task.createDate" :task="task"
                               :hideAssign="true"></task-line>
                </ul>
            </section>
            <section>
                <div class="sectionTitle taskListTitle">My Doing Tasks</div>
                <ul class="taskList" id="doingList">
                    <task-line v-for="task in doingList" :key="task.createDate" :task="task"
                               :hideAssign="true"></task-line>
                </ul>
            </section>
        </div>
    </div>
</template>

<script lang="ts">
    import ProjectApi from "@scripts/api/ProjectApi"
    import {TaskStatus} from "@scripts/enums/TaskStatus"
    import Project from "@scripts/model/api/project/Project"
    import ProjectLine from "@components/Common/ProjectLine"
    import TaskLine from "@components/TaskBoard/TaskLine"
    import Task from "../scripts/model/api/task/Task"

    export default {
        name: "Dashboard",
        components: {
            ProjectLine,
            TaskLine
        },
        data: () => {
            const projects: Array<Project> = []
            const todoList: Array<Task> = []
            const doingList: Array<Task> = []

            return {
                todoList: todoList,
                doingList: doingList,
                projects: projects,
            }
        },
        methods: {
            loadTodoTask() {
                this.$store.commit("incrementLoadingCount")

                ProjectApi.getList(TaskStatus.TODO).then(res => {
                    if (!res.success) {
                        console.error("API ERROR")
                        return
                    }
                    this.todoList = res.task
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            },
            loadDoingTask() {
                this.$store.commit("incrementLoadingCount")

                ProjectApi.getList(TaskStatus.DOING).then(res => {
                    if (!res.success) {
                        console.error("API ERROR")
                        return
                    }
                    this.doingList = res.task
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            },
            loadingJoinedProject() {
                this.$store.commit("incrementLoadingCount")

                ProjectApi.getProject().then(res => {
                    if (!res.success) {
                        return
                    }
                    this.projects = res.project
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            },
        },
        created() {
            this.loadTodoTask()
            this.loadDoingTask()
            this.loadingJoinedProject()
        }
    }
</script>

<style scoped>

</style>