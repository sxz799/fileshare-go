const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  //publicPath: "/dist", //加上这一行即可
  devServer: {
    port: 4001,
    proxy: {
      '/api': {
        ws: false,
        target: process.env.VUE_APP_BASE_URL,
        changeOrigin: true
      }
    }
  },
  pwa: {
    iconPaths: {
      favicon32: 'favicon.ico',
      favicon16: 'favicon.ico',
      appleTouchIcon: 'appleIcon.png',
      maskIcon: 'favicon.ico',
      msTileImage: 'favicon.ico'
    }
  }
})
