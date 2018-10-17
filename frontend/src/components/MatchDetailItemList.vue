<template>
  <span v-if="items">
    <ol v-if="items.length > 1">
      <li v-for="(item, index) in items" :key="index">
        <a :href="item.url" target="_blank">{{ item.url }}</a>
      </li>
    </ol>
    <template v-else>
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
