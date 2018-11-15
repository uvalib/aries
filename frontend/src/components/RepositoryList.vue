<template>
  <div class="list-wrapper">
    <h4>Aries Repositories</h4>  
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
export default {
  name: "RepositoryList",
  props: {
    repositories: {
      type: Array,
      required: true
    }
  }
};
</script>

<style scoped>
h4 {
  margin: 4px 0;
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
