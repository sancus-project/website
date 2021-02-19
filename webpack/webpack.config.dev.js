// vim: ts=2 sw=2 et
const Path = require('path');
const Webpack = require('webpack');

const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');

// extend webpack.common.js
module.exports = merge(common, {
  mode: 'development',
  devtool: 'cheap-source-map',
  // webpack serve
  devServer: {
    inline: true,
    index: '',
    contentBase: Path.join(__dirname, '../assets'),
    host: process.env.HOST || '127.0.0.1',
    port: process.env.PORT,
    // backend
    proxy: {
      context: () => true,
      target: {
        host: '127.0.0.1',
        port: process.env.BACKEND
      }
    }
  }
});
