<template>
  <span v-if="files">
    <ol v-if="files.length > 1">
      <li v-for="(file, index) in files" :key="index">
        {{ file }}
      </li>
    </ol>
    <span v-else>{{ files[0] }}</span>
  </span>
  <span v-else class="na">N/A</span>
</template>

<script>
  export default {
    name: 'file-list',
    props: {
      files: Array
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
</style>
