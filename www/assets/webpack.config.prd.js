const MODE = 'production';
const enabledSourceMap = (MODE === 'development');

const webpack = require('webpack');
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')

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
                        loader: 'style-loader?modules',
                        options: {
                            hmr: true,
                            singleton: true,
                        }
                    },
                    {
                        loader: 'css-loader?modules',
                        options: {
                            url: false,
                            sourceMap: enabledSourceMap,
                            minimize: true,
                            importLoaders: 2
                        },
                    },
                    {
                        loader: 'sass-loader?modules',
                        options: {
                            // ソースマップの利用有無
                            sourceMap: enabledSourceMap,
                        }
                    }
                ],
            },
        ]
    },
    plugins: [
        new webpack.NamedModulesPlugin(),
        new webpack.HotModuleReplacementPlugin(),
        new UglifyJsPlugin({
            uglifyOptions: {
                compress: {
                    top_retain: true,
                    warnings: false,
                    drop_console: true,
                }
            }
        }),
    ],
};
if (module.hot) {
    module.hot.accept();
}
