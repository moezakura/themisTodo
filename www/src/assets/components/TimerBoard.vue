<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <div class="task-timer-board">
        <project-header>
            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-menu>
                <li @click="loadPage"><i class="fas fa-redo"></i>RELOAD</li>
                <li @click="loadPage"><i class="fas fa-user-clock"></i>OTHER TIMERS</li>
            </template>

            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:inner-content>
                <h2 class="task-timer-board-title">Your task timer list</h2>
            </template>

            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-right-menu>
                <li @click="addEntryFocus"><i class="fas fa-plus-circle"></i>ADD TIMER</li>
            </template>
        </project-header>

        <form class="task-timer-add-container basicForm" @submit.prevent="startTimer">
            <div class="task-timer-add" :class="{'active': taskTimerTopFocus > 0}">
                <label class="task-timer-entry-name">
                    <input placeholder="What did you task on?" @focus="taskTimerTopFocus++"
                           @blur="taskTimerTopFocus--" ref="task-timer-entry-name" v-model="search.text">
                </label>
                <input class="task-timer-entry-submit" type="submit" value="Add timer">
            </div>
            <div class="task-timer-exclude-done-task-container">
                <input id="exclude-done-task" type="checkbox" v-model="search.excludeDone">
                <label for="exclude-done-task" class="fas fa-check"></label>
                <label for="exclude-done-task">exclude done task</label>
            </div>
        </form>
        <ul class="taskList searched-task-list" :style="search.position" @mouseenter="taskTimerTopFocus++"
            @mouseleave="taskTimerTopFocus--">
            <li v-for="task in searchedTasks" :data-task-id="task.createDate" @click="selectedSearchTask(task)">
                <task-line :task="task" :allowShowDetail="false"></task-line>
            </li>
        </ul>

        <div class="task-timer-history">
            <div class="task-timer-history-title-section">
                <h3 class="task-timer-history-title"><i class="fas fa-calendar-alt"></i>{{ displayDate.string }}</h3>
                <div class="task-timer-history-title-detail" v-if="isStartToday">{{ displayDate.stringMD }}</div>
                <div class="task-timer-history-title-action">
                    <i class="fas fa-chevron-left" @click="moveDisplayDate(-1)"></i>
                    <i class="fas fa-chevron-right" @click="moveDisplayDate(1)"></i>
                </div>
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

    interface TimerBoardData {
        displayDate: {
            stringMD: string,
            string: string,
        },
        tasks: Array<Task>,
        search: {
            text: string,
            selectedTask: undefined | Task
            excludeDone: boolean,
            position: object
        },
        taskTimerTopFocus: number,
        timeHistories: Array<TaskTimer>,
        taskReloadTimer: undefined | number,
        totalTimeHM: string,
        startDate: Date,
        endDate: Date,
    }

    export default {
        name: "TimerBoard",
        components: {TaskLine, TaskTimerLine, ProjectHeader},
        data(): TimerBoardData {
            return {
                displayDate: {
                    stringMD: "",
                    string: "",
                },
                tasks: [],
                search: {
                    text: "",
                    selectedTask: undefined,
                    excludeDone: true,
                    position: {},
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
            isStartToday(): boolean {
                const startDate: Date = this.startDate
                const now: Date = new Date()
                const startDateString = startDate.getFullYear() + "/" + ("0" + (startDate.getMonth() + 1)).slice(-2) + "/" + ("0" + startDate.getDate()).slice(-2)
                const nowDateString = now.getFullYear() + "/" + ("0" + (now.getMonth() + 1)).slice(-2) + "/" + ("0" + now.getDate()).slice(-2)

                if (startDateString == nowDateString) {
                    this.$set(this.displayDate, 'string', 'Today')
                    this.$set(this.displayDate, 'stringMD', startDateString)
                    return true
                } else {
                    this.$set(this.displayDate, 'string', startDateString)
                    return false
                }
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
            startDate(): void {
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
            loadTasks(): void {
                this.$store.commit("incrementLoadingCount")
                ProjectApi.getTasks(this.projectId).then(res => {
                    this.tasks = res.task
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
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

                    // Tasks
                    this.loadTasks()
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
            moveDisplayDate(diff: number): void {
                const startDate: Date = this.startDate
                const endDate: Date = this.endDate
                startDate.setDate(startDate.getDate() + diff)
                endDate.setDate(endDate.getDate() + diff)

                this.$router.push({
                    name: "timerBoard",
                    query: {
                        start: startDate.getTime().toString(),
                        end: endDate.getTime().toString(),
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

    .task-timer-add-container {
        width: 75%;
        height: 110px;
        margin: 10px auto;

        .task-timer-add {
            display: flex;
            height: $buttonHeight + 30px;
            padding: 0 15px;
            background-color: rgba(black, .4);
            box-shadow: 2px 1px 3px rgba(black, 0.5);
            border-bottom: solid 2px transparent;
            -webkit-transition: border-bottom-color .3s ease;
            -moz-transition: border-bottom-color .3s ease;
            -o-transition: border-bottom-color .3s ease;
            transition: border-bottom-color .3s ease;
            $task-timer-entry-submit-width: 200px;

            &.active {
                border-bottom-color: $accentColor;
            }

            .task-timer-entry-name {
                width: calc(100% - #{$task-timer-entry-submit-width + 15px});
                margin-right: 15px;
                height: $buttonHeight;

                input {
                    display: block;
                    width: 100%;
                    height: $buttonHeight;
                    background-color: transparent;
                    border: 0;
                    letter-spacing: 1.5px;
                }
            }

            .task-timer-entry-submit {
                width: $task-timer-entry-submit-width;
                height: $buttonHeight;
            }
        }

        .task-timer-exclude-done-task-container {
            $height: 20px;
            display: flex;
            height: $height;
            line-height: $height;
            margin: 10px 0 0 0;
            -webkit-user-select: none;
            -moz-user-select: none;
            -ms-user-select: none;
            user-select: none;

            #exclude-done-task {
                display: none;

                & + label {
                    display: block;
                    width: $height;
                    height: $height;
                    box-sizing: border-box;
                    border: solid thin rgba(white, .5);
                    margin: 0 8px 0 auto;

                    &::before {
                        font-size: 14px;
                        display: block;
                        text-align: center;
                        width: $height;
                        height: $height;
                        line-height: $height;
                        opacity: 0;
                        -webkit-transition: opacity ease .3s;
                        -moz-transition: opacity ease .3s;
                        -ms-transition: opacity ease .3s;
                        -o-transition: opacity ease .3s;
                        transition: opacity ease .3s;
                    }
                }

                &:checked + label::before {
                    opacity: 1;
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