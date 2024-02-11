async function buttonSendClick(){
    let elem = document.querySelector('#task-input')
    let response = await fetch('/api/task', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json;charset=utf-8'
        },
        body: elem.value
      });
    let result = await response.text()
    
    alert(result)
}