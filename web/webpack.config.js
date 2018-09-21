const path = require('path');
const webpack = require("webpack")
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');

module.exports = {
    mode: 'production',
    entry: {
        app: './index.js'
    },
    devtool: 'inline-source-map',
    devServer: {
        contentBase: './dist',
        port:3000,
    },
    output: {
        filename: '[name].bundle-[hash:6].js',
        chunkFilename: '[name].bundle.js',
        path: path.resolve(__dirname, 'dist')
    },
    plugins: [
        new CleanWebpackPlugin(['dist']),
        new HtmlWebpackPlugin({title: "GORD",template:"./public/entry.html",appMountId:"root"}),
        new webpack.optimize.SplitChunksPlugin({names:["common","vendor","manifest"]}),
    ],
    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: path.resolve(__dirname, 'node_modules'),
                include: path.resolve(__dirname, 'src'),
                loader: "babel-loader"
            },
            {test: /\.jsx?$/, loader: 'babel-loader'},
            {test: /\.css$/, loader: 'style-loader!css-loader'},
            {test: /\.less$/, loader: 'style-loader!css-loader!less-loader'},
            {test: /\.(png|jpg)$/, loader: 'url-loader?limit=25000'}
        ]
    }
}
