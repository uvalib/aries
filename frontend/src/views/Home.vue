<template>
  <div class="home main-content">
    <div class="search-panel">
      <input id="target-id" type="text" placeholder="Search all repositories">
      <button class="pure-button">Search</button>
      <p class="instructions">
        Enter an identifer in the box above and hit search to find all occurrences in the UVA Library repositories
      </p>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'home',
  }
</script>

<style scoped>
  div.search-panel {
    width: 70%;
    margin: 2% auto;
    padding: 25px;
    text-align: center;
  }
  div.search-panel .instructions {
    width:50%;
    margin:20px auto;
    color: #666;
  }
  #target-id {
    margin: 10px;
    width: 50%;
    padding: 5px 10px;
    border-radius: 2px;
    outline: none;
    border: 1px solid #ccc;
  }
</style>
