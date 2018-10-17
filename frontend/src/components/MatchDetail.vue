<template>
  <div class="result">
    <p class="hit-header">Match found in {{ this.match.system }} <span class="elapsed">{{match.response_time_ms}}ms</span></p>
    <div class="pure-g result-detail">
      <div class="pure-u-1-4 label"><span>Identifiers:</span></div>
      <div class="pure-u-3-4 data">{{ identifiers }}</div>

      <template v-if="publicURLs">
        <div class="pure-u-1-4 label"><span>Public URL:</span></div>
        <div class="pure-u-3-4 data">
          <MatchDetailURLList :urls="publicURLs"/>
        </div>
      </template>

      <template v-if="adminURLs">
        <div class="pure-u-1-4 label"><span>Admin URL:</span></div>
        <div class="pure-u-3-4 data">
          <MatchDetailURLList :urls="adminURLs"/>
        </div>
      </template>

      <template v-if="serviceURLs">
        <div class="pure-u-1-4 label"><span>Service URL:</span></div>
        <div class="pure-u-3-4 data">
          <MatchDetailItemList :items="serviceURLs"/>
        </div>
      </template>

      <template v-if="metadataURLs">
        <div class="pure-u-1-4 label"><span>Metadata URL:</span></div>
        <div class="pure-u-3-4 data">
          <MatchDetailItemList :items="metadataURLs"/>
        </div>
      </template>

      <template v-if="masterFiles">
        <div class="pure-u-1-4 label"><span>Master File:</span></div>
        <div class="pure-u-3-4 data">
          <MatchDetailFileList :files="masterFiles"/>
        </div >
      </template>

      <template v-if="derivatives">
        <div class="pure-u-1-4 label"><span>Derivative File:</span></div>
        <div class="pure-u-3-4 data">
          <MatchDetailFileList :files="derivatives"/>
        </div>
      </template>

      <template v-if="match.response.access_restriction">
        <div class="pure-u-1-4 label"><span>Access Restriction:</span></div>
        <div class="pure-u-3-4 data" v-html="accessRestriction"/>
      </template>
    </div>
  </div>
</template>

<script>
  import MatchDetailURLList from '@/components/MatchDetailURLList'
  import MatchDetailFileList from '@/components/MatchDetailFileList'
  import MatchDetailItemList from '@/components/MatchDetailItemList'

  export default {
    components: {
      MatchDetailURLList,
      MatchDetailItemList,
      MatchDetailFileList
    },
    props: {
      match: Object
    },
    computed: {
      renderNA: function() {
        return "<span style='color:#aaa;font-style: italic'>N/A</span>"
      },
      identifiers: function() {
        return this.match.response.identifier.join(", ")
      },
      adminURLs: function() {
        return this.match.response.administrative_url
      },
      publicURLs: function() {
        return this.match.response.access_url
      },
      metadataURLs: function() {
        return this.match.response.metadata_url
      },
      serviceURLs: function() {
        return this.match.response.service_url
      },
      accessRestriction: function() {
        let resp = this.match.response
        if (resp.access_restriction) {
          return "<span>"+resp.access_restriction+"</span>"
        } else {
          return this.renderNA
        }
      },
      derivatives: function() {
        return this.match.response.derivative_file
      },
      masterFiles: function() {
        return this.match.response.master_file
      }
    }
  }
</script>

<style scoped>
  ol {
    margin: 0;
    padding-inline-start: 20px;
  }
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
  span.na {
    color:#aaa;
    font-style: italic;
  }
  a {
    color: #6495ed;
    text-decoration: none;
    font-weight: 500;
  }
  a:hover {
   text-decoration: underline;
 }
</style>
