<template>
  <div id="app">
    <aries-header/>
    <router-view/>
    <aries-footer/>
  </div>
</template>


<script>
  import AriesHeader from './components/AriesHeader'
  import AriesFooter from './components/AriesFooter'

  export default {
    components: {
      'aries-header': AriesHeader,
      'aries-footer': AriesFooter
    }
  }
</script>


<style>
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
  }
</style>
