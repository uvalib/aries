<template>
  <div class="list-wrapper">
    <h4>Aries Repositories<span @click="closeRepoList" class="hide-repos">Close</span></h4>  
    <table>
      <tr>
        <th>Service</th><th>URL</th><th>Alive</th>
      </tr>
      <tr v-for="repo in repositories"
          v-bind:key="repo.name">
        <td>{{repo.name}}</td>
        <td>{{repo.url}}</td>
        <td>{{repo.alive}}</td>
      </tr>
    </table>
  </div>
</template>

<script>
import EventBus from '@/EventBus'
export default {
  name: "RepositoryList",
  props: {
    repositories: {
      type: Array,
      required: true
    }
  },
  methods: { 
    closeRepoList: function() {
      EventBus.$emit("close-repos-clicked")
    }
  }
}
</script>

<style scoped>
h4 {
  margin: 4px 0;
  position: relative;
}
.hide-repos {
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
.hide-repos:hover {
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
