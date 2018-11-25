<template>
    <div id="task-board">
        <div class="project-title-container">
            <h2 class="project-title"><i class="fas fa-tasks"></i><span>{{ storeProject.name }}</span></h2>
            <form class="project-task-search" @submit.prevent>
                <input type="search" placeholder="TASK SEARCH" v-model="searchText"/>
            </form>
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
                        <li v-for="task in tasks.todo" :data-task-id="task.createDate">
                            <task-line :key="task.createDate" :task="task"></task-line>
                        </li>
                    </ul>
                </section>
                <section id="doing">
                    <div class="statusName">Doing</div>
                    <ul class="taskList" data-status="DOING" ref="task-line-doing">
                        <li v-for="task in tasks.doing" :data-task-id="task.createDate">
                            <task-line :key="task.createDate" :task="task"></task-line>
                        </li>
                    </ul>
                </section>
                <section id="pr">
                    <div class="statusName">PullRequest</div>
                    <ul class="taskList" data-status="PULL_REQUEST" ref="task-line-pr">
                        <li v-for="task in tasks.pullRequest" :data-task-id="task.createDate">
                            <task-line :key="task.createDate" :task="task"></task-line>
                        </li>
                    </ul>
                </section>
                <section id="done">
                    <div class="statusName">Done</div>
                    <ul class="taskList" data-status="DONE" ref="task-line-done">
                        <li v-for="task in tasks.done" :data-task-id="task.createDate">
                            <task-line :key="task.createDate" :task="task"></task-line>
                        </li>
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
            const todoOrigin: Array<Task> = []
            const doing: Array<Task> = []
            const doingOrigin: Array<Task> = []
            const pullRequest: Array<Task> = []
            const pullRequestOrigin: Array<Task> = []
            const done: Array<Task> = []
            const doneOrigin: Array<Task> = []

            return {
                project: project,
                isShowTaskAdd: false,
                isShowProjectSettings: false,
                isShowOtherMenu: false,
                searchText: "",
                tasks: {
                    todo: todo,
                    doing: doing,
                    pullRequest: pullRequest,
                    done: done,
                },
                tasksOrigin: {
                    todo: todoOrigin,
                    doing: doingOrigin,
                    pullRequest: pullRequestOrigin,
                    done: doneOrigin,
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

                if (this.$route.query["search"]) {
                    this.searchText = this.$route.query["search"]
                }
            },
            searchText(value) {
                for (let key in this.tasksOrigin) {
                    let tasksCopy = this.tasksOrigin[key].slice(0, this.tasksOrigin[key].length)
                    let isNext = true
                    tasksCopy = tasksCopy.filter(task => {
                        let searchStr = value.toLowerCase()

                        if (searchStr.length <= 0) {
                            return true
                        }

                        if (searchStr.startsWith("#")) {
                            const taskId = searchStr.slice(1)
                            if (task.taskId == taskId) {
                                isNext = false
                                return true
                            }
                        }

                        if (searchStr.startsWith("@")) {
                            const userId = searchStr.slice(1)
                            if (task.assignName == userId) {
                                isNext = false
                                return true
                            }
                        }

                        if (!isNext) {
                            return false
                        }
                        searchStr = searchStr.replace("\\", "")

                        return (task.name.toLowerCase().indexOf(searchStr) > -1 ||
                            task.description.toLowerCase().indexOf(searchStr) > -1 ||
                            task.assignName.toLowerCase().indexOf(searchStr) > -1)
                    })

                    this.$router.replace({query: {search: value}})
                    this.$set(this.tasks, key, tasksCopy)
                }
            }
        },
        methods: {
            async runInit() {
                this.loadProjectInfo()
                await this.loadTasks()

                this.setCurrentTask()
                if (this.$route.query["search"]) {
                    this.searchText = this.$route.query["search"]
                }
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

                for (let key in this.tasks) {
                    this.tasks[key].splice(0, this.tasks[key].length)
                }
                for (let key in this.tasksOrigin) {
                    this.tasksOrigin[key].splice(0, this.tasksOrigin[key].length)
                }

                await ProjectApi.getTasks(this.projectId).then(res => {
                    if (res.success) {
                        for (let task of res.task) {
                            task.projectId = this.projectId
                            switch (task.status) {
                                case TaskStatusConvert.toNumber(TaskStatus.TODO):
                                    this.tasks.todo.push(task)
                                    this.tasksOrigin.todo.push(task)
                                    break
                                case TaskStatusConvert.toNumber(TaskStatus.DOING):
                                    this.tasks.doing.push(task)
                                    this.tasksOrigin.doing.push(task)
                                    break
                                case TaskStatusConvert.toNumber(TaskStatus.PULL_REQUEST):
                                    this.tasks.pullRequest.push(task)
                                    this.tasksOrigin.pullRequest.push(task)
                                    break
                                case TaskStatusConvert.toNumber(TaskStatus.DONE):
                                    this.tasks.done.push(task)
                                    this.tasksOrigin.done.push(task)
                                    break
                            }

                            const currentTask: Task | undefined = this.$store.getters.getCurrentTask
                            if (currentTask != undefined && currentTask.createDate == task.createDate) {
                                this.$store.commit("setCurrentTask", task)
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

                TaskApi.update(createDate, task).then(res => {
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