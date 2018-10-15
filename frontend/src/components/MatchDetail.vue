<template>
  <div class="result">
    <p class="hit-header">Match found in {{ this.match.system }} <span class="elapsed">{{match.response_time_ms}}ms</span></p>
    <div class="pure-g result-detail">
      <div class="pure-u-1-4 label"><span>Identifiers:</span></div>
      <div class="pure-u-3-4 data">{{ identifiers }}</div>

      <div class="pure-u-1-4 label"><span>Public URL:</span></div>
      <div class="pure-u-3-4 data" v-html="publicURL"/>

      <div class="pure-u-1-4 label"><span>Admin URL:</span></div>
      <div class="pure-u-3-4 data"><a :href="adminURL" target="_blank">{{ adminURL }}</a></div>

      <div class="pure-u-1-4 label"><span>Service URL:</span></div>
      <div class="pure-u-3-4 data"></div>

      <div class="pure-u-1-4 label"><span>Metadata URL:</span></div>
      <div class="pure-u-3-4 data"></div>

      <div class="pure-u-1-4 label"><span>Master File:</span></div>
      <div class="pure-u-3-4 data"></div>

      <div class="pure-u-1-4 label"><span>Derivitaive File:</span></div>
      <div class="pure-u-3-4 data"></div>

      <div class="pure-u-1-4 label"><span>Access Restriction:</span></div>
      <div class="pure-u-3-4 data"></div>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'match-detail',
    props: {
      match: Object
    },
    computed: {
      identifiers: function() {
        return this.match.response.identifier.join(", ")
      },
      adminURL: function() {
        return this.match.response.administrative_url[0]
      },
      publicURL: function() {
        let resp = this.match.response
        if ( !resp.access_url || resp.access_url.length == 0) {
          return "<span style='color:#aaa;font-style: italic'>N/A</span>"
        }
        return "FOO"
      }
    }
  }
</script>

<style scoped>
  div.result {
    width: 75%;
    margin: 0 auto;
    text-align: left;
    font-size: 0.9em;
  }
  div.row {
    margin-left: 15px;
  }
  div.result p.hit-header {
    padding: 4px 12px;
    background: #E57200;
    color: white;
    font-weight: bold;
    text-align: left;
  }
  span.elapsed {
    float: right;
    font-size: 0.8em;
    font-style: italic;
    margin-top: 2px;
  }
  div.result-detail {
    margin-left: 15px;
    color: #666;
  }
  div.result-detail div {
    margin: 5px 0;
  }
  div.label {
    font-weight: bold;
    text-align: right;
  }
  div.label span {
    padding-right: 10px;
  }
  a {
    color: #6495ed;
    text-decoration: none;
    font-weight: 700;
  }
  a:hover {
   text-decoration: underline;
 }
</style>
