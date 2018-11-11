const MODE = 'development';
const enabledSourceMap = (MODE === 'development');

const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const {VueLoaderPlugin} = require('vue-loader');

const path = require('path');
const dist = path.resolve(__dirname, 'dist');

module.exports = {
    entry: './src/main.ts',
    mode: "development",
    output: {
        filename: '[name].js',
        path: dist,
        publicPath: '/'
    },
    devServer: {
        port: 8652,
        historyApiFallback: true,
        proxy: {
            '/api': {
                target: 'http://localhost:31204',
                pathRewrite: {'^/api' : ''}
            }
        }
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
                test: /\.css/,
                use: [
                    'style-loader',
                    {loader: 'css-loader', options: {url: false}},
                ],
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
                            // オプションでCSS内のurl()メソッドの取り込みを禁止する
                            url: false,
                            // ソースマップの利用有無
                            sourceMap: enabledSourceMap,

                            // 0 => no loaders (default);
                            // 1 => postcss-loader;
                            // 2 => postcss-loader, sass-loader
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
                test: /\.html$/,
                loader: 'html-loader'
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
        extensions: ['*', '.js', '.vue', '.json', '.tsx', '.ts'],
        alias: {
            vue: 'vue/dist/vue.js',
            '@components': path.resolve(__dirname, './src/assets/components'),
            '@scripts': path.resolve(__dirname, './src/assets/scripts'),
        }
    },
    plugins: [
        new webpack.NamedModulesPlugin(),
        new webpack.HotModuleReplacementPlugin(),
        new VueLoaderPlugin(),
        new HtmlWebpackPlugin({
            template: './src/index.html',
            inject: 'body'
        }),
        new CopyWebpackPlugin([
            {
                from: path.resolve(__dirname, 'src/assets/fontawesome/web-fonts-with-css/webfonts/'),
                to: path.resolve(dist, './webfonts/'),
            },
            {
                from: path.resolve(__dirname, 'src/assets/images/'),
                to: path.resolve(dist, './assets/images/'),
            }
        ]),
    ],
};
if (module.hot) {
    module.hot.accept();
}
