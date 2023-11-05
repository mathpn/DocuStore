<script>
import { ReadTextFile } from '../wailsjs/go/main/App';
import ErrorModal from './components/ErrorModal.vue';
import Search from './components/Search.vue';
import Markdown from './components/Markdown.vue';

export default {
  name: 'DocuStore',
  data() {
    return {
      search: true,
      textContent: '',
      loaded: false,
      error: false,
      errorMsg: '',
    }
  },
  methods: {
    showSearch(show) {
      this.search = show;
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
          setTimeout(() => this.error = false, 2000);
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
  <Markdown v-on:show-search="showSearch" :content="this.textContent"
    v-if="this.loaded & !this.search"></Markdown>
  <Search v-on:markdown-doc-id="loadText" v-on:show-search="showSearch" v-show="search"></Search>
</template>

<style scoped></style>
