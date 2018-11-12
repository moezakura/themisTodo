<template>
    <transition>
        <div v-if="value">
            <div id="project-settings">
                <h3>Project Config</h3>
                <div id="project-settings-container">
                    <ul id="project-settings-menus">
                        <li @click="changeContent('overview')">
                            <i class="fas fa-cog"></i>
                            <span>overview</span>
                        </li>
                        <li @click="changeContent('member')">
                            <i class="fas fa-users"></i>
                            <span>members</span>
                        </li>
                        <li @click="changeContent('danger')">
                            <i class="fas fa-exclamation-triangle"></i>
                            <span>danger</span>
                        </li>
                    </ul>
                    <div id="project-settings-content">
                        <div class="error" v-show="errorMessage !== undefined && errorMessage.length > 0"
                             @click="clearMessage">{{errorMessage }}
                        </div>
                        <component :is="componentName"></component>
                    </div>
                </div>
            </div>
            <div class="backView" @click="hideSettings"></div>
        </div>
    </transition>
</template>

<script lang="ts">
    import Overview from "@components/TaskBoard/Settings/Overview"
    import Member from "@components/TaskBoard/Settings/Member"
    import Danger from "@components/TaskBoard/Settings/Danger"

    export default {
        name: "ProjectSettings",
        props: ['value'],
        components: {Overview, Member, Danger},
        data: () => {
            return {
                componentName: "overview"
            }
        },
        computed: {
            projectId(): number {
                if (this.$store.getters.getCurrentProject == undefined ||
                    this.$store.getters.getCurrentProject.uuid == undefined) {
                    return 0
                }
                return this.$store.getters.getCurrentProject.uuid
            },
            errorMessage(): string {
                return this.$store.getters.getProjectSettings.errorMessage
            }
        },
        methods: {
            hideSettings() {
                this.$router.back()
            },
            changeContent(name) {
                this.componentName = name
            },
            clearMessage() {
                this.$store.commit("setProjectSettingsProps", {
                    key: "errorMessage",
                    value: "",
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

    #project-settings-container {
        display: flex;
        justify-content: space-between;
    }

    #project-settings-menus {
        list-style: none;
        margin: 0 10px 0 0;
        padding: 0;
        height: 100%;
        width: 190px;

        li {
            height: 45px;
            line-height: 45px;
            font-size: 16px;
            letter-spacing: 1px;
            width: 100%;
            box-sizing: border-box;
            padding: 0 8px;

            &:hover {
                background-color: rgba(67, 160, 71, 0.15);
            }

            i {
                margin-right: 10px;
            }
        }
    }

    #project-settings-content {
        width: calc(100% - 200px);
    }
</style>