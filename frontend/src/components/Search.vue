<template>
    <div class="container">
        <InputFields v-on:search-results="addResults" />
        <div class="search-results" id="search-results">
            <search-result
                v-for="result in searchResults"
                :docID="result.DocID"
                :title="result.Title"
                :score="result.Score"
                :type="result.Type"
                :identifier="result.Identifier"
                :key="result.DocID"
            ></search-result>
        </div>
    </div>
</template>

<script>
import InputFields from "./InputFields.vue";
import SearchResult from "./SearchResult.vue";

export default {
    data() {
        return {
            searchResults: [],
        }
    },
    mounted() {
        window.scrollTo(0,0);
    },
    methods: {
        addResults(results) {
            this.searchResults = results
        },
    },
    components: {
        InputFields,
        SearchResult,
    }
}
</script>

<style scoped>
.search-results {
    display: flex;
    flex-direction: column;
    margin-top: 20px;
}
</style>