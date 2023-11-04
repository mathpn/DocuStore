<template>
    <nav class="pagination">
        <ul>
            <li @click="changePage(1)">&hookleftarrow;</li>
            <li @click="changePage(currentPage - 1)">Prev</li>
            <li v-for="page in pages" :key="page" @click="changePage(page)" :class="{ active: currentPage === page }"
                class="page-button">
                {{ page }}
            </li>
            <li @click="changePage(currentPage + 1)">Next
            </li>
            <li @click="changePage(totalPages)">&hookrightarrow;</li>
        </ul>
    </nav>
</template>
  
<script>
export default {
    props: {
        currentPage: Number,
        totalPages: Number,
    },
    computed: {
        pages() {
            const maxVisiblePages = 5;
            const startPage = Math.max(1, this.currentPage - Math.floor(maxVisiblePages / 2));
            const endPage = Math.min(this.totalPages, startPage + maxVisiblePages - 1);

            return Array.from({ length: endPage - startPage + 1 }, (_, i) => startPage + i);
        },
    },
    methods: {
        changePage(page) {
            if (page >= 1 && page <= this.totalPages) {
                this.$emit('page-changed', page);
            }
        },
    },
};
</script>
  
<style>
.pagination {
    display: flex;
    justify-content: center;
}

.pagination ul {
    list-style: none;
    display: flex;
    padding: 0;
    margin-top: -5px;
}

.pagination li {
    cursor: pointer;
    padding: 0.25rem 0.5rem;
    margin: 0.1rem;
    border: 1px solid #ccc;
    border-radius: 5px;
    user-select: none;
}

.pagination .active {
    background-color: #169ba0;
    color: #fff;
}

.page-button {
    min-width: 20px;
    text-align: center;
}
</style>
