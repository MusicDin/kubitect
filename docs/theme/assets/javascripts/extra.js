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

// terminal animation
function terminalAnimation() {

    const appendLine = async (target, line) => {
        target.innerHTML += line + '<br>'
    }

    // recursively appends given lines to the given target's innerHtml (uses appendLine() function)
    const appendMultipleLines = async (target, lines, speed) => {
        if (lines.length !== 0)
            appendLine(target, lines[0])
            setTimeout(() => {
                appendMultipleLines(target, lines.slice(1), speed)
            }, speed)
    }

    // append command to the target's innterHtml (adds a styled '$' at the beggining)
    // uses typeSequence to simulate typing the command
    const printCommand = async (target, command) => {
        target.innerHTML += '<span class="command-dollar-sign">$</span> '
        await typeSequence(target, command, 20)
        target.innerHTML += '<br>'
    }

    //smoothly appends characters of the given sequence to given target's innerHtml
    const typeSequence = async (target, sequence, speed) => {
        if (sequence.length > 0) {
            target.innerHTML += sequence[0]
            await delay(speed)
            await typeSequence(target, sequence.substr(1), speed)
        }
    }

    const delay = async (ms) => new Promise(resolve => setTimeout(resolve, ms));
    

    // miliseconds between each command typed
    let timeBetween = 500

    // get the terminal nativeElement
    let terminal = document.getElementById('command_line')

    // Prevent animation if terminal element is not found
    if (terminal == null)
        return
    
    // output lines to be appended at the end of commands
    const terminalOutputLines = [
        '',
        "Preparing cluster 'default'...",
        "Setting up 'main' virtual environment...",
        "Creating virtual environment...",
        "Installing pip3 dependencies...",
        "This can take up to a minute when the virtual environment is initialized for the first time...<br>",
        "PLAY [localhost]<br>",
        "TASK [cluster-config/copy : Make sure config directory exists]",
        "<span style=\"color:green\">ok: [127.0.0.1]</span><br>",
        "...",
    ]

    const animate = async () => {

        await printCommand(terminal, 'curl -o kubitect.tar.gz -L https://github.com/MusicDin/kubitect/releases/...')
        await delay(timeBetween)
        await printCommand(terminal, 'tar -xzf kubitect.tar.gz')
        await delay(timeBetween)
        await printCommand(terminal, 'sudo mv kubitect /usr/local/bin/')
        await delay(timeBetween)
        await printCommand(terminal, 'kubitect apply')
        await delay(timeBetween)

        await appendMultipleLines(terminal, terminalOutputLines, 50)

    }

    const animated = false
    const scrollEvent = () => {
        var element = document.getElementById('terminal');
        var position = element.getBoundingClientRect();

        // check for partial visibility
        if (position.top < window.innerHeight && position.bottom >= 0) {
            animate()
            document.getElementById('main-box').removeEventListener('scroll', scrollEvent)
        }
    }

    // vieport width lower than 768px --> phone --> wait for scroll
    const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
    if (vw < 768) {
        document.getElementById('main-box').addEventListener('scroll', scrollEvent)
    } else {
        animate()
    }
}