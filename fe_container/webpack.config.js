import HtmlWebpackPlugin from 'html-webpack-plugin'
import path from 'path'
import webpack from 'webpack'
import dotenv from 'dotenv'

const env = process.env.NODE_ENV || 'dev'
const Dotenv = dotenv.config({
  path: `./.env.${env}`,
})

const __dirname = path.resolve()

const port =
  Number(
    process.env.BLOCK_ENV_URL_container.substr(
      process.env.BLOCK_ENV_URL_container.length - 4,
    ),
  ) || 3011

export default {
  entry: './src/index',
  mode: 'development',
  resolve: {
    alias: {
      core: path.join(__dirname, 'core'),
    },
  },
  devServer: {
    historyApiFallback: true,
    hot: true,
    static: path.join(__dirname, 'dist'),
    port,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS',
      'Access-Control-Allow-Headers':
        'X-Requested-With, content-type, Authorization',
    },
  },
  target: 'web',
  externals: {
    env: JSON.stringify(process.env),
  },
  output: {
    publicPath: '/',
    crossOriginLoading: 'anonymous',
  },
  module: {
    rules: [
      {
        test: /.js$/,
        loader: 'babel-loader',
        options: {
          presets: ['@babel/preset-react'],
        },
      },
      {
        test: /\.(jpg|png|svg|eot|svg|ttf|woff|woff2|ico|gif)$/,
        use: {
          loader: 'url-loader',
        },
      },
      {
        test: /\.s[ac]ss$/i,
        use: [
          'style-loader',
          'css-loader',
          {
            loader: 'sass-loader',
          },
        ],
      },
      {
        test: /\.css$/i,
        use: ['style-loader', 'css-loader'],
      },
      {
        test: /.m?js/,
        type: 'javascript/auto',
      },
      {
        test: /.m?js/,
        resolve: {
          fullySpecified: false,
        },
      },
    ],
  },
  optimization: {
    minimize: false,
  },
  plugins: [
    new webpack.DefinePlugin({
      process: { env: JSON.stringify(Dotenv.parsed) },
    }),
    new HtmlWebpackPlugin({
      template: './public/index.html',
      favicon: './public/favicon.ico',
    }),
  ],
}
