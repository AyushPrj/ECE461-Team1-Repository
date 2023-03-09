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

function addTable(name, ru, cs, rm, bs, rc, ls, ns) {
  // creates a <table> element and a <tbody> element
  const tbl = document.createElement("table");
  const tblBody = document.createElement("tbody");
  const headers = ['Name', 'Ramp','Correctness', 'Responsive Maintainer', 'Bus Factor', 'Review Coverage', 'License', 'Net']
  const vals = [name, ru, cs, rm, bs, rc, ls, ns]

  // creating all cells
  for (let i = 0; i < 2; i++) {
    // creates a table row
    const row = document.createElement("tr");

    for (let j = 0; j < 8; j++) {
      // Create a <td> element and a text node, make the text
      // node the contents of the <td>, and put the <td> at
      // the end of the table row
      const cell = document.createElement("td");
      cellText = '';
      if (i == 0) {
        cellText = document.createTextNode(headers[j]);
      } else {
        cellText = document.createTextNode(vals[j]);
      }
      cell.appendChild(cellText);
      row.appendChild(cell);
    }

    // add the row to the end of the table body
    tblBody.appendChild(row);
  }

  // put the <tbody> in the <table>
  tbl.appendChild(tblBody);
  // appends <table> into <body>
  document.body.appendChild(tbl);
  // sets the border attribute of tbl to '2'
  tbl.setAttribute("border", "2");
}

// document.getElementById("search_output").style.display = "flex"; 
//       document.getElementById("name").innerHTML = "Package Name";
//       document.getElementById("ru").innerHTML ="Ramp-Up";
//       document.getElementById("cs").innerHTML="Correctness";
//       document.getElementById("rm").innerHTML="Responsive Maintainer";
//       document.getElementById("bf").innerHTML="Bus Factor";
//       document.getElementById("rc").innerHTML="Review Coverage";
//       document.getElementById("ls").innerHTML="License";
//       document.getElementById("ns").innerHTML="Net";
//       document.getElementById("name_search").innerHTML = data[i].URL;
//       document.getElementById("ru_search").innerHTML =data[i].RAMP_UP_SCORE;
//       document.getElementById("cs_search").innerHTML=data[i].CORRECTNESS_SCORE;
//       document.getElementById("rm_search").innerHTML=data[i].RESPONSIVE_MAINTAINER_SCORE;
//       document.getElementById("bf_search").innerHTML=data[i].BUS_FACTOR_SCORE;
//       document.getElementById("rc_search").innerHTML=data[i].REVIEW_COVERAGE_SCORE;
//       document.getElementById("ls_search").innerHTML=data[i].LICENSE_SCORE;
//       document.getElementById("ns_search").innerHTML=data[i].NET_SCORE;

function callAPI() {
  //console.log("hello")
  const inputText = document.getElementById("input-text").value;
  fetch("http://localhost:5500/repos")
    .then(response => response.json())
    .then(data => {
      // Handle API response data
      console.log(inputText);
      for (var i = 0; i < data.length; i++) {
        if (data[i].URL == inputText) {
          addTable(data[i].URL, data[i].RAMP_UP_SCORE, data[i].CORRECTNESS_SCORE, data[i].RESPONSIVE_MAINTAINER_SCORE, data[i].BUS_FACTOR_SCORE, data[i].REVIEW_COVERAGE_SCORE, data[i].LICENSE_SCORE, data[i].NET_SCORE)
          break;
        }
      }
      console.log(data);
      const targetURL = inputText;
      //const targetObject = data.find(obj => obj.URL === targetURL);
      //console.log(targetObject); // outputs the target object
    })
    .catch(error => console.error(error));
}