// trigger terminal animation when window loads
window.addEventListener('load', (event) => {
    
    const appendLine = async (target,line) => {
        target.innerHTML += line+'<br>'
    }
    
    // recursively appends given lines to the given target's innerHtml (uses appendLine() function)
    const appendMultipleLines = async (target,lines,speed,callback)=>{
        if(lines.length===0)callback()
        else{
            appendLine(target,lines[0])
            setTimeout(()=>{
                appendMultipleLines(target,lines.slice(1),speed,()=>{
                    callback()
                })
            },speed)
        }
    }
    
    // append command to the target's innterHtml (adds a styled '$' at the beggining)
    // uses typeSequence to simulate typing the command
    const printCommand = async (target, command, callback) => {
        target.innerHTML += '<span class="command-dollar-sign">$</span> '
        typeSequence(target,command,20,()=>{
            target.innerHTML += '<br>'
            callback()
        })
    }
    
    //smoothly appends characters of the given sequence to given target's innerHtml
    const typeSequence = async (target, sequence, speed, callback) => {
        if(sequence.length===0)callback()
        else{
            target.innerHTML += sequence[0]
            setTimeout(()=>{
                typeSequence(target,sequence.substr(1),speed,()=>{
                    callback()
                })
            },speed)
        }
    }
    
    // miliseconds between each command typed
    let timeBetween = 500

    // get the terminal nativeElement
    let terminal = document.getElementById('command_line')

    // output lines to be appended at the end of all 3 commands
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

    // start animation
    setTimeout(()=>{
        printCommand(terminal,'curl -o kubitect -L https://dl.kubitect.com',()=>{
            setTimeout(()=>{
                printCommand(terminal,'sudo mv kubitect /usr/local/bin/',()=>{
                    setTimeout(()=>{
                        printCommand(terminal,'kubitect apply',()=>{
                            setTimeout(()=>{
                                appendMultipleLines(terminal,terminalOutputLines,50,()=>{
                                })
                            },timeBetween)
                        })
                    },timeBetween)
                })
            },timeBetween)
        })
    },timeBetween)
});

function learnMore(){
    console.log("hello")
    window.location.href = '/user-guide/installation'
}