<script>
import { ReadTextFile } from '../wailsjs/go/main/App';
import ErrorModal from './components/ErrorModal.vue';
import Search from './components/Search.vue';
import Markdown from './components/Markdown.vue';

export default {
  name: 'DocuStore',
  data() {
    return {
      globalComponent: 'search',
      textContent: '',
      loaded: false,
      error: false,
      errorMsg: '',
    }
  },
  methods: {
    changeComponent(componentName) {
      this.globalComponent = componentName
    },
    loadText(docID) {
      this.loaded = false;
      ReadTextFile(docID)
        .then(c => {
          this.textContent = c;
          this.loaded = true;
        })
        .catch(err => {
          console.log("loadText failed: ", err);
          this.errorMsg = err;
          this.error = true;
          setTimeout(() => this.error = false, 1000);
        })
    }
  },
  components: {
    ErrorModal,
    Markdown,
    Search,
  }
}
</script>

<template>
  <ErrorModal v-if="error" :errorMsg="errorMsg"></ErrorModal>
  <Markdown v-on:global-component="changeComponent" :content="this.textContent"
    v-if="this.loaded & this.globalComponent === 'markdown'"></Markdown>
  <Search v-on:markdown-doc-id="loadText" v-on:global-component="changeComponent" v-else></Search>
</template>

<style scoped></style>
