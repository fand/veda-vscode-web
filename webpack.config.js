const path = require('path');

module.exports = {
  mode: 'development',
  entry: path.resolve(__dirname, 'src/content.ts'),
  output: {
    path: path.resolve(__dirname, 'build'),
    filename: 'content.js',
  },
  resolve: {
    extensions: ['.ts', '.js'],
  },
  module: {
    rules: [{ test: /\.ts$/, loader: 'ts-loader' }],
  },
};
