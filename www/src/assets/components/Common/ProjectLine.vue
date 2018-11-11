<template>
    <li class="project-line">
        <div class="name">
            <i class="fas fa-th-list"></i>
            <span>{{ name }}</span>
        </div>
        <p class="description">{{ description }}</p>
    </li>
</template>

<script lang="ts">
    export default {
        name: "ProjectLine",
        props: ['project'],
        computed: {
            uuid(): Number {
                return this.project.uuid
            },
            name(): string {
                return this.project.name
            },
            description(): string {
                return this.project.description
            }
        }
    }
</script>

<style lang="scss" scoped>
    .project-line {
        border-bottom: solid 1px rgb(0, 57, 34);
        padding: 10px 8px 10px 16px;
        transition: background-color ease .3s;
        cursor: default;

        &:hover {
            background-color: rgba(black, .5);
        }
        &:last-child {
            border: 0;
        }

        .name {
            height: 40px;
            line-height: 40px;

            i, span {
                font-size: 17px;
                font-weight: 900;
            }

            i {
                margin-right: 15px;
            }
        }

        .description {
            font-size: 16px;
            color: rgb(200, 200, 200);
        }
    }
</style>