<template>
    <div>
        <h2 id="taskboardTitle"><i class="fas fa-tasks"></i><span>{{ project.name }}</span></h2>
        <a id="taskboardConfig" href="#"><i class="fas fa-cog"></i>Config</a>
        <a id="taskboardAdd" href="#"><i class="fas fa-plus-circle"></i>ADD</a>
        <div id="taskboard">
            <div id="taskBoardMinSized">
                <section id="todo">
                    <div class="statusName">Todo</div>
                    <ul class="taskList" data-status="todo">
                        <task-line v-for="task in tasks.todo" :key="task.createDate" :task="task"></task-line>
                    </ul>
                </section>
                <section id="doing">
                    <div class="statusName">Doing</div>
                    <ul class="taskList" data-status="doing">
                        <task-line v-for="task in tasks.doing" :key="task.createDate" :task="task"></task-line>
                    </ul>
                </section>
                <section id="pr">
                    <div class="statusName">PullRequest</div>
                    <ul class="taskList" data-status="pr">
                        <task-line v-for="task in tasks.pullRequest" :key="task.createDate" :task="task"></task-line>
                    </ul>
                </section>
                <section id="done">
                    <div class="statusName">Done</div>
                    <ul class="taskList" data-status="done">
                        <task-line v-for="task in tasks.done" :key="task.createDate" :task="task"></task-line>
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
    import TaskLine from "./TaskBoard/TaskLine";

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

                this.tasks.todo.splice(0, this.tasks.todo.length);
                this.tasks.doing.splice(0, this.tasks.doing.length);
                this.tasks.pullRequest.splice(0, this.tasks.pullRequest.length);
                this.tasks.done.splice(0, this.tasks.done.length);

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
            }
        },
        created() {
            this.loadProjectInfo()
            this.loadTasks()
        }
    }
</script>

<style scoped>

</style>