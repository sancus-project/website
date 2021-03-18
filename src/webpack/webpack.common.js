// vim: ts=2 sw=2 et
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const RemoveEmptyScriptsPlugin = require('webpack-remove-empty-scripts');
const Path = require('path');
const Webpack = require('webpack');

module.exports = {
  entry: {
    'js/index': Path.resolve(__dirname, '../../src/js/index.js'),
    'css/index': Path.resolve(__dirname, '../../src/scss/index.scss'),
  },
  output: {
    path: Path.join(__dirname, '../../assets'),
    filename: '[name].js',
  },
  plugins: [
    new RemoveEmptyScriptsPlugin(),
    new MiniCssExtractPlugin({
      filename: '[name].css',
    }),
  ],
  resolve: {
    alias: {
      '~': Path.resolve(__dirname, '../../src')
    },
  },
  module: {
    rules: [
      {
        test: /\.s?css$/,
        use: [MiniCssExtractPlugin.loader, 'css-loader', 'postcss-loader', 'sass-loader'],
      },
    ],
  },
}
