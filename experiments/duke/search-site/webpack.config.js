// the html plugin will dynamically add the bundle script tags to the main index.html file
// it also allows us to use template to build the rest of that file
var HtmlWebpackPlugin = require('html-webpack-plugin')
var webpack = require('webpack')

var environment = process.env.NODE_ENV || 'development'
console.log(environment)

//var config = require('dotenv').config({path: __dirname + '/.env.'+ environment})
var toml = require('toml')
var fs = require('fs')
const config = toml.parse(fs.readFileSync('./config.toml', 'utf-8'));

module.exports = {
  // start an main.js and follow requires to build the 'app' bundle in the 'dist' directory
  entry: {
    app: ['babel-polyfill', "./src/main.js"]
  },
  // put all built files in dist
  // use 'name' variable to make 
  // bundles named after the entryoints
  // above
  output: {
    path: __dirname + "/dist/",
    filename: "[name].js",
    library: "VivoSearch",
    libraryTarget: "umd"
  },
  module: {
    rules: [
      { test: /\.(ttf|eot|svg|woff(2)?)(\?[a-z0-9]+)?$/, loader: 'file-loader' },
      // style pre-processing
      { test: /\.less$/, loader: 'style-loader!css-loader!less-loader' }, // use ! to chain loaders
      { test: /\.css$/, loader: 'style-loader!css-loader' },
      { test: /\.(scss)$/,
        use: [
          {
            // Adds CSS to the DOM by injecting a `<style>` tag
            loader: 'style-loader'
          },
          {
            // Interprets `@import` and `url()` like `import/require()` and will resolve them
            loader: 'css-loader'
          },
          {
            // Loader for webpack to process CSS with PostCSS
            loader: 'postcss-loader',
            options: {
              plugins: function () {
                return [
                  require('autoprefixer')
                ];
              }
            }
          },
          {
            // Loads a SASS/SCSS file and compiles it to CSS
            loader: 'sass-loader'
          }
        ]
      },
      { test: /\.(png|gif|jpg)$/, loader: 'file-loader' },
      { test: /jquery/, loader: 'expose-loader?$!expose-loader?jQuery' },
      { test: /\.json$/, loader: 'json' },
       // react/jsx and es6/2015 transpiling
      {
        test: /\.js$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
        // http://asdfsafds.blogspot.com/2016/02/referenceerror-regeneratorruntime-is.html
        // http://stackoverflow.com/questions/33527653/babel-6-regeneratorruntime-is-not-defined-with-async-await
        query: {
          presets: ['react','es2015'],
          //plugins: ["transform-runtime"]
        }
      }
    ]
  },
  // make sourcemaps in separate files
  devtool: 'source-map',
  plugins: [
    new HtmlWebpackPlugin({
      inject: 'head',
      hash: true,
      title: "Vivo Search",
      template: 'src/index.ejs/'
    }),
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(environment),
      'process.env.ELASTIC_URL': JSON.stringify(config.elastic.url),
      //'process.env.ORG_URL':  JSON.stringify(process.env.ORG_URL)
    })
  ]
}

