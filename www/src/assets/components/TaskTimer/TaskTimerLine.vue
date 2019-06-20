<template>
    <li class="task-timer-entry">
        <div class="label"></div>
        <div class="name">{{ task.name }}</div>
        <div class="note">{{ taskTimer.note }}</div>
        <div class="date">
            <div class="time-between">
                <span>{{ taskTimer.startDateHM }}</span>
                <span>〜</span>
                <span>{{ taskTimer.endDateHM }}</span>
            </div>
            <div class="time-total">
                <span class="time-total-text">{{ taskTimer.totalHM }}</span>
            </div>
        </div>
        <div class="actions">
            <i class="fas fa-stop" @click="stopTask(task)" v-if="taskTimer.endDateUnix === 0"></i>
            <i class="fas fa-trash" @click="deleteTask(taskTimer)" v-else></i>
            <i class="fas fa-edit" @click="isEdit=true"></i>
            <i class="fas fa-info"></i>
        </div>

        <transition>
            <div class="timer-edit-modal-container" v-show="isEdit">
                <form class="timer-edit-modal basicForm" @submit.prevent="applyChange">
                    <h2>Task Timer Edit</h2>
                    <div class="timer-inputs">
                        <input type="time" v-model="edit.startDate">
                        <div class="split">〜</div>
                        <input type="time" v-model="edit.endDate">
                    </div>
                    <label class="timer-note">
                        <textarea placeholder="note"></textarea>
                    </label>
                    <div class="timer-actions">
                        <input type="submit" value="CANCEL" @click.prevent="isEdit=false">
                        <input type="submit" value="CHANGE">
                    </div>
                </form>
                <div class="timer-edit-modal-background" @click="isEdit=false"></div>
            </div>
        </transition>
    </li>
</template>

<script lang="ts">
    import Task from "@scripts/model/api/task/Task";
    import TaskTimerApi from "@scripts/api/TaskTimer";
    import TaskTimer from "@scripts/model/api/taskTimer/TaskTimer";

    export default {
        name: "TaskTimerLine",
        props: {
            taskTimer: TaskTimer
        },
        data() {
            return {
                isEdit: false,
                edit: {
                    startDate: "",
                    endDate: "",
                    note: "",
                }
            }
        },
        computed: {
            task(): Task {
                return this.taskTimer.task
            }
        },
        watch: {
            isEdit(value: boolean, oldValue: boolean): void {
                if (value && !oldValue) {
                    this.$set(this.edit, "startDate", this.taskTimer.startDateHM)
                    this.$set(this.edit, "endDate", this.taskTimer.endDateHM)
                    this.$set(this.edit, "note", this.taskTimer.note)
                }
            },
            taskTimer(value: TaskTimer): void {
                if (this.isEdit) {
                    return
                }

                this.$set(this.edit, "startDate", value.startDateHM)
                this.$set(this.edit, "endDate", value.endDateHM)
                this.$set(this.edit, "note", value.note)
            }
        },
        methods: {
            loadPage(): void {
                this.$emit("load-page")
            },
            async stopTask(task: Task): Promise<void> {
                this.$store.commit("incrementLoadingCount")
                const res = await TaskTimerApi.getTaskTimerStatus(task.createDate)

                if (res.start) {
                    const toggleRes = await TaskTimerApi.toggleTimer(task.createDate)
                    if (toggleRes.success) {
                        this.loadPage()
                    }
                }
                this.$store.commit("decrementLoadingCount")
            },
            async deleteTask(taskTimer: TaskTimer): Promise<void> {
                this.$store.commit("incrementLoadingCount")

                const toggleRes = await TaskTimerApi.deleteTaskTimer(taskTimer.id)
                if (toggleRes.success) {
                    this.loadPage()
                }
                this.$store.commit("decrementLoadingCount")
            }
        }
    }
</script>

<style scoped lang="scss">
    .v-enter,
    .v-leave-to {
        opacity: 0;
    }

    .v-enter-active,
    .v-leave-active {
        transition: opacity .2s;
    }

    .task-timer-entry {
        $height: 55px;
        display: flex;
        height: $height;
        margin-bottom: 5px;
        padding-bottom: 5px;
        border-bottom: solid 1px rgba(white, .3);
        flex-wrap: wrap;

        &::after {
            position: relative;
            display: block;
            content: " ";
            border-bottom: solid 3px transparent;
            width: 100%;
            height: 0;
            margin: 3px 0 0 0;
        }

        &.active::after {
            @include animation(blinkBorderBottomAnimation 1s infinite, linear);
        }

        .label {
            width: 5px;
            background-color: red;
            margin-right: 10px;
            height: $height;
        }

        .name {
            width: 40%;
            line-height: $height;
            overflow: hidden;
            font-weight: bold;
            font-size: 16px;
        }

        .note {
            width: 25%;
            overflow-y: auto;
        }

        .date {
            margin: 0 auto;

            .time-between,
            .time-total {
                height: $height / 2;
                line-height: $height / 2;
            }

            .time-total {
                text-align: right;

                .time-total-text::before {
                    content: "total:";
                    margin-right: 5px;
                }
            }
        }

        .actions {
            display: flex;
            margin-left: 5px;

            i {
                $margin: 5px;
                width: $height;
                height: $height - $margin * 2;
                line-height: $height - $margin * 2;
                margin: $margin 0;
                font-size: 16px;
                text-align: center;
                box-sizing: border-box;
                background-color: transparent;
                border: solid 1px $accentColor;
                color: white;
                padding: 0 15px;
                letter-spacing: 1px;
                transition: ease background-color .3s;

                &:hover {
                    background-color: accentColor(.1);
                }

                &:not(:last-child) {
                    border-right: 0;
                }
            }
        }

        .timer-edit-modal-container {
            .timer-edit-modal {
                position: fixed;
                left: 50%;
                top: 50%;
                transform: translate(-50%, -50%);
                padding: 10px 20px;
                box-sizing: border-box;
                box-shadow: 0 0 6px rgba(black, .5);
                background-color: $backgroundColor;
                width: 500px;
                z-index: 99;

                h2 {
                    height: $buttonHeight;
                    line-height: $buttonHeight;
                    font-size: 18px;
                    padding-left: 15px;
                    background-color: rgba(black, 0.5);
                    letter-spacing: 1.5px;
                    margin: -10px -20px 0 -20px;
                    text-align: left;
                }

                .timer-inputs,
                .timer-actions {
                    display: flex;
                    height: $lineTextHeight;
                    line-height: $lineTextHeight;
                    margin-top: 15px;

                    input {
                        margin: 0;
                    }

                    .split {
                        margin: 0 10px;
                    }
                }

                .timer-note {
                    display: block;
                    width: 100%;
                    margin: 10px 0;

                    textarea {
                        display: block;
                        width: 100%;
                        min-width: 100%;
                        max-width: 100%;
                        height: 250px;
                        resize: none;
                    }
                }

                .timer-actions {
                    margin-bottom: 15px;

                    input:first-child {
                        font-size: 12px;
                        margin: 0 10px 0 auto;
                    }

                    input:last-child {
                        font-size: 12px;
                        margin: 0 0 0 10px;
                    }
                }
            }

            .timer-edit-modal-background {
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                background-color: rgba(black, .7);
                z-index: 98;
            }
        }
    }
</style>