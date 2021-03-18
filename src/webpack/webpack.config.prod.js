// vim: ts=2 sw=2 et
const Webpack = require('webpack');

const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');

// extend webpack.common.js
module.exports = merge(common, {
  mode: 'production',
  devtool: 'source-map',
  stats: 'errors-only',
  bail: true,
});
