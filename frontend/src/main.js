import './style.css';

import {Search} from '../wailsjs/go/main/App';
import {AddContent} from '../wailsjs/go/main/App';

const maxLength = 20000;
var inputBox = document.getElementById('input-box');
var counter = document.getElementById("char-count")
var successPopup = document.getElementById('success-popup');
var slideBar = document.getElementById('slide-bar');
document.getElementById("register-popup-content").classList.add("not-loading");
document.getElementById("register-popup-loading").classList.add("not-loading");
document.getElementById("search-button").addEventListener("click", callSearch);
document.getElementById("content-button").addEventListener("click", addDoc);

// Set the initial placeholder
inputBox.innerHTML = inputBox.dataset.placeholder;

inputBox.addEventListener('focus', function() {
    if (inputBox.innerHTML == inputBox.dataset.placeholder) {
        inputBox.innerHTML = '';
    }
});

inputBox.addEventListener('blur', function() {
    if (inputBox.innerHTML == '') {
        inputBox.innerHTML = inputBox.dataset.placeholder;
    }
});

inputBox.addEventListener("input", () => {
    const text = inputBox.innerText;
    counter.innerHTML = text.length + '/' + maxLength;
    if (text.length > maxLength) {
        // truncate the text to the maximum size
        inputBox.innerText = text.substring(0, maxLength);
    }
});

function sleep (time) {
    return new Promise((resolve) => setTimeout(resolve, time));
}

function addDoc() {
    var txt = document.getElementById("input-box").innerText;
    if (txt == inputBox.dataset.placeholder) {
        return
    }
    console.log("Adding document: " + txt);
    // call add
    document.getElementById("register-popup-content").classList.remove("not-loading");
    document.getElementById("register-popup-loading").classList.remove("not-loading");
    AddContent(txt).then(() => {
        document.getElementById("register-popup-content").classList.add("not-loading");
        document.getElementById("register-popup-loading").classList.add("not-loading");
        // Show the success popup
        successPopup.style.display = 'block';
        
        // Start the slide bar animation
        slideBar.style.animationPlayState = 'running';
        // Hide the success popup and reset the slide bar animation after 5 seconds
        setTimeout(function() {
            successPopup.style.display = 'none';
            slideBar.style.animation = 'none';
            slideBar.offsetHeight; // Force a reflow to restart the animation from the beginning
            slideBar.style.animation = 'slideOut 1s linear forwards';
        }, 1000);
    });
}

function callSearch() {
    var searchTerm = document.getElementById("search-box").value;
    console.log("Searching for: " + searchTerm);
    Search(searchTerm).then(results => {
        var searchResults = document.getElementById("search-results");
        searchResults.innerHTML = "";
        for (var i = 0; i < results.length; i++) {
            var result = results[i];
            var container = document.createElement("div");
            container.className = "search-result";
            var title = document.createElement("span");
            title.className = "search-result-title";
            title.textContent = result.Title;
            var link = document.createElement("a");
            console.log(result);
            link.href = result.Identifier; // FIXME
            link.textContent = "Open";
            link.className = "search-result-button";
            container.appendChild(title);
            container.appendChild(link);
            searchResults.appendChild(container);
        }
    })
}