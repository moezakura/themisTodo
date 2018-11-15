<template>
    <div class="admin-top">
        <div class="left">
            <ul>
                <li>
                    <div class="admin-menu-title">Accounts</div>
                    <ul class="admin-section-menu">
                        <li>
                            <router-link to="/admin/account/list">Account List</router-link>
                        </li>
                        <li>
                            <router-link to="/admin/account/new">New Account</router-link>
                        </li>
                    </ul>
                </li>
            </ul>
        </div>
        <div class="right">
            <router-view name="admin-pages"></router-view>
        </div>
    </div>
</template>

<script>
    export default {
        name: "AdminTop"
    }
</script>

<style lang="scss" scoped>
    .admin-top {
        display: flex;

        .left {
            width: 200px;
            margin-left: 30px;
            letter-spacing: 1px;
        }
        .right {
            width: calc(100% - #{200px + 30px});
        }

        .admin-menu-title {
            font-size: 16px;
            font-weight: 900;
        }

        .admin-section-menu {
            font-size: 14px;

            a {
                position: relative;
                display: block;
                width: 100%;
                height: 45px;
                line-height: 45px;
                padding-left: 25px;
                color: white;
                text-decoration: none;
                transition: ease background-color .3s;
                box-sizing: border-box;
                z-index: 1;

                &:hover {
                    background-color: rgba(67, 160, 71, .3);
                }
                &.router-link-active::after {
                    position: relative;
                    content: ' ';
                    display: block;
                    width: 100%;
                    height: 100%;
                    margin: -45px 0 0 -25px;
                    background: -moz-linear-gradient(left, rgba(67,160,71,1) 0%, rgba(67,160,71,0) 100%);
                    background: -webkit-linear-gradient(left, rgba(67,160,71,1) 0%,rgba(67,160,71,0) 100%);
                    background: linear-gradient(to right, rgba(67,160,71,1) 0%,rgba(67,160,71,0) 100%);
                    filter: progid:DXImageTransform.Microsoft.gradient( startColorstr='#43a047', endColorstr='#0043a047',GradientType=1 );
                    z-index: -1;
                }
            }
        }
    }
</style>