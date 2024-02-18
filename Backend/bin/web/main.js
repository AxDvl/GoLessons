async function buttonSendClick(){
    let elem = document.querySelector('#task-input')
    await fetch('/api/task', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json;charset=utf-8'
        },
        body: elem.value
      });
    buttonResultsClick()
}

async function buttonResultsClick(){
  let output = document.querySelector('#results');
  if (!output) return
  let response = await fetch('/api/task');
  let result = await response.json();
  while (output.firstChild) {
    output.removeChild(output.firstChild);
  }

  for (let task of result){
    let resultText = task.TaskText;
    if (task.Status == 0) resultText = "[pending] " + resultText;
    else if (task.Status == 1) resultText = "[done] " + resultText + "=" + task.Result;
    else if (task.Status == 2) resultText = "[error] " + resultText;
    li = document.createElement('li')
    li.innerHTML = resultText
    output.append(li)
  }

    
  
  //alert(result)
}