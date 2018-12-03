import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

// root state object. Holds all of the state for the system
const state = {
  services: [],
  error: null
}

// state getter functions. All are functions that take state as the first param 
// and the getters themselves as the second param. Getter params are passed 
// as a function. Access as a property like: this.$store.getters.NAME
const getters = {
  services: state => {
    return state.services
  },

  error: state => {
    return state.error
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
  setServices (state, services) {
    state.services = services
  },

  setError (state, error) {
    state.error = error
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
    axios.get("/api/services").then((response)  =>  {
      ctx.commit('setServices', response.data )
    }).catch((error) => {
      ctx.commit('setServices', []) 
      ctx.commit('setError', "Internal Error: Unable to reach any services") 
    })
  },

  updateService( ctx, updatedService ) {
    axios.post("/api/services", updatedService).then((response)  =>  {
      ctx.commit('updateService', updatedService )
    }).catch( error => {
      ctx.commit('setError', "Update Failed: "+ error.response.data) 
    })
  }
}

// A Vuex instance is created by combining state, getters, actions and mutations
const store = new Vuex.Store({
  state,
  getters,
  actions,
  mutations
})
export default store