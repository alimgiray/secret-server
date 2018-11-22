import Vue from 'vue'
import Router from 'vue-router'
import Create from './components/Create.vue'
import Retrieve from './components/Retrieve.vue'

Vue.use(Router)

export default new Router({
  routes: [
    { path: '/', component: Create },
    { path: '/retrieve', component: Retrieve }
  ]
})
