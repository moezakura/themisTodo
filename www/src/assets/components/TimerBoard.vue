<template>
    <div>
        <project-header>
            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-menu>
                <li @click=""><i class="fas fa-redo"></i>RELOAD</li>
            </template>

            <!--suppress HtmlUnknownBooleanAttribute -->
            <template v-slot:ex-right-menu>
                <li @click=""><i class="fas fa-plus-circle"></i>ADD TIMER</li>
            </template>
        </project-header>
    </div>
</template>

<script lang="ts">
    import Project from "@scripts/model/api/project/Project";
    import ProjectApi from "@scripts/api/ProjectApi";
    import ProjectHeader from "@components/Project/ProjectHeader.vue";

    export default {
        name: "TimerBoard",
        components: {ProjectHeader},
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

<style scoped>

</style>