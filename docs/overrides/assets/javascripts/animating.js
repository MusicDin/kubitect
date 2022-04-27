
const command1 = 'wget -O kubitect https://download.kubitect.io'
const commandResult1 = [
    '$'
]
const command2 = 'sudo mv kubitect /usr/local/bin/'
const commandResult2 = [
    '$'
]
const command3 = "kubitect apply"
const commandResult3 = [
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

window.addEventListener('load', (event) => {
    console.log('page is fully loaded');
    let command_top = document.getElementById('top-command-fg')
    let top_command_bg = document.getElementById('top-command-bg')

    let command_anim = document.getElementById('command-anim')
    let command_bg = document.getElementById('command-anim-bg')

    let command_result = document.getElementById('command-result-fg')
    let main_box = document.getElementById('main-box')

    let command_description = document.getElementById('command-description')

    let boundingRect = main_box.getBoundingClientRect()
    console.log(boundingRect.width)
    console.log(boundingRect.height)

    // console.log(command_anim)
    // console.log(command_result)
    // setTimeout(()=>{
    //     command_anim.innerHTML = commandText
    // },1000)

    
    const printLines = async () => {
        appendLine(command_result,'$ ')
        appendLine(command_result,'hello marko1')
        appendLine(command_result,'hello marko2')
        appendLine(command_result,'hello marko3')
    }
    
    const clearCommands =  () => {
        command_top.innerHTML = ''
        command_anim.innerHTML = ''
        command_result.innerHTML = ''
    }

    const loopCommands = async () => {
        clearCommands()
        //print first sequence
        top_command_bg.innerHTML = "$ "+command1
        command_bg.innerHTML = command1
        command_description.innerHTML = 'Downloads kubitect.'
        printSequence(command_top,command_anim,command1,20,()=>{
            console.log("printed first sequence")
            setTimeout(()=>{
                //print result

                printCommandResult(command_result,commandResult1,50,()=>{
                    console.log("printed result1")

                    setTimeout(()=>{
                        clearCommands()
                        
                        top_command_bg.innerHTML = command2
                        command_bg.innerHTML = command2
                        
                        command_description.innerHTML = 'Installs kubitect.'
                        printSequence(command_top,command_anim,command2,20,()=>{
                            console.log("printed second command")
                            setTimeout(()=>{
                                printCommandResult(command_result,commandResult2,50,()=>{
                                    console.log("printed result2")

                                    setTimeout(()=>{
                                        clearCommands()
                                        
                                        command_description.innerHTML = 'Creates a kubernetes cluster.'
                                        command_bg.innerHTML = command3
                                        top_command_bg.innerHTML = command3
                                        printSequence(command_top,command_anim,command3,20,()=>{
                                            console.log("printed third command")
                                            setTimeout(()=>{
                                                printCommandResult(command_result,commandResult3,50,()=>{
                                                    console.log("printed result3")
                                                    
                                                    setTimeout(()=>{
                                                        loopCommands()
                                                    },4000)
                
                                                })
                                            },400)
                                        })
                                    },4000)
                                })
                            },400)
                        })
                    },4000)
                })
            },400)
        })
    
    }

    // printSequence()
    let timeBetween = 500
    let commandLine = document.getElementById('command_line')
    setTimeout(()=>{
        // appendNext(command_top3,command_anim3,commandText,command_result3)
        // animation_done = true
        printCommand(commandLine,'curl -o kubitect -L https://dl.kubitect.com',()=>{
            setTimeout(()=>{
                printCommand(commandLine,'sudo mv kubitect /usr/local/bin/',()=>{
                    setTimeout(()=>{
                        printCommand(commandLine,'kubitect apply',()=>{
                            console.log("all commands printed")
                            setTimeout(()=>{
                                appendMultipleLines(commandLine,commandResult3,50,()=>{
                                    console.log(commandLine.innerHTML)
                                })
                            },timeBetween)
                        })
                    },timeBetween)
                })
            },timeBetween)
        })
    },timeBetween)

    // window.addEventListener('scroll', function() {
    //     var position = command_anim3.getBoundingClientRect();
    
    //     // checking whether fully visible
    //     if(position.top >= 0 && position.bottom <= (3/4)*window.innerHeight) {
    //         if(!animation_done){

                
    //         }
    //     }
    // });


});



animation_done = false

// let commands = [
//     {
//         command:'wget -O kubitect https://download.kubitect.com',
//         result:'$'
//     },
//     { 
//         command:'sudo mv kubitect /usr/local/bin/',
//         result:'$'
//     },
//     {
//         command:'kubitect apply',
//         result: resultText
//     },
// ]

let commandText = 'kubitect apply'
let resultText = "<p class=\"px-2 m-0\">"
// Preparing cluster '<span style=\"background-color:red\">default</span>'...<br>
// Setting up 'main' virtual environment...<br>
// Creating virtual environment...<br>
// Installing pip3 dependencies...<br>
// This can take up to a minute when the virtual environment is initialized for the first time...<br><br>
// PLAY [localhost]<br><br>
// TASK [cluster-config/copy : Make sure config directory exists]<br>
// <span style=\"color:green\">ok: [127.0.0.1]</span><br><br>
// ...


const printSequence = async (top_element,main_element,sequence,speed, callback) => {
    if(sequence.length===0)callback()
    else{
        top_element.innerHTML = top_element.innerHTML += sequence[0]
        main_element.innerHTML = main_element.innerHTML += sequence[0]
        setTimeout(()=>{
            printSequence(top_element,main_element,sequence.substr(1),speed,()=>{
                callback()
            })
        },speed)
    }
}

const printCommandResult = async (result_element,lines,speed,callback) => {
    if(lines.length===0)callback()
    else{
        result_element.innerHTML+=lines[0]+'<br>'
        setTimeout(()=>{
            printCommandResult(result_element,lines.slice(1),speed,()=>{
                callback()
            })
        },speed)
    }

}


const appendLine = async (target,line) => {
    target.innerHTML += line+'<br>'
}

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


const appendNext = async (top_element,element,remainer,result) => {
    if(remainer===''){
        // result.innerHTML = resultText
        setTimeout(()=>{
            result.innerHTML = resultText
        },500)
        
    }else{
        top_element.innerHTML=top_element.innerHTML+=remainer[0]
        element.innerHTML=element.innerHTML+=remainer[0]
        setTimeout(()=>{
            appendNext(top_element,element,remainer.substr(1),result)
        },20)
    }

}