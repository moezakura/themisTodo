const MODE = 'development';
const enabledSourceMap = (MODE === 'development');

const webpack = require('webpack');

module.exports = {
    entry: './js/main.js',
    mode: "development",
    output: {
        // ディレクトリ名
        path: `${__dirname}/`,
        // ファイル名
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
        ]
    },
    plugins: [
        new webpack.NamedModulesPlugin(),
        new webpack.HotModuleReplacementPlugin(),
    ],
};
if (module.hot) {
    module.hot.accept();
}
