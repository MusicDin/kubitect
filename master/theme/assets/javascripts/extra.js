/*
 * Landing page terminal
 */

// initial url path
let prevPath = window.location.pathname;

// trigger terminal animation when window loads
window.addEventListener('DOMContentLoaded', terminalAnimation)

window.onload = function () {
    let bodyList = document.querySelector("body")

    let observer = new MutationObserver(() => {

        // check if paths differ
        if (prevPath !== window.location.pathname) {

            // update href to match current location
            prevPath = window.location.pathname;

            if (window.location.pathname === '/') {
                terminalAnimation()
            }
        }
    })

    // start observing body for mutations
    observer.observe(bodyList, { childList: true, subtree: true });
}

// terminal animation
function terminalAnimation() {

    // milliseconds between each output line printed
    const outputDelay = 50

    // milliseconds before command is applied
    const applyCommandDelay = 250

    // milliseconds before command is typed
    const startCommandDelay = 1000

    // milliseconds between each command character typed
    const commandCharDelay = 20

    const Output = Symbol("output")
    const Command = Symbol("command")

    const content = [
        { type: Command, value: "curl -o kubitect.tar.gz -L https://dl.kubitect.io/linux/amd64/latest" },
        { type: Command, value: "tar -xzf kubitect.tar.gz" },
        { type: Command, value: "sudo mv kubitect /usr/local/bin/" },
        { type: Command, value: "" },
        { type: Command, value: "kubitect export preset --name getting-started > cluster.yaml" },
        { type: Command, value: "kubitect apply --config cluster.yaml" },
        { type: Output, value: "" },
        { type: Output, value: "Ensuring SSH keys are present..." },
        { type: Output, value: "Ensuring Terraform is installed..." },
        { type: Output, value: "Terraform could not be found locally." },
        { type: Output, value: "Installing Terraform in '/home/user/.kubitect/share/terraform'..." },
        { type: Output, value: "Initializing modules..." },
        { type: Output, value: "" },
        { type: Output, value: "Initializing the backend..." },
        { type: Output, value: "" },
        { type: Output, value: "Initializing provider plugins..." },
        { type: Output, value: "" },
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
        value += '<span class="terminal-command">' + command + '</span>'

        return value
    }

    // wraps the output into span element
    function wrapOutput(line) {

        let value = "";

        // add command element
        value += '<span class="terminal-output">' + line + '</span><br>'

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

                    await delay(startCommandDelay)
                    await typeSequence(cmdElement, line.value)
                    await delay(applyCommandDelay)

                    // remove cursor when command is "applied"
                    cmdElement.classList.remove("terminal-cursor")

                    target.innerHTML += "<br>"
                    break

                case Output:

                    target.innerHTML += wrapOutput(line.value)
                    await delay(outputDelay)
                    break
            }
        }
    }

    // sets placeholder (transparent) content
    const setPlaceholder = async(target) => {

        let placeholder = ""

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

    // viewport width lower than 768px --> mobile --> wait for scroll
    const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
    if (vw < 768) {
        document.getElementById('main-box').addEventListener('scroll', scrollEvent)
    } else {
        printContent(terminalContent)
    }
}
