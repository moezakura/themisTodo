<template>
    <transition>
        <div v-show="isShow" class="task-delete-popup-container">
            <div id="taskDeletePopup">
                <div id="taskDeletePopupTitle">Delete or Hide confirm</div>
                <div id="taskDeletePopupActions">
                    <i class="fas fa-times" id="taskDeletePopupClose" @click="clickClose"></i>
                </div>
                <p class="taskDeletePopupText" id="taskDeletePopupTextCaution">Would you like to hide or delete
                    ${taskName}?</p>
                <p class="taskDeletePopupText"><strong>HIDE</strong>: Hide the task from the task board.</p>
                <p class="taskDeletePopupText"><strong>DELETE</strong>: Delete the task from the task board.<br>This
                    operation
                    can not be undone.</p>
                <div id="taskDeletePopupSelectButtons">
                    <div id="taskDeletePopupHide" class="deleteButton" @click="clickHide">HIDE</div>
                    <div id="taskDeletePopupDelete" class="deleteButton" @click="clickDelete">DELETE</div>
                </div>
                <div class="both"></div>
            </div>
            <div class="backView" @click="clickClose"></div>
        </div>
    </transition>
</template>

<script lang="ts">
    import {ProjectDetailStatus} from "@scripts/enums/ProjectDetailStatus"
    import TaskApi from "../../scripts/api/TaskApi"
    import Task from "../../scripts/model/api/task/Task"
    import TaskStatusConvert, {TaskStatus} from "../../scripts/enums/TaskStatus"

    export default {
        name: "TaskDeleteOrHide",
        computed: {
            isShow(): boolean {
                return this.$store.getters.getProjectDetailStatus == ProjectDetailStatus.DELETE_OR_HIDE
            },
            task(): Task {
                return this.$store.getters.getCurrentTask
            },
            projectId() {
                return this.$store.getters.getCurrentProject.uuid
            }
        },
        methods: {
            clickClose() {
                this.$store.commit("setProjectDetailStatus", ProjectDetailStatus.HIDE)
            },
            clickHide() {
                let task = this.task
                task.status = TaskStatusConvert.toNumber(TaskStatus.HIDE)
                this.$store.commit("incrementLoadingCount")
                TaskApi.Update(task.createDate, task).then(res => {
                    if (res.success) {
                        this.$store.commit("setProjectDetailStatus", ProjectDetailStatus.HIDE)
                        this.$store.commit("setCurrentTask", new Task())
                        this.$emit("reload")
                        this.$router.push({name: "taskBoard", params: {projectId: this.projectId}})
                    }
                }).finally(() => {
                    this.$store.commit("decrementLoadingCount")
                })
            },
            clickDelete() {
                this.$store.commit("setProjectDetailStatus", ProjectDetailStatus.DELETE_CONFIRM)
            },
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