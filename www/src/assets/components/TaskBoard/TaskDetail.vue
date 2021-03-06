<template>
    <div>
        <transition>
            <div v-if="isShowTaskDetail">
                <form id="taskPopup" ref="taskPopup" @submit.prevent="submitEdit">
                    <h2 :class="[limitAddClass]">{{ taskStatusText }} Task</h2>
                    <div id="taskPopupActions">
                        <i class="fas fa-history" id="taskPopupHistoryButton" @click="showToggleTaskHistoryList"></i>
                        <i class="fas fa-trash" id="taskPopupTrashButton" @click="showTaskDeleteOrHide"></i>
                        <i class="fas fa-edit" id="taskPopupEditButton" @click="setEditing(!isEditing)"></i>
                        <i class="fas fa-times" id="taskPopupCloseButton" @click="hideTaskDetail"></i>
                    </div>
                    <div class="both"></div>
                    <div class="success" v-show="isSuccess">Update Success.</div>
                    <div class="error" v-show="errorMessage !== undefined && errorMessage.length > 0">{{ errorMessage }}
                    </div>
                    <div id="taskPopupTaskIdTitle">
                        <label id="taskPopupTaskId">#{{ task.taskId }}</label>
                        <input value="TITLE" id="taskPopupTitle" v-model="taskCache.name" :readonly="!isEditing">
                    </div>
                    <div class="both"></div>
                    <div id="taskPopupAssignCreatorLine">
                        <div class="taskPopupAssignCreatorColumn">
                            <p>Assign</p>
                            <div id="taskPopupAssignIcon"
                                 :style="{ 'background-image': `url('${assignIconPath}')` }"></div>
                            <user-select :is-show="true" :is-in-project="true" v-model="selectUser"
                                         :readonly="!isEditing"></user-select>
                        </div>
                        <div class="taskPopupAssignCreatorColumn">
                            <p>Creator</p>
                            <div id="taskPopupCreatorIcon"
                                 :style="{ 'background-image': `url('${creatorIconPath}')` }"></div>
                            <label id="taskPopupCreator">{{ task.creatorName }}</label>
                        </div>
                    </div>
                    <div class="both"></div>
                    <div class="timer-button">
                        <task-timer-simple-controller></task-timer-simple-controller>
                    </div>
                    <div id="taskPopupProgressBar">
                        <div id="taskPopupProgressText" v-show="!isEditing">
                            <i class="fas fa-calendar-alt"></i>
                            <span>{{ task.deadline }}</span>
                            <span v-if="!isCompleted">(あと{{ task.limitDate }}日)</span>
                            <span v-if="isCompleted">(Already Completed!)</span>
                        </div>
                        <input type="date" id="taskPopupDeadlineChange" v-show="isEditing" v-model="taskCache.deadline">
                        <div id="taskPopupProgressCurrent" :style="{'width': `${currentProgress}%`}"
                             :class="[limitAddClass]">&nbsp;
                        </div>
                    </div>
                    <div>
                        <task-detail-description v-model="taskCache.description"
                                                 :readonly="!isEditing" v-if="isEditing"></task-detail-description>
                        <task-detail-description-rich :task="task" v-if="!isEditing"
                                                      class="task-detail-description-rich"
                                                      @detailUpdate="detailUpdate"></task-detail-description-rich>
                    </div>
                    <div class="input-box" v-show="isEditing" ref="changeActionsBox">
                        <input type="button" value="CANCEL" @click="setEditing(false)">
                        <input type="submit" value="CHANGE">
                    </div>
                </form>

                <div class="backView" @click="hideTaskDetail"></div>
            </div>
        </transition>
        <transition name="slide">
            <div id="taskHistoryList" v-if="isShowTaskDetail && taskHistoryList.isShow"
                 :style="{top: `calc(${taskHistoryList.topPos})`, height: `calc(${taskHistoryList.height})`}">
                <div class="title-bar">
                    <p>Task History</p>
                    <i class="fas fa-chevron-left" @click="hideTaskHistoryList"></i>
                </div>
                <ul :style="{height: `calc(${taskHistoryList.height} - (45px * 2 + 5px + 10px))` }">
                    <li v-for="i in taskHistoryList.list" :class="{ selected: currentTaskUpdateDate === i.updateDate }"
                        @click="selectHistory(i.updateDate)">
                        {{ i.updateDateFormat }}
                    </li>
                </ul>
                <div class="apply-button" @click="applyClick">APPLY</div>
            </div>
        </transition>
        <task-delete-or-hide @reload="commitLoadTasks"></task-delete-or-hide>
        <task-delete-confirm @reload="commitLoadTasks"></task-delete-confirm>
    </div>
</template>

