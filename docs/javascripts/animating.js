// initial url
var oldPath = window.location.pathname;

/**/
window.onload = function () {

    var bodyList = document.querySelector("body")

    var observer = new MutationObserver((mutations) => {

        //check if paths differ
        if (oldPath !== window.location.pathname) {

            // update href to match current location
            oldPath = window.location.pathname; 

            // set up cli animation
            cli()
        }
        
    });

    // start observing body for mutations
    observer.observe(bodyList, {childList: true, subtree: true}); 
};
/**/
window.addEventListener('DOMContentLoaded', cli)
/** 
window.addEventListener('load', (event) => {
    cli()
    //document.querySelector("#btn_learn-more").addEventListener("load", cli())
});
/** */

// trigger terminal animation when window loads
function cli() {

    const appendLine = async (target,line) => {
        // console.log(target)        

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

    if(terminal == null) {
        return
    }

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

    // animation 
    const animate = () => {
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
    }

    const animated = false
    const scrollEvent = () => {
        var element = document.getElementById('terminal');
        var position = element.getBoundingClientRect();
    
        // checking for partial visibility
        if(position.top < window.innerHeight && position.bottom >= 0) {
            animate()
            document.getElementById('main-box').removeEventListener('scroll',scrollEvent)
        }
    }

    //vieport width lower than 768px --> phone --> wait for scroll
    const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
    if(vw<768){
        document.getElementById('main-box').addEventListener('scroll',scrollEvent)
    }else{
        animate()
    }

    
    
}