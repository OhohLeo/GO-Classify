var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var CopyWebpackPlugin = require('copy-webpack-plugin');
var ProvidePlugin = require('webpack/lib/ProvidePlugin');
var helpers = require('./helpers');


module.exports = {
    entry: {
        'polyfills': './src/polyfills.ts',
        'vendor': './src/vendor.ts',
        'app': './src/main.ts'
    },

    resolve: {
        extensions: ['', '.js', '.ts']
    },

    module: {
        loaders: [
            {
                test: /\.ts$/,
                loader: 'ts'
            }
        ]
    },

    plugins: [
        new webpack.optimize.CommonsChunkPlugin({
            name: ['app', 'vendor', 'polyfills']
        }),

        new HtmlWebpackPlugin({
            template: 'src/index.html'
        }),

        new ProvidePlugin({
            jQuery: 'jquery',
            $: 'jquery',
            jquery: 'jquery'
        }),

        // Copy all files
        new CopyWebpackPlugin([

            // Copy all css
            {
                context: 'src/css',
                from: '**/*.css',
                to: 'css'
            },

            // Copy all fonts
            {
                context: 'src/fonts',
                from: '**/*.otf',
                to: 'fonts'
            },

            // Copy manifest.json
            { from: 'manifest.json' },

            // Copy all templates && css and img
            {
                context: 'src/app',
                from: '**/*.html',
                to: 'app',
            },
            {
                context: 'src/app',
                from: '**/*.css',
                to: 'app',
            },
            {
                context: 'src/app',
                from: '**/*.png',
                to: 'app',
            },

            // Copy all images
            {
                context: 'src/img',
                from: '**/*.png',
                to: 'img',
            },
            {
                context: 'src/img',
                from: '**/*.svg',
                to: 'img',
            },

            // Copy all translations
            {
                context: 'src/i18n',
                from: '**/*.json',
                to: 'i18n',
            }
        ]),
    ]
};
