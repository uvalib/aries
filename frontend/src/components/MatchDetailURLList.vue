<template>
  <span v-if="urls">
    <ol v-if="urls.length > 1">
      <li v-for="(url, index) in urls" :key="index">
        <a :href="url" target="_blank">{{ url }}</a>
      </li>
    </ol>
    <a v-else :href="urls[0]" target="_blank">{{ urls[0] }}</a>
  </span>
  <span v-else class="na">N/A</span>
</template>

<script>
  export default {
    props: {
      urls: {
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
