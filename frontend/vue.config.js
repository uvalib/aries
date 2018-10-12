// NOTES on this are found here:
//    https://cli.vuejs.org/config/#devserver
//    https://github.com/chimurai/http-proxy-middleware#proxycontext-config
module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: process.env.ARIES_API, // or 'http://localhost:8095',
        changeOrigin: true,
        logLevel: 'debug'
      }
    }
  },
  configureWebpack: {
    performance: {
      // bump max sizes to 512k
      maxEntrypointSize: 512000,
      maxAssetSize: 512000
    }
  }
}
