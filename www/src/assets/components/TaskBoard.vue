<template>
    <div>
        <h2 id="taskboardTitle"><i class="fas fa-tasks"></i><span>{{ project.name }}</span></h2>
        <a id="taskboardConfig" href="#"><i class="fas fa-cog"></i>Config</a>
        <a id="taskboardAdd" href="#"><i class="fas fa-plus-circle"></i>ADD</a>
        <div id="taskboard">
            <div id="taskBoardMinSized">
                <section id="todo">
                    <div class="statusName">Todo</div>
                    <ul class="taskList" data-status="TODO" ref="task-line-todo">
                        <task-line v-for="task in tasks.todo" :key="task.createDate" :task="task"
                                   :data-task-id="task.createDate"></task-line>
                    </ul>
                </section>
                <section id="doing">
                    <div class="statusName">Doing</div>
                    <ul class="taskList" data-status="DOING" ref="task-line-doing">
                        <task-line v-for="task in tasks.doing" :key="task.createDate" :task="task"
                                   :data-task-id="task.createDate"></task-line>
                    </ul>
                </section>
                <section id="pr">
                    <div class="statusName">PullRequest</div>
                    <ul class="taskList" data-status="PULL_REQUEST" ref="task-line-pr">
                        <task-line v-for="task in tasks.pullRequest" :key="task.createDate" :task="task"
                                   :data-task-id="task.createDate"></task-line>
                    </ul>
                </section>
                <section id="done">
                    <div class="statusName">Done</div>
                    <ul class="taskList" data-status="DONE" ref="task-line-done">
                        <task-line v-for="task in tasks.done" :key="task.createDate" :task="task"
                                   :data-task-id="task.createDate"></task-line>
                    </ul>
                </section>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
    import Project from "../scripts/model/api/project/Project"
    import ProjectApi from "../scripts/api/ProjectApi"
    import Task from "../scripts/model/api/task/Task"
    import TaskStatusConvert, {TaskStatus} from "../scripts/enums/TaskStatus"
    import TaskLine from "./TaskBoard/TaskLine"
    import Sortable from "sortablejs/Sortable"
    import TaskApi from "../scripts/api/TaskApi"

    export default {
        name: "TaskBoard",
        components: {TaskLine},
        data() {
            const project = new Project()
            const todo: Array<Task> = []
            const doing: Array<Task> = []
            const pullRequest: Array<Task> = []
            const done: Array<Task> = []

            return {
                project: project,
                tasks: {
                    todo: todo,
                    doing: doing,
                    pullRequest: pullRequest,
                    done: done,
                }
            }
        },
        computed: {
            projectId(): number {
                return this.$route.params["projectId"]
            }
        },
        methods: {
            loadProjectInfo() {
                this.$store.commit("incrementLoadingCount")
                ProjectApi.getProject(this.projectId).then(res => {
                    if (res.success) {
                        this.project = res.project
                    }
                    this.$store.commit("setCurrentProject", res.project)
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            },
            loadTasks() {
                this.$store.commit("incrementLoadingCount")

                this.tasks.todo.splice(0, this.tasks.todo.length)
                this.tasks.doing.splice(0, this.tasks.doing.length)
                this.tasks.pullRequest.splice(0, this.tasks.pullRequest.length)
                this.tasks.done.splice(0, this.tasks.done.length)

                ProjectApi.getTasks(this.projectId).then(res => {
                    if (res.success) {
                        for (const task of res.task) {
                            switch (task.status) {
                                case TaskStatusConvert.toNumber(TaskStatus.TODO):
                                    this.tasks.todo.push(task)
                                    break
                                case TaskStatusConvert.toNumber(TaskStatus.DOING):
                                    this.tasks.doing.push(task)
                                    break
                                case TaskStatusConvert.toNumber(TaskStatus.PULL_REQUEST):
                                    this.tasks.pullRequest.push(task)
                                    break
                                case TaskStatusConvert.toNumber(TaskStatus.DONE):
                                    this.tasks.done.push(task)
                                    break
                            }
                        }
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            },
            changedStatus(changeStatus: number, createDate) {
                let task = new Task()
                task.status = changeStatus
                this.$store.commit("incrementLoadingCount")

                TaskApi.Update(createDate, task).then(res => {
                    if (res.success) {
                        this.loadTasks()
                    }
                    //TaskBoard.loadTask(targetCreateId)
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            }
        },
        created() {
            this.loadProjectInfo()
            this.loadTasks()
        },
        mounted() {
            const taskLists = [
                this.$refs["task-line-todo"],
                this.$refs["task-line-doing"],
                this.$refs["task-line-pr"],
                this.$refs["task-line-done"],
            ]

            for (let taskLine of taskLists) {
                Sortable.create(taskLine, {
                    group: "shares",
                    onEnd: evt => {
                        const status = TaskStatusConvert.toNumber(evt.item.parentNode.dataset.status)
                        const targetCreateId = evt.item.dataset.taskId
                        this.changedStatus(status, targetCreateId)
                    },
                    animation: 100
                })
            }
        }
    }
</script>

<style scoped>

</style>