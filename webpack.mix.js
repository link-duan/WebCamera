let mix = require('laravel-mix');

mix.webpackConfig({
  externals: {
    'element-ui': 'Element',
    'axios': 'axios',
    'vue': 'Vue',
    'vue-router': 'VueRouter',
    'lodash':'_',
  },
});

mix.js('assets/js/app.js', 'app/view/js/app.js');