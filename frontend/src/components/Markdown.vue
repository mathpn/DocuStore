<template>
    <button class="back-button" @click="goBack">Go back</button>
    <div class="markdown-content" v-html="parsedMarkdown"></div>
</template>

<script>
import DOMPurify from 'dompurify';
import { marked } from 'marked';

export default {
    props: ['content'],
    methods: {
        goBack() {
            this.$emit('global-component', 'search');
        },
    },
    computed: {
        parsedMarkdown() {
            return DOMPurify.sanitize(marked.parse(this.content, { mangle: false }))
        }
    }
}
</script>

<style>
.back-button {
    position: relative;
    width: 15%;
    padding-top: 10px;
    padding-bottom: 10px;
    font-size: 16px;
    font-weight: bold;
    color: #ffffff;
    background-color: #169ba0;
    border: none;
    border-radius: 4px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    cursor: pointer;
}
</style>