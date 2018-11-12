<template>
    <transition name="task-detail">
        <div v-if="isShowTaskDetail">
            <form id="taskPopup">
                <h2>Doing Task</h2>
                <div id="taskPopupActions">
                    <i class="fas fa-trash" id="taskPopupTrashButton"></i>
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
                        <div id="taskPopupAssignIcon"><!-- TODO: ユーザーアイコンの設定 --></div>
                        <input value="ASSIGN" id="taskPopupAssign" autocomplete="off">
                    </div>
                    <div class="taskPopupAssignCreatorColumn">
                        <p>Creator</p>
                        <div id="taskPopupCreatorIcon"><!-- TODO: ユーザーアイコンの設定 --></div>
                        <label id="taskPopupCreator">{{ task.creatorName }}</label>
                    </div>
                </div>
                <div class="both"></div>
                <div id="taskPopupProgressBar">
                    <div id="taskPopupProgressText" v-show="!isEditing">
                        <i class="fas fa-calendar-alt"></i>
                        <span>{{ task.deadline }}</span>
                        <span>(あと{{ task.limitDate }}日)</span>
                    </div>
                    <input type="date" id="taskPopupDeadlineChange" v-show="isEditing">
                    <div id="taskPopupProgressCurrent">&nbsp;</div>
                </div>
                <textarea id="taskPopupDescription" v-model="taskCache.description" :readonly="!isEditing"></textarea>
                <div class="input-box" v-show="isEditing">
                    <input type="button" value="CANCEL" @click="setEditing(false)">
                    <input type="submit" value="CHANGE">
                </div>
            </form>
            <div class="backView" v-if="isShowTaskDetail" @click="hideTaskDetail"></div>
        </div>
    </transition>
</template>

<script lang="ts">
    import Task from "../../scripts/model/api/task/Task"

    export default {
        name: "TaskDetail",
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
            task(): Task | undefined {
                this.taskCache = this.$store.getters.getCurrentTask
                return this.taskCache
            },
            isShowTaskDetail(): boolean {
                return this.task !== undefined
            }
        },
        methods: {
            hideTaskDetail() {
                this.$router.back()
                this.$store.commit("setCurrentTask", undefined)
            },
            setEditing(isEditing: boolean) {
                this.isEditing = isEditing
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