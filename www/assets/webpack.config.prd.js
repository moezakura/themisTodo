const MODE = 'production';
const enabledSourceMap = (MODE === 'development');

const webpack = require('webpack');

module.exports = {
    entry: './js/main.js',
    mode: "production",
    output: {
        path: `${__dirname}/`,
        filename: 'bundle.js'
    },

    devServer: {
        port: 8652,
    },
    module: {
        rules: [
            {
                test: /\.scss/, // 対象となるファイルの拡張子
                use: [
                    {
                        loader: 'style-loader',
                        options: {
                            hmr: true,
                            singleton: true,
                        }
                    },
                    {
                        loader: 'css-loader',
                        options: {
                            url: false,
                            sourceMap: enabledSourceMap,
                            minimize: true,
                            importLoaders: 2
                        },
                    },
                    {
                        loader: 'sass-loader',
                        options: {
                            // ソースマップの利用有無
                            sourceMap: enabledSourceMap,
                        }
                    }
                ],
            },
            {
                test: /\.vue$/,
                loader: 'vue-loader'
            },
            {
                test: /\.js$/,
                loader: 'babel-loader?optional[]=runtime',
                options: {
                    presets: [
                        ['env', { 'modules': false }]
                    ]
                },
                exclude: /node_modules/
            },
        ]
    },
    resolve: {
        extensions: ['*', '.js', '.vue', '.json'],
        alias: {
            vue: 'vue/dist/vue.js'
        }
    },
    plugins: [
        new webpack.NamedModulesPlugin(),
        new webpack.HotModuleReplacementPlugin(),
    ]
}
if (module.hot) {
    module.hot.accept();
}
