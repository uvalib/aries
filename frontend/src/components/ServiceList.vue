<template>
  <LoadingSpinner message="Loading services information..." v-if="loading"/>
  <div v-else class="list-wrapper">
    <h4>Aries Services<span @click="closeRepoList" class="hide-services">Close</span></h4>  
    <table>
      <tr>
        <th>Service</th><th>URL</th><th>Alive</th>
      </tr>
      <tr v-for="svc in services"
          v-bind:key="svc.name">
        <td>{{svc.name}}</td>
        <td>{{svc.url}}</td>
        <td>{{svc.alive}}</td>
      </tr>
    </table>
  </div>
</template>

<script>
import EventBus from '@/EventBus'
import LoadingSpinner from '@/components/LoadingSpinner'

export default {
  name: "ServiceList",

  components: {
    LoadingSpinner
  },

  computed: {
    services: function() {
        return this.$store.getters.services
    },
    loading: function() {
        return this.$store.getters.isLoading
    }
  },

  created: function () {
    this.$store.dispatch('getServices')
  },

  methods: { 
    closeRepoList: function() {
      EventBus.$emit("close-services-clicked")
    }
  }
}
</script>

<style scoped>
h4 {
  margin: 4px 0;
  position: relative;
}
.hide-services {
  position: absolute;
  font-weight: 200;
  font-size: 0.8em;
  background: #efefef;
  padding: 1px 9px;
  border-radius: 10px;
  right: 0;
  opacity: 0.6;
  cursor: pointer;
}
.hide-services:hover {
  opacity: 1;
}
.list-wrapper {
  width: 50%;
  margin: 0 auto;
}
table {
  width: 100%;
  margin: 0;
  border-collapse: collapse;
}
table td,
table th {
  text-align: left;
  padding: 3px 9px;
}
table th {
  background-color: #f5f5f5;
  border-bottom: 1px solid #ccc;
  border-top: 1px solid #ccc;
  color: #666;
}
</style>
