<template>
  <div id="app">
    <AriesHeader/>
    <router-view/>
    <AriesFooter/>
  </div>
</template>


<script>
  import AriesHeader from './components/AriesHeader'
  import AriesFooter from './components/AriesFooter'

  export default {
    components: {
      AriesHeader,
      AriesFooter
    }
  }
</script>


<style>
  html,body{
      height: 100%
  }
  div.main-content {
    padding-bottom: 210px;
  }
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    color: #2c3e50;
    min-height: 100%;
    position:relative;
  }
</style>
