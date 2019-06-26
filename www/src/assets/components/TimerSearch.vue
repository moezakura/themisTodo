<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <div class="task-timer-board">
        <project-header>
            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-menu>
                <li @click="loadPage"><i class="fas fa-redo"></i>RELOAD</li>
            </template>

            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:inner-content>
                <h2 class="task-timer-board-title">Task timer search</h2>
            </template>

            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-right-menu>
                <li @click="moveTasks"><i class="fas fa-tasks"></i>TASKS</li>
            </template>
        </project-header>

        <form class="task-timer-search-box basicForm" @submit.prevent="startTimer">
            <div class="user-select-container">
                <ul class="selected-user-icons">
                    <li v-for="(user, index) in search.user" :key="user.uuid"
                        :style="{'background-image': `url(/api/account/icon/${user.iconPath})`}"
                        v-if="index < 5"></li>
                    <li v-if="search.user.length > 5" class="user-icon-more">...</li>
                </ul>
                <user-multi-select :isInProject="true" v-model="search.user"></user-multi-select>
            </div>
        </form>

        <div class="task-timer-history">
            <div class="task-timer-history-title-section">
                <h3 class="task-timer-history-title"><i class="fas fa-calendar-alt"></i>{{ displayDate.start }} 〜 {{
                    displayDate.end }}</h3>
            </div>

            <ul class="task-timer-entry-container">
                <li is="task-timer-line" v-for="i in timeHistories" :class="{ active: i.endDateUnix === 0 }"
                    :task-timer="i" :key="i.id" @load-page="loadPage"></li>
            </ul>

            <div class="total-task-timer">
                <div class="total-task-timer-title">Total</div>
                <div class="total-task-timer-time">{{ totalTimeHM }}</div>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
    import Project from "@scripts/model/api/project/Project";
    import ProjectApi from "@scripts/api/ProjectApi";
    import ProjectHeader from "@components/Project/ProjectHeader.vue";
    import TaskTimerApi from "@scripts/api/TaskTimer";
    import TaskTimerGetMyListRequest from "@scripts/model/api/taskTimer/TaskTimerGetMyListRequest";
    import Task from "@scripts/model/api/task/Task";
    import TaskStatusConvert, {TaskStatus} from "@scripts/enums/TaskStatus";
    import TaskTimerLine from "@components/TaskTimer/TaskTimerLine.vue"
    import TaskLine from "@components/TaskBoard/TaskLine.vue"
    import TaskTimer from "@scripts/model/api/taskTimer/TaskTimer";
    import User from "@scripts/model/api/user/User";
    import UserMultiSelect from "@components/Common/UserMultiSelect.vue";

    interface TimerBoardData {
        displayDate: {
            start: string,
            end: string,
        },
        tasks: Array<Task>,
        search: {
            user: Array<User>
        },
        taskTimerTopFocus: number,
        timeHistories: Array<TaskTimer>,
        taskReloadTimer: undefined | number,
        totalTimeHM: string,
        startDate: Date,
        endDate: Date,
    }

    export default {
        name: "TimerSearch",
        components: {UserMultiSelect, TaskLine, TaskTimerLine, ProjectHeader},
        data(): TimerBoardData {
            return {
                displayDate: {
                    start: "",
                    end: "",
                },
                tasks: [],
                search: {
                    user: new Array<User>()
                },
                taskTimerTopFocus: 0,
                timeHistories: [],
                taskReloadTimer: undefined,
                totalTimeHM: "",
                startDate: new Date(),
                endDate: new Date(),
            }
        },
        computed: {
            projectId(): number | undefined {
                if (this.$route.params === undefined) {
                    return
                }
                return this.$route.params["projectId"]
            },
            storeProject(): Project {
                if (this.$store.getters.getCurrentProject == undefined) {
                    return new Project()
                }
                return this.$store.getters.getCurrentProject
            },
            searchedTasks(): Array<Task> {
                const searchText = this.search.text
                const excludeDone = this.search.excludeDone

                let taskIdSearch: undefined | number = undefined
                if (searchText.startsWith("#")) {
                    let taskIdSearchString = searchText.slice(1)
                    taskIdSearch = Number(taskIdSearchString)
                    if (isNaN(taskIdSearch)) {
                        taskIdSearch = undefined
                    }
                }

                let searchedTasks = new Array<Task>()
                this.tasks.forEach((task: Task) => {
                    if (excludeDone && task.status === TaskStatusConvert.toNumber(TaskStatus.DONE)) {
                        return
                    }

                    if (typeof taskIdSearch !== "undefined" && task.taskId === taskIdSearch) {
                        searchedTasks.push(task)
                        return
                    }

                    if (task.name.includes(searchText) || task.description.includes(searchText)) {
                        searchedTasks.push(task)
                        return
                    }
                })

                return searchedTasks
            },
        },
        watch: {
            startDate(startDate: Date): void {
                const startDateString = startDate.getFullYear() + "/" +
                    ("00" + (startDate.getMonth() + 1)).slice(-2) + "/" +
                    ("00" + startDate.getDate()).slice(-2) + " " +
                    ("00" + startDate.getHours()).slice(-2) + ":" +
                    ("00" + startDate.getMinutes()).slice(-2) + ":" +
                    ("00" + startDate.getSeconds()).slice(-2)
                this.$set(this.displayDate, "start", startDateString)
                this.loadPage()
            },
            endDate(endDate: Date): void {
                const endDateString = endDate.getFullYear() + "/" +
                    ("00" + (endDate.getMonth() + 1)).slice(-2) + "/" +
                    ("00" + endDate.getDate()).slice(-2) + " " +
                    ("00" + endDate.getHours()).slice(-2) + ":" +
                    ("00" + endDate.getMinutes()).slice(-2) + ":" +
                    ("00" + endDate.getSeconds()).slice(-2)
                this.$set(this.displayDate, "end", endDateString)
                this.loadPage()
            },
            taskTimerTopFocus(value: number): void {
                if (value < 0) {
                    this.taskTimerTopFocus = 0
                }

                if (value > 0) {
                    const el = this.$refs["task-timer-entry-name"]
                    const rect = el.getBoundingClientRect()
                    const pos = {
                        top: (rect.bottom + 15) + "px",
                        left: (rect.left - 15) + "px",
                        width: (rect.right - rect.left) + "px",
                        display: "block"
                    }
                    this.$set(this.search, "position", pos)
                } else {
                    this.$set(this.search, "position", {
                        display: "none"
                    })
                }
            },
            '$route'(): void {
                this.parseQuery()
            }
        },
        methods: {
            moveTasks(): void {
                this.$router.push({name: "taskBoard", params: {projectId: this.projectId}})
            },
            parseQuery(): void {
                const now = new Date()
                let startUnix: undefined | string = this.$route.query.start
                let endUnix: undefined | string = this.$route.query.end

                let startDate: undefined | Date = undefined
                let endDate: undefined | Date = undefined

                try {
                    startDate = this.parseStartDate(startUnix)
                } catch (e) {
                    startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0)
                }
                this.startDate = startDate


                try {
                    endDate = this.parseStartDate(endUnix)
                } catch (e) {
                    endDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59)
                }
                this.endDate = endDate
            },
            parseStartDate(startUnix: string): Date {
                let startDate: undefined | Date = undefined

                if (typeof startUnix !== "undefined") {
                    let startUnixNumber = Number(startUnix)
                    if (!isNaN(startUnixNumber)) {
                        startDate = new Date(startUnixNumber)
                    }
                }

                if (typeof startDate === "undefined" || startDate === null) {
                    throw "undefined date"
                }

                if (startDate.toString() === "Invalid Date") {
                    throw "Invalid date"
                }

                return startDate
            },
            parseEndDate(endUnix: string): Date {
                let endDate: undefined | Date = undefined

                if (typeof endUnix !== "undefined") {
                    let startUnixNumber = Number(endUnix)
                    if (!isNaN(startUnixNumber)) {
                        endDate = new Date(startUnixNumber)
                    }
                }

                if (typeof endDate === "undefined" || endDate === null) {
                    throw "undefined date"
                }

                if (endDate.toString() === "Invalid Date") {
                    throw "Invalid date"
                }

                return endDate
            },
            addEntryFocus(): void {
                this.$refs['task-timer-entry-name'].focus()
            },
            loadPage(isLoadingShow: boolean = true, startDate?: Date, endDate?: Date): void {
                if (isLoadingShow) {
                    this.$store.commit("incrementLoadingCount")

                    // Project info
                    ProjectApi.getProject(this.projectId).then(res => {
                        this.$store.commit("setCurrentProject", res.project)
                    }).finally(() => {
                        this.$store.commit("decrementLoadingCount")
                    })
                }

                if (isLoadingShow) {
                    this.$store.commit("incrementLoadingCount")
                }

                let req = new TaskTimerGetMyListRequest()
                if (typeof startDate === "undefined") {
                    req.startDate = this.startDate
                } else {
                    req.startDate = startDate
                }
                if (typeof endDate === "undefined") {
                    req.endDate = this.endDate
                } else {
                    req.endDate = endDate
                }
                TaskTimerApi.getMyList(this.projectId, req).then(res => {
                    this.timeHistories = res.list
                    let totalTimeSec = 0
                    let startFlag = false
                    for (const item of res.list) {
                        if (item.endDateUnix === 0) {
                            totalTimeSec += new Date().getTime() / 1000 - item.startDateUnix
                            startFlag = true
                        } else {
                            totalTimeSec += item.endDateUnix - item.startDateUnix
                        }
                    }
                    if (startFlag) {
                        this.taskReloadTimer = setTimeout(() => {
                            this.loadPage(false)
                        }, 10 * 1000)
                    }
                    let totalTimeH = Math.floor(totalTimeSec / 3600)
                    let totalTimeM = Math.floor(totalTimeSec / 60 - totalTimeH * 60)
                    this.totalTimeHM = ("00" + totalTimeH).slice(-2) + ":" + ("00" + totalTimeM).slice(-2)
                }).finally(() => {
                    if (isLoadingShow) {
                        this.$store.commit("decrementLoadingCount")
                    }
                })
            },
            selectedSearchTask(task: Task): void {
                this.$set(this.search, "text", task.name)
                this.$set(this.search, "selectedTask", task)
                this.taskTimerTopFocus--
            },
            async startTimer(): Promise<void> {
                if (typeof this.search.selectedTask === "undefined" || this.$store.getters.isLoadingShow) {
                    return
                }

                const task: Task = this.search.selectedTask

                if (task.name !== this.search.text) {
                    return
                }

                this.$store.commit("incrementLoadingCount")
                const res = await TaskTimerApi.getTaskTimerStatus(task.createDate)

                if (!res.start) {
                    const toggleRes = await TaskTimerApi.toggleTimer(task.createDate)
                    if (toggleRes.success) {
                        this.$set(this.search, "selectedTask", undefined)
                        this.$set(this.search, "text", "")
                        this.loadPage()
                    }
                }
                this.$store.commit("decrementLoadingCount")
            },
        },
        created(): void {
            this.parseQuery()
            this.loadPage(true)
        },
        destroyed(): void {
            if (typeof this.taskReloadTimer !== "undefined") {
                clearTimeout(this.taskReloadTimer)
            }
        }
    }
