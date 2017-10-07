const path = require("path");
const HtmlWebpackPlugin = require('html-webpack-plugin');
const WebpackNotifierPlugin = require('webpack-notifier');

const isProd = process.env.NODE_ENV === 'production';

function getEntry() {
    const result = [];

    // the entry point of our app
    result.push('./src/index.js');

    return result;
}

function getSettings() {
    let settings = {};

    settings.apiUrl = isProd ? '/api' : 'http://localhost:3030/api';

    return settings;
}

function getPlugins() {
    const plugins = [
        new HtmlWebpackPlugin(({
            template: './src/index.html',
            settings: JSON.stringify(getSettings())
        }))
    ];

    if (process.env.NODE_ENV === 'development') {
        plugins.push(new WebpackNotifierPlugin());
    }

    return plugins;
}

const config = function (env) {
    let additionalElmFlags = '';

    if (process.env.NODE_ENV === 'development') {
        additionalElmFlags = '&debug=true'
    }

    return {
        entry: getEntry(),

        output: {
            path: path.resolve(__dirname + '/dist'),
            filename: 'app.js',
            publicPath: isProd ? '/assets/' : '/'
        },

        plugins: getPlugins(),

        module: {
            rules: [{
                test: /\.(css|scss)$/,
                use: [
                    'style-loader',
                    'css-loader',
                    'sass-loader',
                ]
            },
                {
                    test: /\.elm$/,
                    exclude: [/elm-stuff/, /node_modules/],
                    loader: 'elm-webpack-loader?verbose=true&warn=true' + additionalElmFlags,
                },
                {
                    test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
                    loader: 'url-loader?limit=10000&mimetype=application/font-woff',
                },
                {
                    test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
                    loader: 'file-loader',
                },
            ],

            noParse: /\.elm$/,
        },

        devServer: {
            inline: true,
            stats: {
                colors: true
            },
            historyApiFallback: {
                index: '/'
            }
        },


    };
};

module.exports = function (env) {
    if (!env)
        env = {};

    console.log('Node Env: ' + process.env.NODE_ENV);
    console.log(env);

    return config(env);
};
