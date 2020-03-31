const path = require("path");
const webpack = require("webpack");
const CopyWebpackPlugin = require("copy-webpack-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");

function resolve(dir) {
  return path.join(__dirname, "..", dir);
}

module.exports = {
  target: "node",
  context: resolve("src"),
  entry: "./module.ts",
  output: {
    filename: "module.js",
    path: resolve("dist"),
    libraryTarget: "amd"
  },
  externals: [
    "jquery",
    "lodash",
    "moment",
    "angular",
    function(context, request, callback) {
      var prefix = "grafana/";
      if (request.indexOf(prefix) === 0) {
        return callback(null, request.substr(prefix.length));
      }
      callback();
    }
  ],
  plugins: [
    new webpack.optimize.OccurrenceOrderPlugin(),
    new CopyWebpackPlugin([
      { from: "../LICENSE.txt" },
      { from: "../README.md" },
      { from: "plugin.json" },
      { from: "img/*" },
      { from: "partials/*" }
    ]),
    new CleanWebpackPlugin({
      cleanStaleWebpackAssets: false,
      cleanAfterEveryBuildPatterns: [
        "!README.md",
        "!LICENSE",
        "!plugin.json",
        "!img/*",
        "!partials/*",
        "!simple-json-plugin_linux_amd64",
        "!simple-json-plugin_windows_amd64.exe"
      ]
    })
  ],
  resolve: {
    extensions: [".ts", ".tsx", ".js"]
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        loaders: ["ts-loader"],
        exclude: /node_modules/
      },
      {
        test: /\.css$/,
        use: ["style-loader", "css-loader"]
      }
    ]
  }
};
