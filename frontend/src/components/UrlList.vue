<template>
  <span v-if="urls">
    <ol v-if="urls.length > 1">
      <li v-for="(item, index) in urls" :key="index">
        <a :href="getURL(item)" target="_blank">{{ getURL(item) }}</a>
      </li>
    </ol>
    <a v-else :href="getURL(urls[0])" target="_blank">{{ getURL(urls[0]) }}</a>
  </span>
  <span v-else class="na">N/A</span>
</template>

<script>
  export default {
    name: 'url-list',
    props: {
      urls: Array
    },
    methods: {
      getURL: function(item) {
        if ( typeof  item == "string" ) {
          return item
        }
        return item.url
      }
    }
  }
</script>

<style scoped>
  ol {
    margin: 0;
    padding-inline-start: 20px;
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
