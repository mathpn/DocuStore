<template>
    <div class="container">
        <InputFields v-on:search-results="addResults" />
        <Pagination v-if="searchResults.length > this.resultsPerPage" :current-page="page" :total-pages="totalPages"
            @page-changed="changePage" />
        <div class="search-results" id="search-results">
            <search-result v-for="result in pageResults" :docID="result.DocID" :title="result.Title" :score="result.Score"
                :type="result.Type" :identifier="result.Identifier" :key="result.DocID"></search-result>
        </div>
    </div>
</template>

<script>
import InputFields from "./InputFields.vue";
import SearchResult from "./SearchResult.vue";
import Pagination from "./Pagination.vue";

export default {
    data() {
        return {
            searchResults: [],
            pageResults: [],
            page: 1,
            resultsPerPage: 8,
        };
    },
    mounted() {
        window.scrollTo(0, 0);
    },
    methods: {
        addResults(results) {
            this.searchResults = results;
            this.pageResults = results.slice(0, this.resultsPerPage);
        },
        changePage(page) {
            this.page = page;
            this.pageResults = this.searchResults.slice((page - 1) * this.resultsPerPage, page * this.resultsPerPage);
        },
    },
    computed: {
        totalPages() {
            return Math.ceil(this.searchResults.length / this.resultsPerPage);
        },
    },
    components: {
        InputFields,
        SearchResult,
        Pagination,
    },
};
</script>

<style scoped>
.search-results {
    display: flex;
    flex-direction: column;
    margin-top: 0px;
}
</style>
