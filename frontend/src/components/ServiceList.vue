<template>
  <div class="list-wrapper">
    <h4>Aries Services<span @click="closeRepoList" class="hide-services">Close</span></h4>  
    <table>
      <tr>
        <th>Service</th><th>URL</th><th style="text-align:center">Status</th>
      </tr>
      <tr v-for="svc in services"
          v-bind:key="svc.id">
        <td>{{svc.name}}</td>
        <td>{{svc.url}}</td>
        <td style="text-align:center">
          <span v-if="svc.alive"  class="indicator online"></span>
          <span v-else class="indicator offline"></span>
        </td>
      </tr>
    </table>
  </div>
</template>

<script>

export default {
  name: "ServiceList",

  components: {
  },

  computed: {
    services: function() {
      return this.$store.getters.services
    }
  },

  methods: { 
    closeRepoList: function() {
      this.$root.$emit("close-services-clicked")
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
span.indicator {
  display: inline-block;
  width: 16px;
  height: 16px;
  border-radius: 15px;
}
span.indicator.online {
  background: #7c7;
}
span.offline {
  background: #c55;
}
</style>
