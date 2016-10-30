var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var CopyWebpackPlugin = require('copy-webpack-plugin');
var ProvidePlugin = require('webpack/lib/ProvidePlugin');
var helpers = require('./helpers');


module.exports = {
    entry: {
        'polyfills': './app/polyfills.ts',
        'vendor': './app/vendor.ts',
        'app': './app/main.ts'
    },

    resolve: {
        extensions: ['', '.js', '.ts']
    },

    module: {
	loaders: [
	    {
		test: /\.ts$/,
		loaders: ['awesome-typescript-loader', 'angular2-template-loader']
	    },
	    {
		test: /\.html$/,
		loader: 'html'
	    },
	    {
		test: /\.(png|jpe?g|gif|svg|woff|woff2|ttf|eot|ico)$/,
		loader: 'file?name=assets/[name].[hash].[ext]'
	    },
	    {
		test: /\.css$/,
		exclude: helpers.root('app'),
		loader: ExtractTextPlugin.extract('style', 'css?sourceMap')
	    },
	    {
		test: /\.css$/,
		include: helpers.root('app'),
		loader: 'raw'
	    }
	]
    },

    plugins: [
        new webpack.optimize.CommonsChunkPlugin({
            name: ['app', 'vendor', 'polyfills']
        }),

        new HtmlWebpackPlugin({
            template: 'app/index.html'
        }),

        new ProvidePlugin({
            jQuery: 'jquery',
            $: 'jquery',
            jquery: 'jquery',
            'window.$': 'jquery',
            'window.jQuery': 'jquery',
            'root.jQuery': 'jquery',
            Hammer : 'hammerjs/hammer'
        })
    ]
};
