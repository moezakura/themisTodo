<template>
    <div id="home">
        <section>
            <div class="sectionTitle">Your Projects</div>
            <ul class="joinProject">
                <project-line v-for="project in projects" :project="project"></project-line>
            </ul>
        </section>
        <section>
            <div class="sectionTitle taskListTitle">My Todo Tasks</div>
            <ul class="taskList" id="todoList">
                <parent :todo="todo" v-for="todo in todoList"/>
            </ul>
        </section>
        <section>
            <div class="sectionTitle taskListTitle">My Doing Tasks</div>
            <ul class="taskList" id="doingList">
                <parent :todo="todo" v-for="todo in doingList"/>
            </ul>
        </section>
    </div>
</template>

<script lang="ts">
    import ProjectApi from "@scripts/api/ProjectApi"
    import {TaskStatus} from "@scripts/enums/TaskStatus"
    import Project from "@scripts/model/api/project/Project"
    import ProjectLine from "@components/Common/ProjectLine";

    export default {
        name: "Dashboard",
        components: {ProjectLine},
        data: () => {
            const projects: Array<Project> = []

            return {
                todoList: [],
                doingList: [],
                projects: projects,
            }
        },
        created() {
            // ProjectApi.getList(TaskStatus.TODO).then(res => {
            //     if (!res.success) {
            //         console.error("API ERROR")
            //         return
            //     }
            //     this.todoList = res
            // })
            // ProjectApi.getList(TaskStatus.DOING).then(json => {
            //     if (!json.success) {
            //         console.error("API ERROR")
            //         return
            //     }
            //     this.doingList = json.task
            // })

            ProjectApi.getProject().then(res => {
                console.log(res)
                if (!res.success) {
                    return
                }
                this.projects = res.project
            })
        },
    }
</script>

<style scoped>

</style>