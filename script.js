//register gsap plugin 
gsap.registerPlugin(MotionPathPlugin);

// Create animations
gsap.set(document.body, {overflow: "hidden"})
gsap.from('.header', {fontSize: 0})
gsap.from('.title',  {fontSize: 0, delay: 1})
gsap.from('.banners', {opacity: 0, stagger: .5, rotation: 720, delay: 2})
gsap.from('.input', {opacity: 0, fontSize: 0, delay: 3})

// Fetch user selected banner and inputted text and run sendText() func
function submitText(){
  
  // Get the selected font from the element
  let banner = ""
  var chk = document.querySelector('input[name="choice"]:checked')
  //Test if a font has been selected else, alert user
  if(chk != null){  
    banner = chk.value
  } else {
    alert("Please select a font first!"); 
  }

  // Get the text the user inputted, if no text is present alert them
  let text = document.getElementById('text-input').value;
  if (text.length == 0) {
    alert("You haven't typed anything!")
  }
    
  // List containing valid banner options
  const banners = ["shadow", "standard", "thinkertoy"];

  /*  Check if selected banner is valid option and that there is text present
      running sendText() func if both are valid*/
  if (banners.includes(banner) && text.length > 0) {
    sendText(text, banner);
  } 
};

// Make request to the server sending the banner and text, and receiving the ascii-art output
async function sendText(text, banner) {
  // Format data into struct
  let tooSend ={Text: text, Banner: banner};

  // Set request method, headers and the body
  let options = {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(tooSend),
  };
  
  // Make the request and await response
  let response = await fetch('http://localhost:4000/ascii-art', options);
  let result = await response.json();

  // Reveal the results and download buttons
  showResults ()
  
  // Update contents of "result" element with the received ascii-art
  document.getElementById("result").innerHTML = result.Value;

  // Add animation to result and turn overflow back on to enable scrolling
  gsap.from('.titleResult', { fontSize: 0})
  gsap.from('.result', { fontSize: 0})
  gsap.set(document.body, {overflow: "auto"})
}

// Show the results section
function showResults () {
  let x = document.getElementById("results");
  let y = document.getElementById("downloadTxt");
  let z = document.getElementById("downloadHtml");
  x.style.display = 'block';
  y.style.display = 'block';
  z.style.display = 'block';
}