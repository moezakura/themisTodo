const merge = require('webpack-merge');
const baseConfig = require('./webpack.config.base.js');

module.exports = merge(baseConfig(false), {
    resolve: {
        alias: {
            vue: 'vue/dist/vue.min.js',
        }
    },
});
