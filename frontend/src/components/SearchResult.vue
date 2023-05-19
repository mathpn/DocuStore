<template>
    <div class="search-result">
        <span v-if="this.expanded" @click="toggleExpandResult" class="search-result-title">
            <b>Title: </b>{{ this.title }}
            <br>
            <b>Score: </b>{{ Math.round(this.score * 100) / 100 }}
        </span>
        <span v-else @click="toggleExpandResult" class="search-result-title">
            {{ shortenTitle(this.title) }}
        </span>
        <button class="search-result-button" @click="openDocument">Open</button>
    </div>
</template>

<script>
import { BrowserOpenURL } from "../../wailsjs/runtime";

export default {
    data() {
        return {
            expanded: false,
            shortTitleLimit: 50,
        }
    },
    props: ['title', 'score', 'identifier', 'type', 'rawContent'],
    methods: {
        shortenTitle() {
            const words = this.title.split(" ");
            let length = 0;
            for (let i = 0; i < words.length; i++) {
                const word = words[i];
                length += word.length
                if (length > this.shortTitleLimit) {
                    return words.slice(0, Math.max(1, i - 1)).join(" ") + " (...)"
                }
            }
            return this.title
        },
        toggleExpandResult() {
            this.expanded = !this.expanded;
        },
        openDocument() {
            if (this.Type == 0) {
                BrowserOpenURL(this.identifier);
            } else {
                console.log(this.identifier);
                this.$parent.$emit('markdown-content', this.rawContent);
                this.$parent.$emit('global-component', 'markdown');
            }
        },
    },
}
</script>

<style scoped>
.search-result {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    padding: 10px;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    background-color: #ffffff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.search-result a {
    color: #ffffff;
    text-decoration: none;
}

.search-result button:hover {
    background-color: #136063;
}

.search-result-title {
    margin-right: 10px;
    cursor: pointer;
}

.search-result-button {
    color: #ffffff;
    background-color: #169ba0;
    border: none;
    border-radius: 4px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    font-size: 11pt;
    display: inline-block;
    padding: 10px 20px;
    background-color: #169ba0;
    color: #fff;
}
</style>