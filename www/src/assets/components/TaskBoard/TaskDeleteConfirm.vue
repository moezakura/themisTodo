<template>
    <transition>
        <div v-show="isShow" class="task-delete-popup-container">
            <div id="taskDeleteConfirmPopup">
                <div id="taskDeleteConfirmPopupTitle">Delete confirm</div>
                <p class="taskDeleteConfirmPopupText">Delete the task from the task board.</p>
                <p class="taskDeleteConfirmPopupText">This operation can not be undone.</p>
                <div id="taskDeleteConfirmPopupSelectButtons">
                    <div id="taskDeleteConfirmPopupDelete" class="deleteButton" @click="clickDelete">DELETE</div>
                    <div id="taskDeleteConfirmPopupCancel" class="deleteButton" @click="clickClose">CANCEL</div>
                </div>
            </div>
            <div class="backView" @click="clickClose"></div>
        </div>
    </transition>
</template>

<script lang="ts">
    import {ProjectDetailStatus} from "@scripts/enums/ProjectDetailStatus"
    import Task from "../../scripts/model/api/task/Task"
    import TaskApi from "../../scripts/api/TaskApi"

    export default {
        name: "TaskDeleteConfirm",
        computed: {
            isShow(): boolean {
                return this.$store.getters.getProjectDetailStatus == ProjectDetailStatus.DELETE_CONFIRM
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
            clickDelete() {
                let task = this.task
                this.$store.commit("incrementLoadingCount")
                TaskApi.Delete(task.createDate).then(res => {
                    if (res.success) {
                        this.$store.commit("setProjectDetailStatus", ProjectDetailStatus.HIDE)
                        this.$store.commit("setCurrentTask", new Task())
                        this.$emit("reload")
                        this.$router.push({name: "taskBoard", params: {projectId: this.projectId}})
                    }
                }).finally(() =>{
                    this.$store.commit("decrementLoadingCount")
                })
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