<template>
    <transition>
        <div v-if="value">
            <div></div>
            <div class="backView"></div>
        </div>
    </transition>
</template>

<script lang="ts">
    export default {
        name: "ProjectSettings",
        props: ['value'],
        data: () => {
            return {}
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