</script>

<style scoped lang="scss">
    .task-timer-board {
        height: calc(100% - #{$headerHeight + 10px});

        .task-timer-board-title {
            margin: 0 auto;
            text-align: center;
            line-height: $buttonHeight;
            font-size: 16px;
            font-weight: 900;
            letter-spacing: 2px;
        }
    }

    .task-timer-search-box {
        width: 85%;
        height: 105px;
        margin: 10px auto;

        .user-select-container {
            width: 200px;

            .selected-user-icons {
                display: flex;
                height: 30px;

                li {
                    width: 30px;
                    height: 30px;
                    border-radius: 50%;
                    background-size: cover;
                    background-repeat: no-repeat;
                    background-position: center;
                    margin: 0 2px;

                    &.user-icon-more {
                        width: 20px;
                    }
                }
            }
        }
    }

    .searched-task-list {
        position: absolute;
        z-index: 20;
        background-color: $backgroundColor;
        display: none;
        border: solid 2px $accentColor;
        border-top: 0;
        max-height: 400px;
        overflow-y: auto;
    }

    .task-timer-history {
        width: 85%;
        margin: 15px auto;
        height: calc(100% - #{45px + 110px + (15px * 2) + 10px});

        .task-timer-history-title-section {
            $height: 40px;
            display: flex;
            height: $height;
            line-height: $height;
            font-size: 16px;
            margin-top: 10px;

            .task-timer-history-title {
                font-weight: bold;

                i {
                    margin-right: 8px;
                }
            }

            .task-timer-history-title-detail {
                font-weight: bold;
                margin-left: 15px;
                opacity: 0.5;
            }

            .task-timer-history-title-action {
                margin-left: auto;

                i {
                    width: $height;
                    height: $height;
                    line-height: $height;
                    text-align: center;
                }
            }
        }

        .task-timer-entry-container {
            margin-top: 8px;
            height: calc(100% - #{40px + 35px});
            overflow-y: auto;
        }

        .total-task-timer {
            display: flex;
            font-size: 16px;
            letter-spacing: 2px;
            height: 35px;
            line-height: 35px;

            .total-task-timer-title {
                margin: 0 15px 0 auto;
            }
        }
    }
</style>