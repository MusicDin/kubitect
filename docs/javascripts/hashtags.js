/*
 * Remove code block hashtags before annotations
 */

// initial url
var oldHref = document.location.href;

// remove hastags when content is loaded for the first time
window.addEventListener('DOMContentLoaded', removeHashtags)

// waits untill annotations are loaded and remove hashtags
function removeHashtags() {

    document.querySelectorAll("code").forEach((node) => {
        let str = node.innerHTML
        node.innerHTML = str.replace(/(# )(\([0-9]*\))/g, "$2")
    }) 

};

window.onload = function () {

    var bodyList = document.querySelector("body")

    var observer = new MutationObserver((mutations) => {

        mutations.forEach(() => {

            // if user changed location
            if (oldHref != document.location.href) {

                // update href to match current location
                oldHref = document.location.href; 

                // wait for annotations and remove hashtags
                removeHashtags()
            }
        });
    });

    // start observing body for mutations
    observer.observe(bodyList, {childList: true, subtree: true}); 
};
