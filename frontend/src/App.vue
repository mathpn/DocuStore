<script>
import Search from './components/Search.vue';
import Markdown from './components/Markdown.vue';

export default {
  name: 'DocuStore',
  data() {
    return {
      globalComponent: 'search',
      markdownContent: '',
    }
  },
  methods: {
    changeComponent(componentName) {
      this.globalComponent = componentName
    },
    changeContent(content) {
      this.markdownContent = content;
    }
  },
  components: {
    Search,
    Markdown,
  }
}
</script>

<template>
  <Markdown v-on:global-component="changeComponent" :content="this.markdownContent"
    v-if="this.globalComponent === 'markdown'"></Markdown>
  <Search v-on:markdown-content="changeContent" v-on:global-component="changeComponent" v-else></Search>
</template>

<style scoped></style>
