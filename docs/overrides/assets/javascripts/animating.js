

window.addEventListener('load', (event) => {
    console.log('page is fully loaded');
    let command_anim3 = document.getElementById('command-anim')
    let command_result3 = document.getElementById('command-result-fg')
    // console.log(command_anim)
    // console.log(command_result)
    // setTimeout(()=>{
    //     command_anim.innerHTML = commandText
    // },1000)

    setTimeout(()=>{
        appendNext(command_anim3,commandText,command_result3)
        animation_done = true
    },1000)

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

let commandText = 'kubitect apply'
let resultText = "<p>Preparing cluster '<span style=\"background-color:red\">default</span>'...<br>Setting up 'main' virtual environment...<br>Creating virtual environment...<br>Installing pip3 dependencies...<br>This can take up to a minute when the virtual environment is initialized for the first time...<br><br>PLAY [localhost]<br><br>TASK [cluster-config/copy : Make sure config directory exists]<br><span style=\"color:green\">ok: [127.0.0.1]</span><br><br>...</p>"
let currentText = ''
let currentIndex = 0
let length = commandText.length

const appendNext = async (element,remainer,result) => {
    if(remainer===''){
        // result.innerHTML = resultText
        setTimeout(()=>{
            result.innerHTML = resultText
        },500)
        
    }else{

        element.innerHTML=element.innerHTML+=remainer[0]
        setTimeout(()=>{
            appendNext(element,remainer.substr(1),result)
        },20)
    }

}