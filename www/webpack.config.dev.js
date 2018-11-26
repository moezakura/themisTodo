const webpack = require('webpack');
const merge = require('webpack-merge');
const baseConfig = require('./webpack.config.base.js');

module.exports = merge(baseConfig(true), {
    devServer: {
        port: 8652,
        historyApiFallback: true,
        inline: true,
        proxy: {
            '/api': {
                target: 'http://localhost:31204',
                pathRewrite: {'^/api': ''}
            }
        }
    },
    resolve: {
        alias: {
            vue: 'vue/dist/vue.js',
        }
    },
    plugins: [
        new webpack.HotModuleReplacementPlugin(),
    ]
});
