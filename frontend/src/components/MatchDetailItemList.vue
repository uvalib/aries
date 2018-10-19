<template>
  <span v-if="items">
    <ol v-if="items.length > 1">
      <li v-for="(item, index) in items" :key="index">
        <span class="url-type">{{ getItemType(item).value }} :</span>
        <a :href="item.url" target="_blank">{{ item.url }}</a>
      </li>
    </ol>
    <template v-else>
      <span class="url-type">{{ getItemType(items[0]).value }} :</span>
      <a :href="items[0].url" target="_blank">{{ items[0].url }}</a>
    </template>
  </span>
  <span v-else class="na">N/A</span>
</template>

<script>
  export default {
    props: {
      items: {
        type: Array,
        required: true
      }
    },
    methods: {
      getItemType: function(item) {
        for (var key in item) {
          if (key !== "url") {
            return {name: key, value: item[key]}
          }
        }
      }
    }
  }
</script>

<style scoped>
  span.url-type {
    margin-right: 10px;
    color: #999;
  }
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
