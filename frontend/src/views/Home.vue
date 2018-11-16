<template>
  <div class="home main-content">
    <div v-if="searching" class="search-panel">
      <h4>Searching...</h4>
      <img src="../assets/spinner2.gif"/>
    </div>
    <div v-else class="search-panel">
      <h4>Identifier Search</h4>
      <div class="pure-button-group" role="group">
        <input id="target-id" ref="target-id" type="text" @keyup.enter="searchClicked" placeholder="Search all repositories" :value="searchTerm">
        <button id="search-btn" class="pure-button"  @click="searchClicked">
          <i class="fas fa-search"></i>
          Search
        </button>
        </div>
      <div v-if="errorMsg">
        <h4 class>Search Failed!</h4>
        <h4 class="error">{{ this.errorMsg }}</h4>
      </div>
      <div v-else>
        <template v-if="searchTerm">
          <p class="instructions">
            <span @click="showRepos" title="View Repositories" class="view-repo"><b>{{ repoCount }}</b> repositories</span>
             searched in <b>{{searchTime}}ms</b><br/> Matches: <b>{{ hits }}</b>
          </p>
          <RepositoryList v-if="showRepoList"/>
          <MatchDetail
            v-for="hit in matches"
            v-bind:key="hit.system"
            v-bind:match="hit">
          </MatchDetail>
        </template>
        <template v-else>
          <p class="instructions">
            Enter an identifer in the box above and hit search to find all<br/>occurrences in the UVA Library repositories
          </p>
          <p class="instructions" v-if="!showRepoList">
            Aries will search
            <span @click="showRepos" title="View Repositories" class="view-repo"><b>{{ repoCount }}</b> repositories</span>
          </p>
          <RepositoryList v-if="showRepoList"/>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
  import MatchDetail from '@/components/MatchDetail'
  import RepositoryList from '@/components/RepositoryList'
  import EventBus from '@/EventBus'
  import axios from 'axios'

  export default {
    name: 'home',
    components: {
      MatchDetail,
      RepositoryList
    },

    data: function () {
      return {
        repoCount: 0,
        searching: false,
        searchTerm: "",
        matches: [],
        searchTime: 0,
        errorMsg: "",
        showRepoList: false
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
        this.repoCount = response.data.length
      })
    },

    mounted: function (){
      EventBus.$on("close-repos-clicked", this.handleCloseRepoClicked)
    },

    methods: {
      handleCloseRepoClicked: function() {
        this.showRepoList = false
      },

      showRepos: function() {
        this.showRepoList = true
      },

      searchClicked: function() {
        if (this.searching === true) return
        this.searchTerm = this.$refs["target-id"].value.trim()
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
    margin: 2% auto;
    text-align: center;
  }
  h4 {
    width: 50%;
    margin: 10px auto 10px auto;
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
    width: 50%;
    border-radius: 2px;
    outline: none;
    border: 1px solid #ccc;
    padding: 7px 14px;
    border-right: 0;
    border-radius: 20px 0 0 20px;
  }
  #search-btn {
    padding: 8px 14px;
    border-radius: 0 20px 20px 0;
    color: #666;
    font-weight: 500;
  }
  .view-repo {
    font-weight: 500;
     color: cornflowerblue;
     cursor:pointer
  }
  .view-repo:hover {
    text-decoration: underline
  }
</style>
