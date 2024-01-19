const HtmlWebpackPlugin    = require('html-webpack-plugin');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const path                 = require('path');

module.exports = {
    entry: './frontend/main.js',
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'public'),
        publicPath: ""
    },
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: [MiniCssExtractPlugin.loader, 'css-loader', 'postcss-loader'],
            },
        ],
    },
    plugins: [
        new HtmlWebpackPlugin(
            {
                template: "./views/base.html",
                filename: "./../views/index.html",
            }
        ),
        new MiniCssExtractPlugin(),
    ]
};
