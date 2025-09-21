const path = require('path');
const BUID_DIR = path.resolve("../dist");
const Dotenv = require('dotenv-webpack');
const HtmlWebpackPlugin = require("html-webpack-plugin");


module.exports = {
    entry: {
        callback: path.resolve(__dirname, './auth/Callback.tsx'),
        logout: path.resolve(__dirname, './auth/Logout.tsx'),
        home: path.resolve(__dirname, './home/index.tsx'),
        budget: path.resolve(__dirname, './budget/index.js'),
        account: path.resolve(__dirname, './account/index.tsx')
    },
    resolve: {
        extensions: [".js", ".jsx", ".ts", ".tsx"]
    },
    module: {
        rules: [
            {
                test: [/\.tsx?$/, /\.ts?$/,],
                use: ['ts-loader'],
                exclude: /node_modules$/,
            },
            {
                test: [/\.js?$/],
                exclude: path.resolve(__dirname, "node_modules"),
                use: {
                    loader: "babel-loader",
                    options: {
                        presets: ['@babel/env', '@babel/react']
                    }
                }

            }
        ]
    },
    plugins: [
        new Dotenv({
            path: `../environments/.env.${process.env.ENV}`
        }),

        new HtmlWebpackPlugin({
            title: 'OnlyOne Portal',
            filename: 'logout.html',
            chunks: ["logout"],
            template: "template.html"
        }),
        new HtmlWebpackPlugin({
            title: 'OnlyOne Portal',
            filename: 'callback.html',
            chunks: ["callback"],
            template: "template.html"
        }),

        new HtmlWebpackPlugin({
            title: 'OnlyOne Portal',
            filename: 'index.html',
            chunks: ["home"],
            template: "template.html"
        }),
        new HtmlWebpackPlugin({
            title: 'OnlyOne Portal',
            filename: 'budget/index.html',
            chunks: ["budget"],
            template: "template.html"
        }),
        new HtmlWebpackPlugin({
            title: 'OnlyOne Portal',
            filename: 'account/index.html',
            chunks: ["account"],
            template: "template.html"
        })
    ],
    output: {
        filename: '[name].[fullhash].js',
        path: BUID_DIR
    }
}