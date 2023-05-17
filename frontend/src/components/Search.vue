<template>
    <div class="container">
        <div class="search-bar">
            <textarea v-on:keydown.enter="addInput" type="text" class="text-input" ref="input-box" id="input-box" rows="1" @input="resizeTextarea" placeholder="Please enter a URL or raw text" v-model="input"/>
            <div id="char-count"></div>
            <button id="content-button" class="search-button" v-on:click="addInput">Register</button>
        </div>
        <div class="search-bar">
            <input v-debounce:300ms="doSearch" v-on:keydown.enter="doSearch" type="text" class="search-input" id="search-box" placeholder="Search" v-model="searchField"/>
        </div>
        <div class="search-results" id="search-results" v-for="result in searchResults">
            <div class="search-result">
                <span class="search-result-title">{{ result.Title }}</span>
                <button class="search-result-button" v-on:click="openDocument(result)">Open</button>
            </div>
        </div>
    </div>
</template>

<script>
import {Search} from '../../wailsjs/go/main/App';
import {AddContent} from '../../wailsjs/go/main/App';
import { vue3Debounce } from 'vue-debounce'
import {BrowserOpenURL} from "../../wailsjs/runtime";

export default {
    data() {
        return {
            input: '',
            searchField: '',
            searchResults: [],
        }
    },
    methods: {
        doSearch() {
            console.log("searching", this.searchResults);
            Search(this.searchField).then(results => this.searchResults = results);
        },
        addInput() {
            console.log("adding", this.input);
            AddContent(this.input)
        },
        resizeTextarea(e) {
            let element = this.$refs["input-box"];
            element.style.height = "18px";
            element.style.height = element.scrollHeight + "px";
        },
        openDocument(document) {
            // TODO open text files as markdown
            BrowserOpenURL(document.Identifier);
        }
    },
    directives: {
    debounce: vue3Debounce({ lock: true })
  }
}
</script>

<style scoped>
.search-bar {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    position: relative;
}
.search-input {
    width: 100%;
    padding: 10px;
    font-size: 16px;
    border-radius: 4px;
    border: none;
    background-color: #f2f2f2;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
.search-button {
    width: 30%;
    padding: 10px;
    font-size: 16px;
    font-weight: bold;
    color: #ffffff;
    background-color: #169ba0;
    border: none;
    border-radius: 4px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    cursor: pointer;
}
.search-results {
    display: flex;
    flex-direction: column;
    margin-top: 20px;
}
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
.search-result a:hover {
    background-color: #136063;
}
.search-result-url {
    font-size: 14px;
    font-weight: bold;
    margin-right: 10px;
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
.search-result-button hover {
    background-color: #136063;
    font-weight: bold;
}
.text-input {
    width: 66%;
    padding: 10px;
    font-size: 16px;
    border-radius: 4px;
    border: none;
    color: #757575;
    background-color: #f2f2f2;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    /* height: auto; */
    position: relative;
    resize: none;
    overflow: hidden;
    font-variant: monospace;
    font-size: 1rem;
    color: #000;
    min-height: 72px;
}
#char-count {
    position: absolute;
    bottom: 5px;
    right: 32%;
    font-size: 12px;
    color: #bebebe;
}
#success-popup {
    display: none;
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background-color: #fff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    padding: 20px;
    z-index: 999;
}

#slide-bar {
    position: relative;
    height: 5px;
    background-color: #ccc;
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    z-index: 998;
    animation: slideOut 1s linear forwards;
}

@keyframes slideOut {
    from {
        width: 100%;
    }
    to {
        width: 0;
    }
}
#register-popup {
    position: fixed;
    top: 50%;
    left: 50%;
    height: 30px;
    transform: translate(-50%, -50%);
    background-color: #f7f7f7;
    z-index: 9999;
}

#register-popup-loading {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 10px;
    background-color: #f7f7f7;
}

.loader {
    border: 4px solid #f7f7f7;
    border-top: 4px solid #3498db;
    border-radius: 50%;
    width: 16px;
    height: 16px;
    animation: spin 1s linear infinite;
    margin: 0 auto;
}

@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}

.loading {
    display: block;
}

.not-loading {
    display: none;
}
</style>