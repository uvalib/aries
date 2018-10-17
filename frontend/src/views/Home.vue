<template>
  <div class="home main-content">
    <div v-if="searching" class="search-panel">
      <h4>Searching...</h4>
      <img src="../assets/spinner2.gif"/>
    </div>
    <div v-else class="search-panel">
      <input id="target-id" ref="target-id" type="text" @keyup.enter="searchClicked" placeholder="Search all repositories" :value="searchTerm">
      <button class="pure-button"  @click="searchClicked">Search</button>
      <div v-if="errorMsg">
        <h4 class>Search Failed!</h4>
        <h4 class="error">{{ this.errorMsg }}</h4>
      </div>
      <div v-else>
        <template v-if="searchTerm">
          <p class="instructions">
            <b>{{ repoCount }}</b> systems searched in <b>{{searchTime}}ms</b><br/> Matches: <b>{{ hits }}</b>
          </p>
          <match-detail
            v-for="hit in matches"
            v-bind:key="hit.system"
            v-bind:match="hit">
          </match-detail>
        </template>
        <template v-else>
          <p class="instructions">
            Enter an identifer in the box above and hit search to find all occurrences in the UVA Library repositories
          </p>
          <p class="instructions">
            Aries searches <b>{{ repoCount }}</b> UVA Library repositories
          </p>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
  import MatchDetail from '@/components/MatchDetail'
  import axios from 'axios'

  export default {
    name: 'home',
    components: {
      'match-detail': MatchDetail
    },

    data: function () {
      return {
        repositories: [],
        searching: false,
        searchTerm: "",
        matches: [],
        searchTime: 0,
        errorMsg: ""
      }
    },

    computed: {
      repoCount: function() {
        return this.repositories.length
      },
      hasResults: function() {
        return this.matches.length > 0
      },
      hits: function() {
        return this.matches.length
      }
    },

    created: function () {
      this.repositories = []
      axios.get("/api/services").then((response)  =>  {
        this.repositories = response.data
      })
    },

    methods: {
      searchClicked: function() {
        if (this.searching === true) return
        this.searchTerm = this.$refs["target-id"].value
        if ( this.searchTerm.length === 0) return

        this.matches = []
        this.searching = true
        this.errorMsg = ""
        axios.get("/api/resources/"+this.searchTerm).then((response)  =>  {
          this.searchTime = response.data.total_response_time_ms
          for (let i=0; i<response.data.responses.length; i++) {
            let resp = response.data.responses[i]
            if (resp.status == 200) {
              this.matches.push(resp)
            }
          }
        }).catch((error) => {
          if (error.response ) {
            this.errorMsg =  error.response.data
          } else {
            this.errorMsg =  error
          }
          this.matches = []
        }).finally(() => {
          this.searching = false
        })
      }
    }
  }
</script>

<style scoped>
  div.search-panel {
    width: 70%;
    margin: 2% auto;
    text-align: center;
  }
  h4 {
    width: 50%;
    margin: 10px auto 0 auto;
    color: #666
  }
  h4.error {
    margin-top: 5px;
    color: #922;
    font-weight: 200;
    font-style: italic;
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
