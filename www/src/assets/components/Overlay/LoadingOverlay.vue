<template>
    <transition name="overlay">
        <div id="loading-overlay" v-show="isShow">
            <div class="loadingBody"></div>
        </div>
    </transition>
</template>

<script lang="ts">
    export default {
        name: "LoadingOverlay",
        computed: {
            isShow(): boolean {
                return this.$store.getters.isLoadingShow
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