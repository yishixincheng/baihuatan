import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    islogin:false,
    token:'',
    user:{}
  },
  mutations: {

    login_succ(state,obj){

      state.token=obj.token;
      state.islogin=true;
      state.user=obj.user;
      console.log(state);

    }

  },
  actions: {
  },
  modules: {
  }
})