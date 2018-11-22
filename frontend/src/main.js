import Vue from 'vue'
import App from './App.vue'
import router from './router'
import VueResource from 'vue-resource'

Vue.use(VueResource);

Vue.config.productionTip = false
Vue.http.options.root = 'http://localhost:3000/v1'

new Vue({
  router,
  render: function (h) { return h(App) }
}).$mount('#app')