<script lang="ts">
    import Task from "../../scripts/model/api/task/Task"
    import UserSelect from "../Common/UserSelect"
    import User from "../../scripts/model/api/user/User"
    import TaskApi from "../../scripts/api/TaskApi"
    import TaskDeleteOrHide from "./TaskDeleteOrHide"
    import {ProjectDetailStatus} from "../../scripts/enums/ProjectDetailStatus"
    import TaskDeleteConfirm from "./TaskDeleteConfirm"
    import TaskDetailDescription from "./TaskDetailDescription"
    import TaskDetailDescriptionRich from "./TaskDetailDescriptionRich.ts"
    import TaskHistory from "@scripts/model/api/task/TaskHistory"
    import TaskTimerSimpleController from "@components/TaskTimer/TaskTimerSimpleController.vue"

    export default {
        name: "TaskDetail",
        components: {
            TaskTimerSimpleController,
            TaskDetailDescription,
            TaskDeleteConfirm,
            TaskDeleteOrHide,
            UserSelect,
            TaskDetailDescriptionRich
        },
        data: () => {
            const taskCache: Task | undefined = undefined
            const taskHistoryList: Array<TaskHistory> | undefined = undefined
            const taskHistoryTask: TaskHistory | undefined = undefined

            return {
                isSuccess: false,
                errorMessage: "",
                isEditing: false,
                taskCache: taskCache,
                taskHistoryList: {
                    isShow: false,
                    topPos: "0px",
                    height: "0px",
                    selected: "",
                    originTask: taskCache,
                    selectedTask: taskHistoryTask,
                    list: taskHistoryList
                }
            }
        },
        computed: {
            isCompleted(): boolean {
                if (this.task === undefined || this.task.status === undefined || this.task.status === null) {
                    return false
                }
                return this.task.status === 3
            },
            task(): Task | undefined {
                const t = this.$store.getters.getCurrentTask
                this.taskCache = Object.assign({}, t)
                return this.taskCache
            },
            isShowTaskDetail(): boolean {
                return this.task != undefined && this.task.taskId != undefined
            },
            assignIconPath(): string {
                return `/api/account/icon/${this.task.assignIconPath}`
            },
            creatorIconPath(): string {
                return `/api/account/icon/${this.task.creatorIconPath}`
            },
            selectUser: {
                get(): User {
                    const user = new User()
                    user.uuid = this.task.assign
                    user.displayName = this.task.assignName
                    return user
                },
                set(value) {
                    this.task.assign = value.uuid
                    this.task.assignName = value.displayName
                }
            },
            currentProgress(): number {
                let createDate = Math.round(this.task.createDate / 1000000),
                    limitDate = (new Date(this.task.deadline)).getTime(),
                    allDiff = Math.abs(limitDate - createDate),
                    limit = Math.abs(limitDate - (new Date()).getTime()),
                    progress = limit / allDiff * 100

                if (progress > 100) return 100
                else if (progress < 0) return 0
                else return progress
            },
            limitAddClass(): string {
                if (this.task.status == 3) {
                    return "normal"
                }
                if (this.task.limitDate <= 0) {
                    return "over"
                } else if (this.task.limitDate <= 1) {
                    return "limit1"
                } else if (this.task.limitDate <= 2) {
                    return "limit2"
                } else if (this.task.limitDate <= 3) {
                    return "limit3"
                }

                return "normal"
            },
            taskStatusText(): string {
                switch (this.task.status) {
                    case 0:
                        return "Todo"
                    case 1:
                        return "Doing"
                    case 2:
                        return "PullRequest"
                    case 3:
                        return "Done"
                }
                return "Unknown"
            },
            currentTaskCreateDate(): string {
                return this.$store.getters.getCurrentTask.createDate
            },
            currentTaskUpdateDate(): string {
                if (this.taskHistoryList.selected !== undefined && this.taskHistoryList.selected.length > 0) {
                    return this.taskHistoryList.selected
                }
                return this.$store.getters.getCurrentTask.updateDate
            },
            isShowConfirmDelete(): boolean {
                return this.$store.getters.getProjectDetailStatus == ProjectDetailStatus.DELETE_CONFIRM
            },
        },
        watch: {
            "taskHistoryList.isShow"(value): void {
                if (value) {
                    this.$store.commit("incrementLoadingCount")
                    TaskApi.getHistory(this.currentTaskCreateDate).then((list) => {
                        if (list === undefined) {
                            this.$store.commit("decrementLoadingCount")
                            return
                        }
                        for (const i in list) {
                            const l = list[i]
                            let date = new Date(Number(l.updateDate.slice(0, 13)))
                            list[i].updateDateFormat = date.getFullYear() + "/" +
                                ("00" + date.getMonth()).slice(-2) + "/" +
                                ("00" + date.getDate()).slice(-2) + " " +
                                ("00" + date.getHours()).slice(-2) + ":" +
                                ("00" + date.getMinutes()).slice(-2) + ":" +
                                ("00" + date.getSeconds()).slice(-2)
                        }
                        this.$set(this.taskHistoryList, "list", list)
                        this.$store.commit("decrementLoadingCount")
                    })
                }
            }
        },
        methods: {
            hideTaskDetail() {
                this.hideTaskHistoryList()
                let projectId = 0
                if (typeof this.task !== "undefined") {
                    projectId = this.task.projectId
                }
                this.$router.push({name: 'taskBoard', params: {projectId: projectId}});
                this.$store.commit("setCurrentTask", undefined)
                this.setEditing(false)
            },
            setEditing(isEditing: boolean) {
                if (isEditing) {
                    this.hideTaskHistoryList()
                }
                this.isEditing = isEditing
                if (!this.isEditing) {
                    const t = this.$store.getters.getCurrentTask
                    this.taskCache = Object.assign({}, t)
                }
            },
            showTaskDeleteOrHide() {
                this.hideTaskHistoryList()
                this.setEditing(false)
                this.$store.commit("setProjectDetailStatus", ProjectDetailStatus.DELETE_OR_HIDE)
            },
            commitLoadTasks() {
                this.$emit("load-tasks")
            },
            dateFormat(dateStr: string): string {
                const date = new Date(dateStr)
                const y = date.getFullYear()
                const m = date.getMonth() + 1
                const d = date.getDate()

                let str_m: string = m.toString()
                let str_d: string = d.toString()
                if (m < 10) {
                    str_m = '0' + m
                }
                if (d < 10) {
                    str_d = '0' + d
                }

                return y + '-' + str_m + '-' + str_d
            },
            detailUpdate(detailText) {
                this.$set(this.taskCache, 'description', detailText)
                this.submitEdit()
            },
            submitEdit() {
                this.hideTaskHistoryList()
                let task = new Task()
                task.assign = this.selectUser.uuid
                task.deadline = this.dateFormat(this.taskCache.deadline)
                task.description = this.taskCache.description
                task.name = this.taskCache.name
                task.status = -1
                this.$store.commit("incrementLoadingCount")

                TaskApi.update(this.taskCache.createDate, task).then(res => {
                    if (res.success) {
                        this.$store.commit("setCurrentTask", this.taskCache)
                        this.$emit("load-tasks")
                        this.setEditing(false)
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            },
            showToggleTaskHistoryList(): void {
                this.$set(this.taskHistoryList, 'isShow', !this.taskHistoryList.isShow)
                if (!this.taskHistoryList.isShow) {
                    this.$set(this.taskHistoryList, 'selected', "")
                    this.$set(this.taskHistoryList, 'selectedTask', undefined)
                    if (this.taskHistoryList.originTask !== undefined) {
                        this.$store.commit("setCurrentTask", this.taskHistoryList.originTask)
                    }
                    this.$set(this.taskHistoryList, 'originTask', undefined)
                }

                let height = this.$refs["taskPopup"].offsetHeight
                if (this.isEditing) {
                    this.setEditing(false)
                    height -= this.$refs["changeActionsBox"].offsetHeight
                }
                this.$set(this.taskHistoryList, 'topPos', `50% - ${height}px / 2`)
                this.$set(this.taskHistoryList, 'height', `${height}px`)
            },
            hideTaskHistoryList(): void {
                this.$set(this.taskHistoryList, 'selected', "")
                this.$set(this.taskHistoryList, 'isShow', false)
                this.$set(this.taskHistoryList, 'selectedTask', undefined)
                if (this.taskHistoryList.originTask !== undefined) {
                    this.$store.commit("setCurrentTask", this.taskHistoryList.originTask)
                }
                this.$set(this.taskHistoryList, 'originTask', undefined)
            },
            selectHistory(value: string): void {
                this.$set(this.taskHistoryList, 'selected', value)
                let task: TaskHistory | undefined = undefined
                for (const t of this.taskHistoryList.list) {
                    if (value === t.updateDate) {
                        task = t
                    }
                }

                if (task === undefined) {
                    return
                }

                if (this.taskHistoryList.originTask === undefined) {
                    this.$set(this.taskHistoryList, 'originTask', this.task)
                }
                this.$set(this.taskHistoryList, "selectedTask", task)

                const t = task.task
                let currentTask = this.$store.getters.getCurrentTask
                currentTask.name = t.name
                currentTask.status = t.status
                currentTask.assign = t.assign
                currentTask.assignName = t.assignName
                currentTask.assignIconPath = t.assignIconPath
                currentTask.deadline = t.deadline
                currentTask.limitDate = t.limitDate
                currentTask.deadlineMD = t.deadlineMD
                currentTask.description = t.description
                currentTask.createDate = t.createDate
                this.$store.commit("setCurrentTask", currentTask)
            },
            applyClick(): void {
                this.$store.commit("incrementLoadingCount")
                TaskApi.applyHistory(this.currentTaskCreateDate, this.currentTaskUpdateDate).then((res) => {
                    if (res.success) {
                        this.$store.commit("setCurrentTask", this.taskCache)
                        this.$set(this.taskHistoryList, 'originTask', undefined)
                        this.$emit("load-tasks")
                    }
                    this.$store.commit("decrementLoadingCount")
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .v-enter,
    .v-leave-to {
        opacity: 0;
    }

    .v-enter-active,
    .v-leave-active {
        transition: opacity .2s;
    }

    .slide-enter,
    .slide-leave-to {
        opacity: 0;
        transform: translateX(-30px);
    }

    .slide-enter-active,
    .slide-leave-active {
        transition: opacity .2s, transform ease .2s;
    }
</style>