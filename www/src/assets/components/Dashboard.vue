<template>
    <div id="home">
        <div id="home-container">
            <section class="home-left-section">
                <div class="sectionTitle">Your Projects</div>
                <ul class="joinProject">
                    <project-line v-for="project in projects" :key="project.uuid" :project="project"></project-line>
                </ul>
            </section>
            <div class="home-right-section">
                <section class="">
                    <div class="sectionTitle taskListTitle">Current Doing Task</div>
                    <task-timer-user-dashboard></task-timer-user-dashboard>
                </section>
                <div class="home-todo-list">
                    <section class="todo-list">
                        <div class="sectionTitle taskListTitle">My Todo Tasks</div>
                        <ul class="taskList" id="todoList">
                            <li v-for="task in todoList">
                                <task-line :key="task.createDate" :task="task"
                                           :hideAssign="true" :full-deadline="true"></task-line>
                            </li>
                        </ul>
                    </section>
                    <section class="todo-list">
                        <div class="sectionTitle taskListTitle">My Doing Tasks</div>
                        <ul class="taskList" id="doingList">
                            <li v-for="task in doingList">
                                <task-line :key="task.createDate" :task="task" :hideAssign="true"
                                           :full-deadline="true"></task-line>
                            </li>
                        </ul>
                    </section>
                </div>
            </div>
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
    import TaskTimerUserDashboard from "@components/TaskTimer/TaskTimerUserDashboard.vue";

    export default {
        name: "Dashboard",
        components: {
            TaskTimerUserDashboard,
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

                ProjectApi.getProjects().then(res => {
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

<style scoped lang="scss">
    .sectionTitle {
        font-size: 16px;
        letter-spacing: 1.3px;
        height: 35px;
        line-height: 35px;

        i {
            color: white;
            margin-left: 10px;
            line-height: 35px;
        }
    }

    .taskListTitle {
        height: $buttonHeight;
        line-height: $buttonHeight;
        background-color: rgba(black, .3);

        &::before {
            width: 5px;
            height: 100%;
            display: inline-block;
            vertical-align: top;
            margin-right: 10px;
            content: " ";
            background-color: $accentColor;
        }
    }


    #home {
        height: calc(100% - 75px);
        overflow-x: auto;
    }

    #home-container {
        display: flex;
        justify-content: flex-start;
        margin: 0 15px;
        min-width: 1100px;
    }

    .home-left-section, .home-right-section {
        flex-grow: 1;
        margin: 0 10px;
    }

    .home-left-section {
        width: 40%;
    }

    .home-right-section {
        width: 60%;

        * {
            justify-content: unset;
            flex-grow: 0;
        }

    }

    .home-todo-list {
        display: flex;

        .todo-list {
            width: calc((100% - 10px) / 2);

            &:first-child {
                margin-right: 10px;
            }
        }

        ul {
            &.joinProject {
                display: block;
                list-style: none;
                margin: 0 0 0 15px;
                padding: 0;
            }

            &.taskList {
                width: 100%;

                li {
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                }
            }
        }
    }
</style>