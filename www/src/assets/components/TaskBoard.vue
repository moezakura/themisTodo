<template>
    <div id="task-board">
        <div class="project-title-container">
            <h2 class="project-title"><i class="fas fa-tasks"></i><span>{{ storeProject.name }}</span></h2>
            <ul class="project-actions">
                <li @mouseenter="isShowOtherMenu = true" @mouseleave="isShowOtherMenu = false">
                    <i class="fas fa-toolbox"></i>OTHER
                    <ul class="project-actions-other" v-show="isShowOtherMenu">
                        <li @click="reloadProject"><i class="fas fa-redo"></i>RELOAD</li>
                        <li @click="moveHideTasks"><i class="fas fa-eye-slash"></i>HIDE TASKS</li>
                    </ul>
                </li>
                <li @click="moveSettings"><i class="fas fa-cog"></i>SETTING</li>
                <li @click="toggleIsShowTaskAdd"><i class="fas fa-plus-circle"></i>ADD</li>
            </ul>
        </div>
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

        <div>
            <task-detail @load-tasks="loadTasks"></task-detail>
            <task-add :class="{ shown: isShowTaskAdd }" v-model="isShowTaskAdd" @load-tasks="loadTasks"></task-add>
            <project-settings v-model="isShowProjectSettings"></project-settings>
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
    import TaskDetail from "./TaskBoard/TaskDetail"
    import TaskAdd from "./TaskBoard/TaskAdd"
    import ProjectSettings from "./TaskBoard/ProjectSettings"

    export default {
        name: "TaskBoard",
        components: {ProjectSettings, TaskAdd, TaskDetail, TaskLine},
        data() {
            const project = new Project()
            const todo: Array<Task> = []
            const doing: Array<Task> = []
            const pullRequest: Array<Task> = []
            const done: Array<Task> = []

            return {
                project: project,
                isShowTaskAdd: false,
                isShowProjectSettings: false,
                isShowOtherMenu: false,
                tasks: {
                    todo: todo,
                    doing: doing,
                    pullRequest: pullRequest,
                    done: done,
                }
            }
        },
        computed: {
            projectId(): number | undefined {
                if (this.$route.params === undefined) {
                    return
                }
                return this.$route.params["projectId"]
            },
            taskId(): number | undefined {
                if (this.$route.params === undefined) {
                    return
                }
                return this.$route.params["taskId"]
            },
            storeProject(): Project {
                if (this.$store.getters.getCurrentProject == undefined) {
                    return new Project()
                }
                return this.$store.getters.getCurrentProject
            }
        },
        watch: {
            '$route'(to, from) {
                this.setCurrentTask()
            }
        },
        methods: {
            async runInit() {
                this.loadProjectInfo()
                await this.loadTasks()

                this.setCurrentTask()
            },
            setCurrentTask() {
                this.isShowProjectSettings = false
                if (this.$route.name == "taskDetail") {
                    if (this.taskId == undefined) {
                        this.$store.commit("setCurrentTask", undefined)
                        return
                    }
                    const selectedTask = this.findTask(this.taskId)
                    this.$store.commit("setCurrentTask", selectedTask)
                } else if (this.$route.meta.isSettings) {
                    this.isShowProjectSettings = true
                }
            },
            reloadProject() {
                this.runInit()
            },
            moveHideTasks() {
                this.$router.push({name: "hiddenTasks", params: {projectId: this.projectId}})
            },
            async loadProjectInfo() {
                this.$store.commit("incrementLoadingCount")
                await ProjectApi.getProject(this.projectId).then(res => {
                    if (res.success) {
                        this.project = res.project
                    }
                    this.$store.commit("setCurrentProject", res.project)
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            },
            async loadTasks() {
                this.$store.commit("incrementLoadingCount")

                this.tasks.todo.splice(0, this.tasks.todo.length)
                this.tasks.doing.splice(0, this.tasks.doing.length)
                this.tasks.pullRequest.splice(0, this.tasks.pullRequest.length)
                this.tasks.done.splice(0, this.tasks.done.length)

                await ProjectApi.getTasks(this.projectId).then(res => {
                    if (res.success) {
                        for (let task of res.task) {
                            task.projectId = this.projectId
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
            },
            findTask(taskId: number): Task | undefined {
                for (const key in this.tasks) {
                    for (const t of this.tasks[key]) {
                        if (t.taskId == taskId) {
                            return <Task>t
                        }
                    }
                }

                return
            },
            toggleIsShowTaskAdd() {
                this.isShowTaskAdd = !this.isShowTaskAdd
            },
            moveSettings() {
                this.$router.push({name: "projectSettings", params: {projectId: this.projectId}})
            }
        },
        created() {
            this.runInit()
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

<style lang="scss" scoped>
    #task-board {
        height: calc(100% - 80px);

        .project-actions-other {
            position: absolute;
            top: 65px + 10px + 45px;
            right: 20px + (140px + 3px * 2) * 2 - 15px - 40px;
            background-color: rgb(30, 30, 30);
            box-shadow: 3px 6px 6px rgba(0, 0, 0, 0.8);
            z-index: 50;

            li {
                padding-left: 15px;
                text-align: left;
                width: 180px;
            }
        }
    }
</style>