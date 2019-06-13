<template>
    <div>
        <project-header>
            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-menu>
                <li @click=""><i class="fas fa-redo"></i>RELOAD</li>
            </template>

            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-right-menu>
                <li @click=""><i class="fas fa-plus-circle"></i>ADD ENTRY</li>
            </template>
        </project-header>

        <form class="task-timer-add basicForm" :class="{'active': taskTimerTopFocus}">
            <label class="task-timer-entry-name">
                <input placeholder="What did you task on?" @focus="taskTimerTopFocus = true"
                       @blur="taskTimerTopFocus = false">
            </label>
            <input class="task-timer-entry-submit" type="submit" value="Add entry">
        </form>
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
</style>