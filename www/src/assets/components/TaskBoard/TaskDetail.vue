<template>
    <div>
        <transition>
            <div v-if="isShowTaskDetail">
                <form id="taskPopup" @submit.prevent="submitEdit">
                    <h2 :class="[limitAddClass]">Doing Task</h2>
                    <div id="taskPopupActions">
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
                    <textarea id="taskPopupDescription" v-model="taskCache.description"
                              :readonly="!isEditing"></textarea>
                    <div class="input-box" v-show="isEditing">
                        <input type="button" value="CANCEL" @click="setEditing(false)">
                        <input type="submit" value="CHANGE">
                    </div>
                </form>

                <div class="backView" @click="hideTaskDetail"></div>
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

    export default {
        name: "TaskDetail",
        components: {TaskDeleteConfirm, TaskDeleteOrHide, UserSelect},
        data: () => {
            const taskCache: Task | undefined = undefined
            return {
                isSuccess: false,
                errorMessage: "",
                isEditing: false,
                taskCache: taskCache
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
                }
                else if (this.task.limitDate <= 1) {
                    return "limit1"
                }
                else if (this.task.limitDate <= 2) {
                    return "limit2"
                }
                else if (this.task.limitDate <= 3) {
                    return "limit3"
                }

                return "normal"
            },
            isShowConfirmDelete(): boolean {
                return this.$store.getters.getProjectDetailStatus == ProjectDetailStatus.DELETE_CONFIRM
            },
        },
        methods: {
            hideTaskDetail() {
                this.$router.back()
                this.$store.commit("setCurrentTask", undefined)
                this.setEditing(false)
            },
            setEditing(isEditing: boolean) {
                this.isEditing = isEditing
                if (!this.isEditing) {
                    const t = this.$store.getters.getCurrentTask
                    this.taskCache = Object.assign({}, t)
                }
            },
            showTaskDeleteOrHide() {
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
            submitEdit() {
                let task = new Task()
                task.assign = this.selectUser.uuid
                task.deadline = this.dateFormat(this.taskCache.deadline)
                task.description = this.taskCache.description
                task.name = this.taskCache.name
                task.status = -1
                this.$store.commit("incrementLoadingCount")

                TaskApi.Update(this.taskCache.createDate, task).then(res => {
                    if (res.success) {
                        this.$store.commit("setCurrentTask", this.taskCache)
                        this.$emit("load-tasks")
                        this.setEditing(false)
                    }
                }).finally(() => {
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
</style>