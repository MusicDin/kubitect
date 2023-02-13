/*
 * Remove code block hashtags before annotations
 */

// initial url path
var oldHref = document.location.href;
var prevPath = window.location.pathname;

// remove hastags when content is loaded for the first time
window.addEventListener('DOMContentLoaded', removeHashtags)

// trigger terminal animation when window loads
window.addEventListener('DOMContentLoaded', terminalAnimation)

window.onload = function () {

    var bodyList = document.querySelector("body")

    var observer = new MutationObserver((mutations) => {

        // check if paths differ
        if (prevPath !== window.location.pathname) {

            // update href to match current location
            prevPath = window.location.pathname;

            if (window.location.pathname === '/') {

                // trigger terminal animation on landing page
                terminalAnimation()

            } else {
                mutations.forEach(() => {
                    if (oldHref != document.location.href) {

                        // update href to match current location
                        oldHref = document.location.href;

                        // wait for annotations and remove hashtags
                        removeHashtags()
                    }
                });
            }
        }
    });

    // start observing body for mutations
    observer.observe(bodyList, { childList: true, subtree: true });
};

// waits until code block annotations are loaded and remove hashtags
function removeHashtags() {

    document.querySelectorAll("code").forEach((node) => {
        let str = node.innerHTML
        node.innerHTML = str.replace(/(# )(\([0-9]*\))/g, "$2")
    })
};

/*
 * Landing page terminal
 */

// terminal animation
function terminalAnimation() {

    // miliseconds between each output line printed
    const outputDelay = 50

    // miliseconds between each command typed
    const commandDelay = 500

    // miliseconds between each command character typed
    const commandCharDelay = 20
    
    const Output = Symbol("output")
    const Command = Symbol("command")

    const content = [
        { type: Command, value: "curl -o kubitect.tar.gz -L https://github.com/MusicDin/kubitect/releases/..." },
        { type: Command, value: "tar -xzf kubitect.tar.gz" },
        { type: Command, value: "sudo mv kubitect /usr/local/bin/" },
        { type: Command, value: "kubitect apply" },
        { type: Output, value: "" },
        { type: Output, value: "Preparing cluster 'default'..." },
        { type: Output, value: "Setting up 'main' virtual environment..." },
        { type: Output, value: "Creating virtual environment..." },
        { type: Output, value: "Installing pip3 dependencies..." },
        { type: Output, value: "This can take up to a minute when the virtual environment is initialized for the first time...<br>" },
        { type: Output, value: "PLAY [localhost]<br>" },
        { type: Output, value: "TASK [cluster-config/copy : Make sure config directory exists]" },
        { type: Output, value: "<span style=\"color:green\">ok: [127.0.0.1]</span><br>" },
        { type: Output, value: "..." }
    ]

    const delay = async (ms) => new Promise(resolve => setTimeout(resolve, ms));

    // simulate typing by appending characters to the target's innerHtml
    const typeSequence = async (target, sequence) => {

        for (const char of sequence) {

            target.innerHTML += char
            await delay(commandCharDelay)
        }
    }

    // wraps the command into span element and adds "$"" sign in front of it
    function wrapCommand(command) {

        let value = "";

        // add dollar sign
        value += '<span class="terminal-command-dollar-sign">$</span>'

        // add command element
        value += '<span class=\"terminal-command\">' + command + '</span>'

        return value
    }

    // print terminal content one by one line
    const printContent = async (target) => {

        for (const line of content) {
            switch (line.type) {

                case Command:

                    // add empty command element
                    target.innerHTML += wrapCommand("")
                    
                    // get added command element
                    let cmdElement = target.lastChild

                    // add cursor when writing command
                    cmdElement.classList.add("terminal-cursor")
                    
                    await delay(commandDelay)
                    await typeSequence(cmdElement, line.value)
                    await delay(commandDelay)

                    // remove cursor when command is "applied"
                    cmdElement.classList.remove("terminal-cursor")

                    target.innerHTML += "<br>"
                    break

                case Output:
                    
                    target.innerHTML += line.value + "<br>"
                    await delay(outputDelay)
                    break
            }
        }
    }

    // sets placeholder (transparent) content
    const setPlaceholder = async(target) => {

        placeholder = ""

        for (const line of content) {

            if (line.type == Command) {
                placeholder += wrapCommand(line.value)
            } else {
                placeholder += line.value
            }

            placeholder += "<br>"
        }

        target.innerHTML += placeholder
    }
    
    // event that is triggered on scroll
    const scrollEvent = () => {
        
        let terminal = document.getElementById('terminal');
        let position = terminal.getBoundingClientRect();
        
        // check for partial visibility
        if (position.top < window.innerHeight && position.bottom >= 0) {
            printContent(terminalContent)
            document.getElementById('main-box').removeEventListener('scroll', scrollEvent)
        }
    }
    
    let terminalContent = document.getElementById('terminal-content');
    let terminalPlaceholder = document.getElementById('terminal-placeholder')
    
    // Prevent animation if terminal placeholder element is not found
    if (terminalPlaceholder == null) {
        return
    }

    setPlaceholder(terminalPlaceholder)

    // vieport width lower than 768px --> mobile --> wait for scroll
    const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
    if (vw < 768) {
        document.getElementById('main-box').addEventListener('scroll', scrollEvent)
    } else {
        printContent(terminalContent)
    }
}