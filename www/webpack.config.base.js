const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const {VueLoaderPlugin} = require('vue-loader');

const path = require('path');
const dist = path.resolve(__dirname, 'dist');

module.exports = (isDev) => {
    const enabledSourceMap = isDev;

    return {
        entry: './src/main.ts',
        mode: isDev ? "development" : "production",
        output: {
            filename: '[name]_[hash].js',
            path: dist,
            publicPath: '/'
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
                    test: /\.scss/,
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
                                // minimize: !enabledSourceMap,
                                importLoaders: 2
                            },
                        },
                        {
                            loader: 'sass-loader',
                            options: {
                                sourceMap: enabledSourceMap,
                                data: `
                                @import './src/assets/styles/_value.scss';
                                `
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
                '@components': path.resolve(__dirname, './src/assets/components'),
                '@scripts': path.resolve(__dirname, './src/assets/scripts'),
            }
        },
        plugins: [
            new webpack.NamedModulesPlugin(),
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
        ]
    };
};
if (module.hot) {
    module.hot.accept();
}
