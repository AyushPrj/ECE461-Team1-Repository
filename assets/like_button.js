function login() {
  // Get the input values from the form
  var username = document.getElementById("username").value;
  var password = document.getElementById("password").value;

  // Hardcoded credentials
  var validUsername = "a";
  var validPassword = "p";

  // Check if the input values match the hardcoded credentials
  if (username === validUsername && password === validPassword) {
    // Redirect to the new page on successful login
    window.location.href = "loggedin.html";
  } else {
    document.getElementById("login-message").innerHTML = "Invalid username or password.";
  }
}

function callAPI() {
  const inputText = document.getElementById("input-text").value;
  fetch("http://localhost:8080/repos", {
    method: "GET",
    body: JSON.stringify({ input: inputText }),
    headers: {
      "Content-Type": "application/json"
    }
  })
  .then(response => response.json())
  .then(data => {
    // Handle API response data
    
  })
  .catch(error => {
    // Handle API error
  });
}