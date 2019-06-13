<template>
    <div>
        <project-header>
            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-menu>
                <li @click=""><i class="fas fa-redo"></i>RELOAD</li>
            </template>

            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-right-menu>
                <li @click="addEntryFocus"><i class="fas fa-plus-circle"></i>ADD ENTRY</li>
            </template>
        </project-header>

        <form class="task-timer-add basicForm" :class="{'active': taskTimerTopFocus}">
            <label class="task-timer-entry-name">
                <input placeholder="What did you task on?" @focus="taskTimerTopFocus = true"
                       @blur="taskTimerTopFocus = false" ref="task-timer-entry-name">
            </label>
            <input class="task-timer-entry-submit" type="submit" value="Add entry">
        </form>

        <div class="task-timer-history">
            <div class="task-timer-history-title-section">
                <h3 class="task-timer-history-title"><i class="fas fa-calendar-alt"></i>Today</h3>
                <div class="task-timer-history-title-detail">2019/12/04</div>
                <div class="task-timer-history-title-action">
                    <i class="fas fa-chevron-left"></i>
                    <i class="fas fa-chevron-right"></i>
                </div>
            </div>
            <ul class="task-timer-entry-container">
                <li class="task-timer-entry">
                    <div class="label"></div>
                    <div class="name">TASK NAME</div>
                    <div class="note">NOTE NOTE NOTE NOTE NOTE NOTE NOTE NOTE NOTE NOTE NOTE NOTE NOTE NOTE NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE  NOTE </div>
                    <div class="date">
                        <div class="time-between">
                            <span>12:00</span>
                            <span>〜</span>
                            <span>12:01</span>
                        </div>
                        <div class="time-total">
                            <span class="time-total-text">99:59</span>
                        </div>
                    </div>
                    <div class="actions">
                        <i class="fas fa-trash"></i>
                        <i class="fas fa-edit"></i>
                        <i class="fas fa-caret-down"></i>
                    </div>
                </li>
            </ul>
        </div>
    </div>
</template>

<script lang="ts">
    import Project from "@scripts/model/api/project/Project";
    import ProjectApi from "@scripts/api/ProjectApi";
    import ProjectHeader from "@components/Project/ProjectHeader.vue";

    export default {
        name: "TimerBoard",
        components: {ProjectHeader},
        data() {
            return {
                taskTimerTopFocus: false,
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
            }
        },
        methods: {
            addEntryFocus(): void {
                this.$refs['task-timer-entry-name'].focus()
            }
        },
        created(): void {
            this.$store.commit("incrementLoadingCount")
            ProjectApi.getProject(this.projectId).then(res => {
                this.$store.commit("setCurrentProject", res.project)
            }).finally(() => {
                this.$store.commit("decrementLoadingCount")
            })
        }
    }
</script>

<style scoped lang="scss">
    .task-timer-add {
        display: flex;
        width: 75%;
        margin: 10px auto;
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

    .task-timer-history {
        width: 85%;
        margin: 15px auto;

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

            .task-timer-entry {
                $height: 50px;
                display: flex;
                height: $height;
                border-bottom: solid 1px rgba(white, .3);

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
            }
        }
    }
</style>