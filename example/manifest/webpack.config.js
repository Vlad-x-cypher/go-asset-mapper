const MiniCssExtractPlugin = require("mini-css-extract-plugin");

const { WebpackManifestPlugin } = require('webpack-manifest-plugin');
const path = require('path');

module.exports = {
    mode: 'production',
    entry: {
        bundle: './src/bundle.js',
    },
    output: {
        path: path.resolve(__dirname, 'public/bundle'),
        publicPath: '/bundle/',
        filename: '[name].[contenthash].js', // Main bundle with content hash
        chunkFilename: '[id].[contenthash].js', // Asynchronous chunks with content hash
        assetModuleFilename: '[name].[contenthash][ext][query]',
    },
    optimization: {
        runtimeChunk: { name: 'runtime' }
    },
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: [MiniCssExtractPlugin.loader, 'css-loader'],
            }
        ],
    },
    plugins: [
        new WebpackManifestPlugin(),
        new MiniCssExtractPlugin(),
    ]
}
