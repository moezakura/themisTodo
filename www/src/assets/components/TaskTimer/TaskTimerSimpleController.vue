<template>
    <div class="task-timer-simple-controller">

        <div class="task-operation-button base-timer-button" @click="toggle" v-if="this.isStart"><span
                class="operation-button-text">Start Task</span><i
                class="fas fa-play operation-button-icon"></i></div>
        <div class="task-operation-button base-timer-button stop-timer-button" @click="toggle" v-else><span
                class="operation-button-text">Stop Task</span><i class="fas fa-stop operation-button-icon"></i></div>

        <div class="timer-start-and-end-text">
            <div class="timer-start">00:00</div>
            <span class="timer-text-splitter">〜</span>
            <div class="timer-end">00:00</div>
        </div>
        <div class="timer-total-texts">
            <div class="timer-today-text"><span
                    class="timer-today-text-label">today</span><span>:</span><span class="timer-today-text-content">24:00</span>
            </div>
            <div class="timer-total-text"><span
                    class="timer-total-text-label">total</span><span>:</span><span class="timer-total-text-content">9999:59</span>
            </div>
        </div>
        <div class="task-timer-edit-button base-timer-button"><i class="fas fa-edit"></i></div>
        <div class="task-timer-history-button base-timer-button"><i class="fas fa-history"></i></div>
    </div>
</template>

<script lang="ts">
    import TaskTimerApi from "@scripts/api/TaskTimer";
    import Task from "@scripts/model/api/task/Task";

    interface TaskTimerSimpleControllerData {
        isStart: boolean
    }

    export default {
        name: "TaskTimerSimpleController",
        data(): TaskTimerSimpleControllerData {
            return {
                isStart: false
            }
        },
        methods: {
            toggle() {
                const task: Task | undefined = this.$store.getters.getCurrentTask
                if (typeof task !== "undefined" && typeof task.createDate !== "undefined") {
                    TaskTimerApi.toggleTimer(task.createDate).then(res => {
                        this.isStart = res.start;
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .task-timer-simple-controller {
        $height: 45px;
        display: flex;

        .base-timer-button {
            border: thin solid rgba($accentColor, 0.5);
            background-color: rgba($accentColor, 0.05);
            transition: background-color ease .3s;
            -webkit-user-select: none;
            -moz-user-select: none;
            -ms-user-select: none;
            user-select: none;

            &:hover {
                background-color: rgba($accentColor, 0.2);
            }

            &.stop-timer-button {
                background-color: rgba($redColor, 0.05);
                border: 1px solid rgba($redColor, 0.5);

                &:hover {
                    background-color: rgba($redColor, 0.2);
                }
            }
        }

        .task-operation-button {
            $button-height: 40px;
            display: flex;
            width: 180px;
            height: $button-height;
            line-height: $button-height;
            text-align: center;
            font-size: 14px;
            margin: #{($height - $button-height) / 2} 0;

            .operation-button-text {
                margin: 0 5px 0 auto;
                vertical-align: middle;
            }

            .operation-button-icon {
                line-height: $button-height;
                text-align: center;
                margin: 0 auto 0 5px;
                vertical-align: middle;
            }
        }

        .timer-start-and-end-text {
            margin: 0 5px 0 auto;
            display: flex;
            line-height: $height;
            text-align: center;

            .timer-text-splitter {
                width: 35px;
            }

            .timer-start, .timer-end {
                width: 80px;
                letter-spacing: 1px;
            }
        }

        .timer-total-texts {
            margin: 0 auto 0 5px;
            line-height: $height / 2;
            text-align: left;
            letter-spacing: 1px;
            width: 110px;

            .timer-today-text,
            .timer-total-text {
                width: 100%;
                display: flex;

                span {
                    display: block;
                }

                .timer-today-text-label,
                .timer-total-text-label {
                    width: 38%;
                }

                .timer-today-text-content,
                .timer-total-text-content {
                    margin-left: 5px;
                }
            }
        }

        .task-timer-edit-button,
        .task-timer-history-button {
            $button-height: 40px;
            width: 60px;
            height: $button-height;
            line-height: $button-height;
            margin: #{($height - $button-height) / 2} 0;
            text-align: center;
            font-size: 16px;
        }

        .task-timer-history-button {
            margin-right: 0;
            border-left: 0;
        }
    }
</style>