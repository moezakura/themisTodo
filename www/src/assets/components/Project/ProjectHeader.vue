<template>
    <div class="project-title-container">
        <h2 class="project-title"><i class="fas fa-tasks"></i><span>{{ storeProject.name }}</span></h2>

        <slot name="inner-content"></slot>

        <ul class="project-actions">
            <li @mouseenter="isShowOtherMenu = true" @mouseleave="isShowOtherMenu = false">
                <i class="fas fa-toolbox"></i>OTHER
                <ul class="project-actions-other" v-show="isShowOtherMenu">
                    <slot name="ex-menu"></slot>
                    <li @click="moveHideTasks"><i class="fas fa-eye-slash"></i>HIDE TASKS</li>
                    <li @click="moveSettings"><i class="fas fa-cog"></i>SETTING</li>
                </ul>
            </li>
            <li @click="moveTimer"><i class="fas fa-stopwatch"></i>TIMERS</li>
            <slot name="ex-right-menu"></slot>
        </ul>
    </div>
</template>

<script lang="ts">
    import Project from "@scripts/model/api/project/Project";

    export default {
        name: "ProjectHeader",
        data() {
            return {
                isShowOtherMenu: false,
            }
        },
        computed: {
            storeProject(): Project {
                if (this.$store.getters.getCurrentProject == undefined) {
                    return new Project()
                }
                return this.$store.getters.getCurrentProject
            },
            projectId():number{
                return this.storeProject.uuid;
            }
        },
        methods: {
            moveHideTasks() {
                this.$router.push({name: "hiddenTasks", params: {projectId: this.projectId}})
            },
            moveSettings() {
                this.$router.push({name: "projectSettings", params: {projectId: this.projectId}})
            },
            moveTimer() {
                this.$router.push({name: "timerBoard", params: {projectId: this.projectId}})
            }
        }
    }
</script>

<style lang="scss">
    .project-title-container {
        .project-title {
            display: block;
            height: $buttonHeight;
            line-height: $buttonHeight;
            font-size: 18px;
            letter-spacing: 1.2px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            flex-grow: 5;

            i {
                margin-right: 10px;
            }
        }

        .project-actions {
            display: flex;

            li {
                display: block;
                font-size: 16px;
                height: $buttonHeight;
                width: 140px;
                line-height: $buttonHeight;
                text-align: center;
                letter-spacing: 1.1px;
                text-decoration: none;
                margin: 0 3px;
                color: $accentColor;
                cursor: default;
                -webkit-user-select: none;
                -moz-user-select: none;
                -ms-user-select: none;
                user-select: none;
                transition: background-color ease .3s;

                &:hover {
                    background-color: accentColor(.15);
                }

                i {
                    margin-right: 10px;
                }
            }

            .project-actions-other {
                position: absolute;
                top: 65px + 10px + 45px;
                right: 20px + (140px + 3px * 2) * 2 - 15px - 40px;
                background-color: rgb(30, 30, 30);
                box-shadow: 3px 6px 6px rgba(0, 0, 0, 0.8);
                z-index: 50;

                & > li {
                    padding-left: 15px;
                    text-align: left;
                    width: 180px;
                }
            }
        }
    }
</style>