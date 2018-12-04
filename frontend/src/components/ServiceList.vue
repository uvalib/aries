<template>
  <div class="list-wrapper">
    <h4>Aries Services<span @click="closeRepoList" class="hide-services">Close</span></h4>  
    <table>
      <tr>
        <th style="text-align:center"><span @click="addServiceClick" class="icon add"></span></th><th>Service</th><th>URL</th><th style="text-align:center">Status</th>
      </tr>
      <tr v-for="svc in services"
          v-bind:key="svc.id">
        <template v-if="editing(svc.id)">
          <td colspan="2">
            <input id="name-edit" :value="editService.name"/>
          </td>
          <td >
            <input id="url-edit" :value="editService.url"/>
          </td>
          <td>
            <span class="icon accept" @click="acceptEditClick"></span>
            <span class="icon cancel" @click="cancelEditClick"></span>
          </td>
        </template> 
        <template v-else>
          <td style="text-align: center;"><span @click="editServiceClick" class="icon edit" :data-id="svc.id"></span></td>
          <td>{{svc.name}}</td>
          <td>{{svc.url}}</td>
          <td style="text-align:center">
            <span v-if="svc.alive"  class="indicator online"></span>
            <span v-else class="indicator offline"></span>
          </td>
        </template>
      </tr>
      <tr v-if="addingService"> 
        <td colspan="2">
            <input id="name-edit" placeholder="New Service Name"/>
          </td>
          <td >
            <input id="url-edit" placeholder="New Service URL"/>
          </td>
          <td>
            <span class="icon accept" @click="acceptAddClick"></span>
            <span class="icon cancel" @click="cancelAddClick"></span>
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

  data: function () {
    return {
        editService: null,
        addingService: false
    }
  },

  computed: {
    services: function() {
      return this.$store.getters.services
    }
  },
  
  methods: { 
    editing: function(id) {
      return this.editService != null && this.editService.id == id
    },

    closeRepoList: function() {
      this.$root.$emit("close-services-clicked")
    },

    editServiceClick: function(event) {
      let btn = event.target
      this.editService = this.$store.getters.getServiceByID( btn.dataset.id )
    },

    cancelEditClick: function() {
      this.editService = null
    },

    acceptEditClick: function() {
      let input = document.getElementById('name-edit')
      this.editService.name = input.value
      input = document.getElementById('url-edit')
      this.editService.url = input.value
      this.$store.dispatch('updateService', this.editService)
      this.editService = null
    },

    addServiceClick: function() {
      this.addingService = true
    },

    acceptAddClick: function() {
      let input = document.getElementById('name-edit')
      let name = input.value
      input = document.getElementById('url-edit')
      let url = input.value
      this.$store.dispatch('addService', {name: name, url: url})
      this.addingService = false
    },

    cancelAddClick: function() {
      this.addingService = false
    },
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
.icon {
  display: inline-block;
  width: 20px;
  height: 20px;
  opacity: 0.5;
  cursor: pointer;
  vertical-align: inherit;
}
.icon:hover {
  opacity: 1;
}
.icon.edit {
  width: 16px;
  height: 16px;
  background-image: url(../assets/edit.png);
}
.icon.accept {
  background-image: url(../assets/accept.png);
  margin: 0 4px 0 0;
}
.icon.cancel {
  background-image: url(../assets/cancel.png);
  margin: 0 0 0 4px;
}
.icon.add {
  background-image: url(../assets/add.png);
}
input {
  width: 100%;
  box-sizing: border-box;
  font-size: 0.85em;
  border-radius: 4px;
  border: 1px solid #ccc;
  padding: 2px 0 2px 8px;
}
</style>
