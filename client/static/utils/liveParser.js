// Prevent back navigation
// window.history.pushState(null, '', window.location.href); // Add an entry to history

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
        body: JSON.stringify(args), // { email: email, password: password }
        });
        if (!response.ok) {
        throw data;
        }
        const data = await response.text()
        return data
    } catch (error) {
        console.error(error);
        return data.text();
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
    console.log("request params : ", path, args)
    try {
      const response = await fetch(path+'?'+argsString);
      const data = await response.text()
      return data
    } catch (error) {
      console.error(error);
    }
};

const forms = [...document.getElementsByTagName('form')]
console.log("forms", forms)

forms.forEach(form => form.addEventListener('submit', (event) => {
    console.log(form.querySelectorAll('input'))
    event.preventDefault()

    const formInputs = form.querySelectorAll('input')
    let formValues = {}

    formInputs.forEach((input) => {
      formValues[input.name] = input.value
    }) 

    if (form.method === "post") {
      ExecPostRequest(form.action, formValues).then((result) => {
        if (result == false) {
            console.log("eroooooooor", formValues)
            
        } else {
        console.log(result)
        const nav = document.getElementById('navbar-placeholder').innerHTML
        document.documentElement.innerHTML = result
        console.log("heloooo")
        // Manually re-execute any inline scripts (since the browser doesn't do it automatically)
        const scripts = document.querySelectorAll('script');
        console.log("scripts : ", scripts)
        scripts.forEach((script) => {
          {script.textContent ? eval(script.textContent): null}
          {script.src && script.src != "http://localhost:8080/static/utils/theme.js" && script.src != "../static/utils/theme.js"  ? ExecuteExternalJs(script.src): null} 
        });
        document.getElementById('navbar-placeholder').innerHTML = nav
        }
      })
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
        console.log(data); // Contents of the index.js file
        eval(data)
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}

const posts = [...document.getElementsByClassName('post-card')]

posts.forEach((post) => post.addEventListener('click', () => ExecGetRequest("http://localhost:8080/post", {id : post.getAttribute("card-id")})
  .then((result) => {
    
    history.pushState(null, '', "http://localhost:8080/post?id="+post.getAttribute("card-id"))
    console.log(result)
    const nav = document.getElementById('navbar-placeholder').innerHTML
    document.documentElement.innerHTML = result
    console.log("heloooo")

    // Manually re-execute any inline scripts (since the browser doesn't do it automatically)
    const scripts = document.querySelectorAll('script');
    console.log("scripts : ", scripts)
    scripts.forEach((script) => {
        {script.textContent ? eval(script.textContent): null}
        console.log(script.src)
        {script.src && script.src != "http://localhost:8080/static/utils/theme.js" && script.src != "../static/utils/theme.js"  ? ExecuteExternalJs(script.src): null} 
    });
    document.getElementById('navbar-placeholder').innerHTML = nav
    console.log("heloooo")
  })
))