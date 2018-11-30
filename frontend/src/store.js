import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

// root state object. Holds all of the state for the system
const state = {
  services: [],
  loading: false
}

// state getter functions. All are functions that take state as the first param 
// and the getters themselves as the second param. Getter params are passed 
// as a function. Access as a property like: this.$store.getters.NAME
const getters = {
  isLoading: state => {
    return state.loading
  },

  services: state => {
    return state.services
  },

  getServiceByID: state => id => {
    return state.services.find(service => service.id === id)
  },
  
  serviceCount: state => {
    return state.services.length
  }
}

// Synchronous updates to the state. Can be called directly in components like this:
// this.$store.commit('mutation_name') or called from asynchronous actions
const mutations = {
  setLoading(state, isLoading) {
    state.loading = isLoading
  },

  setServices (state, services) {
    state.services = services
  },

  updateService (state, updatedSvc) {
    for (let idx = 0; idx < state.services.length; idx++) {
      let svc = state.services[idx]
      if ( svc.id === updatedSvc.id ) {
        state.services[idx] = updatedSvc
        break
      }
    }
  },

  addService (state, service) {
    state.services.push( service )
  },
}

// Actions are asynchronous calls that commit mutatations to the state.
// All actions get context as a param which is essentially the entirety of the 
// Vuex instance. It has access to all getters, setters and commit. They are 
// called from components like: this.$store.dispatch('action_name', data_object)
const actions = {
  getServices( ctx ) {
    ctx.commit('setLoading', true)
    axios.get("/api/services").then((response)  =>  {
      ctx.commit('setServices', response.data )
      ctx.commit('setLoading', false)
    })
  }
}

// A Vuex instance is created by combining state, getters, actions and mutations
export default new Vuex.Store({
  state,
  getters,
  actions,
  mutations
})