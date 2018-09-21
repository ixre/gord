const path = require('path');
const webpack = require("webpack")
const process = require("process")
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const devMode = process.env.NODE_ENV !== 'production'

module.exports = {
    mode: 'production',
    entry: {
        app: './index.js'
    },
    // devtool: 'inline-source-map',
    devServer: {
        contentBase: './dist',
        port: 3000,
    },
    output: {
        filename: '[name].bundle-[hash:6].js',
        chunkFilename: '[name].bundle.js',
        path: path.resolve(__dirname, 'dist')
    },
    plugins: [
        new CleanWebpackPlugin(['dist']),
        new HtmlWebpackPlugin({title: "GORD", template: "./public/entry.html", appMountId: "root"}),
        new webpack.optimize.SplitChunksPlugin({names: ["common", "vendor", "manifest"]}),
        new MiniCssExtractPlugin({
            filename: devMode ? '[name].css' : '[name].[hash].css',
            chunkFilename: devMode ? '[id].css' : '[id].[hash].css',
        })
    ],
    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: /node_modules/,
                include: path.resolve(__dirname, 'src'),
                loader: "babel-loader"
            },
            {test: /\.jsx?$/, loader: 'babel-loader'},
            {test: /\.css$/, use: [{loader: MiniCssExtractPlugin.loader}, "css-loader"]},
            {test: /\.less$/, loader: 'style-loader!css-loader!less-loader'},
            {test: /\.(png|jpg)$/, loader: 'url-loader?limit=25000'}
        ]
    }
}
