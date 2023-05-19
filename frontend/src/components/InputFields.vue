<template>
    <div class="search-bar">
        <div v-if="error" class="error-popup">
            <p><b>Error: {{ errorMsg }}</b></p>
            <div id="slide-bar"></div>
        </div>
        <textarea type="text" class="text-input" ref="input-box" id="input-box" rows="1"
            @input="resizeTextarea(); limitInput();" placeholder="Please enter a URL or raw text" v-model="input"
            :disabled="addingData" />
        <div id="char-count">{{ charCount }}/{{ maxChars }}</div>
        <div v-if="addingData" id="content-button" class="search-button">
            <div class="loader"></div>
        </div>
        <button v-else id="content-button" class="search-button" @click="addInput">Register</button>
    </div>
    <div class="search-bar">
        <input v-debounce:300ms="doSearch" @keydown.enter="doSearch" @input="resetIsSearched" type="text"
            class="search-input" id="search-box" placeholder="Search" v-model="searchField" />
    </div>
</template>

<script>
import { Search } from '../../wailsjs/go/main/App';
import { AddContent } from '../../wailsjs/go/main/App';
import { vue3Debounce } from 'vue-debounce';

export default {
    data() {
        return {
            maxChars: 20000,
            input: '',
            searchField: '',
            isSearched: false,
            addingData: false,
            errorMsg: '',
            error: false,
        }
    },
    methods: {
        limitInput() {
            if (this.input.length > this.maxChars) {
                // truncate the text to the maximum size
                this.input = this.input.substring(0, this.maxChars);
            }
        },
        shortenTitle(title) {
            const words = title.split(" ");
            let length = 0;
            for (let i = 0; i < words.length; i++) {
                const word = words[i];
                length += word.length
                if (length > 50) {
                    return words.slice(0, Math.max(1, i - 1)).join(" ") + " (...)"
                }
            }
            return title
        },
        doSearch() {
            if (this.isSearched | this.searchField === '') {
                return
            };
            this.isSearched = true;
            console.log("searching", this.searchField);
            Search(this.searchField)
                .then(
                    results => {
                        results.forEach(result => result.expanded = false);
                        this.$emit('search-results', results);
                    })
                .catch(err => {
                    console.log("doSearch failed: ", err);
                    this.errorMsg = err;
                    this.error = true;
                    setTimeout(() => this.error = false, 1000);
                })
        },
        resetIsSearched() {
            this.isSearched = false;
        },
        addInput() {
            this.addingData = true;
            if (this.input === '') {
                this.addingData = false;
                return
            }
            console.log("adding", this.input);
            AddContent(this.input)
                .then(() => {
                    this.input = '';
                })
                .catch(err => {
                    console.log("addInput failed: ", err);
                    this.errorMsg = err;
                    this.error = true;
                    setTimeout(() => this.error = false, 1000);
                })
                .finally(() => this.addingData = false);
        },
        resizeTextarea() {
            let element = this.$refs["input-box"];
            element.style.height = "20px";
            element.style.height = (element.scrollHeight - 20) + "px";
        },
    },
    computed: {
        charCount() {
            return this.input.length;
        }
    },
    directives: {
        debounce: vue3Debounce({ lock: true }),
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
    position: absolute;
    bottom: 0px;
    right: 0px;
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

.text-input {
    width: 80%;
    padding: 10px;
    font-size: 16px;
    border-radius: 4px;
    border: none;
    color: #757575;
    background-color: #f2f2f2;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    position: relative;
    resize: none;
    overflow: hidden;
    font-variant: monospace;
    font-size: 1rem;
    color: #000;
    min-height: 76px;
}

#char-count {
    position: absolute;
    bottom: 5px;
    right: 18%;
    font-size: 12px;
    color: #bebebe;
}

.error-popup {
    border: none;
    border-radius: 4px;
    position: fixed;
    top: 50%;
    left: 50%;
    color: #521414;
    transform: translate(-50%, -50%);
    background-color: #ffa29f;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    padding: 20px;
    z-index: 999;
}

#slide-bar {
    border: none;
    border-radius: 4px;
    position: relative;
    height: 5px;
    background-color: #521414;
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
    width: 10px;
    height: 10px;
    animation: spin 1s linear infinite;
    margin: 0 auto;
}

@keyframes spin {
    0% {
        transform: rotate(0deg);
    }

    100% {
        transform: rotate(360deg);
    }
}
</style>