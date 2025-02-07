// Prevent back navigation

window.onpopstate = function() {
  // This will be triggered when the back button is pressed
  window.location.href = "http://localhost:8080/"
  window.reload()
};

const ExecPostRequest = async (path, args) => {
    console.log("request params : ", path, args)
    try {
        const response = await fetch(path, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(args),
        });
        const data = await response.text()
        console.log(data)
        return data
    } catch (error) {
        console.error(error);
        return false;
    }
};

const ExecGetRequest = async (path, args) => {
    let argsString = ""
    Object.entries(args).forEach(([key, val]) => {
      argsString += key+"="+val+"+"
    })
    args = argsString.split("+")
    if (args[args.length-1] === "") args.pop()
    argsString = args.join('+')
    try {
      const response = await fetch(path+'?'+argsString);
      const data = await response.text()
      return data
    } catch (error) {
      console.error(error);
    }
};

const forms = [...document.getElementsByTagName('form')]

forms.forEach(form => form.addEventListener('submit', async (event) => {
    event.preventDefault()

    const formInputs = [...form.querySelectorAll('input'), ...form.querySelectorAll('textarea')]
    let formValues = {}

    formInputs.forEach((input) => {
      formValues[input.name] = input.value
    }) 

    if (form.method === "post") {
      const result = await Promise.resolve(ExecPostRequest(form.action, formValues))
        if (result === false) {
            window.alert("There was an error")
        } else {
        const nav = document.getElementById('navbar-placeholder').innerHTML
        document.documentElement.innerHTML = result
        // Manually re-execute any inline scripts (since the browser doesn't do it automatically)
        const scripts = document.querySelectorAll('script');
        scripts.forEach((script) => {
          {script.textContent ? eval(script.textContent): null}
          {script.src ? ExecuteExternalJs(script.src): null} 
        });
        document.getElementById('navbar-placeholder').innerHTML = nav
        }
    }
}));

const ExecuteExternalJs = (src) => {
    fetch(src).then(response => {
        if (!response.ok) {
        throw new Error('Network response was not ok');
        }
        return response.text();
    })
    .then(data => {
        eval(data)
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}

const posts = [...document.getElementsByClassName('post-card')]

posts.forEach((post) => post.addEventListener('click', async () => {
  const result = await Promise.resolve(ExecGetRequest("http://localhost:8080/post", {id : post.getAttribute("card-id")}))
  history.pushState(null, '', "http://localhost:8080/post?id="+post.getAttribute("card-id"))
  if (result === false) {
    location.reload()
    return
  }
  const nav = document.getElementById('navbar-placeholder').innerHTML
  document.documentElement.innerHTML = result

  // Manually re-execute any inline scripts (since the browser doesn't do it automatically)
  const scripts = document.querySelectorAll('script');
  scripts.forEach((script) => {
      {script.textContent ? eval(script.textContent): null}
      {script.src ? ExecuteExternalJs(script.src): null} 
  });
  document.getElementById('navbar-placeholder').innerHTML = nav
}
))