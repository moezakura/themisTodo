const MODE = 'production';
const enabledSourceMap = (MODE === 'development');

const webpack = require('webpack');
const {VueLoaderPlugin} = require('vue-loader');

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
                test: /\.ts$/,
                use: [
                    {
                        loader: 'ts-loader',
                        options: {
                            appendTsSuffixTo: [/\.vue$/]
                        }
                    }
                ]
            },
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
                        ['env', {'modules': false}]
                    ]
                },
                exclude: /node_modules/
            },
        ]
    },
    resolve: {
        extensions: ['*', '.js', '.vue', '.json', '.ts'],
        alias: {
            vue: 'vue/dist/vue.js'
        }
    },
    plugins: [
        new webpack.NamedModulesPlugin(),
        new webpack.HotModuleReplacementPlugin(),
        new VueLoaderPlugin(),
    ]
};
if (module.hot) {
    module.hot.accept();
}
