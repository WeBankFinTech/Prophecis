module.exports = {
  presets: [
    '@vue/app',
    ['@babel/preset-env', // 添加 babel-preset-env 配置
      {
        'modules': false
      }
    ]
  ],
  plugins: [ // element官方教程
    [
      'component',
      {
        'libraryName': 'element-ui',
        'styleLibraryName': 'theme-chalk'
      }
    ]
  ]
}
