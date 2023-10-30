<template>
    <InputModal v-if="showModal" v-on:show-modal="toggleModal" v-on:input-content="addText"
        :message="'Please provide a title:'"></InputModal>
    <div class="search-bar">
        <ErrorPopup v-if="error" :errorMsg="errorMsg"></ErrorPopup>
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
            class="search-input" id="search-box" ref="searchInput" placeholder="Search" v-model="searchField" />
    </div>
</template>

<script>
import ErrorPopup from './ErrorModal.vue';
import { Search } from '../../wailsjs/go/main/App';
import { AddURL } from '../../wailsjs/go/main/App';
import { AddText } from '../../wailsjs/go/main/App';
import { vue3Debounce } from 'vue-debounce';
import InputModal from './InputModal.vue';

const URLRegex = /^htt(p|ps):\/\/(.*)(\s|$)/i;

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
            showModal: false,
        }
    },
    components: {
        ErrorPopup,
        InputModal
    },
    mounted() {
        this.$refs.searchInput.focus();
    },
    methods: {
        limitInput() {
            if (this.input.length > this.maxChars) {
                // truncate the text to the maximum size
                this.input = this.input.substring(0, this.maxChars);
            }
        },
        updateTitle(title) {
            this.title = title;
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
                    setTimeout(() => this.error = false, 2000);
                })
        },
        resetIsSearched() {
            this.isSearched = false;
        },
        addInput() {
            const input = this.input.trim();
            if (input === '') {
                this.addingData = false;
                return
            }
            const type = URLRegex.test(input) ? 0 : 1;
            if (type === 1) {
                this.toggleModal(true);
                return
            } else {
                this.addURL()
            }
        },
        addURL() {
            const encodedInput = btoa(this.input); // base64 encoding
            this.addingData = true;
            const promise = AddURL(encodedInput);
            this.resolveAddPromise(promise);
        },
        addText(title) {
            this.addingData = true;
            const encodedInput = btoa(this.input); // base64 encoding
            const encodedTitle = btoa(title);
            const promise = AddText(encodedInput, encodedTitle);
            this.resolveAddPromise(promise);
        },
        resolveAddPromise(promise) {
            promise
                .then(() => {
                    this.input = '';
                })
                .catch(err => {
                    console.log("addInput failed: ", err);
                    this.errorMsg = err;
                    this.error = true;
                    setTimeout(() => this.error = false, 2000);
                })
                .finally(() => {
                    this.addingData = false;
                    this.resizeTextarea();
                });
        },
        resizeTextarea() {
            let element = this.$refs["input-box"];
            element.style.height = "20px";
            element.style.height = (element.scrollHeight - 20) + "px";
        },
        toggleModal(show) {
            if (show === true | show === false) {
                this.showModal = show;
            } else {
                this.showModal = !this.showModal;
            }
        }
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
