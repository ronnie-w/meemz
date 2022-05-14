const imageFileInput = document.querySelector("#imageFileInput");
const canvas = document.querySelector("#meme");
const topTextInput = document.querySelector("#topTextInput");
const bottomTextInput = document.querySelector("#bottomTextInput");
const textColorInput = document.querySelector("#textColorInput");
var canvasTxt = window.canvasTxt;

let image;

imageFileInput.addEventListener("change", (e) => {
  const imageDataUrl = URL.createObjectURL(e.target.files[0]);

  image = new Image();
  image.src = imageDataUrl;

  image.addEventListener(
    "load",
    () => {
      updateMemeCanvas(
        canvas,
        image,
        topTextInput.value,
        bottomTextInput.value,
        textColorInput.value
      );
    },
    { once: true }
  );
});

topTextInput.addEventListener("change", () => {
  updateMemeCanvas(canvas, image, topTextInput.value, bottomTextInput.value, textColorInput.value);
});

bottomTextInput.addEventListener("change", () => {
  updateMemeCanvas(canvas, image, topTextInput.value, bottomTextInput.value, textColorInput.value);
});

textColorInput.addEventListener("change", ()=>{
  updateMemeCanvas(canvas, image, topTextInput.value, bottomTextInput.value, textColorInput.value);
});

function updateMemeCanvas(canvas, image, topText, bottomText, textColor) {
  const ctx = canvas.getContext("2d");
  const width = image.width;
  const height = image.height;
  const fontSize = Math.floor(width / 20);
  const yOffset = height / 25;

  // Update canvas background
  canvas.width = width;
  canvas.height = height;
  ctx.drawImage(image, 0, 0);

  // Prepare text
  ctx.strokeStyle = "black";
  ctx.fillStyle = textColor;
  ctx.lineJoin = "round";
  ctx.font = `${fontSize}px sans-serif`;


// Custom function for wrapping text
const wrapText = function(ctx, text, x, y, maxWidth, lineHeight) {
   //splitting all of the text into words, but splitting it into an array split by spaces
    let words = text.split(' ');
    let line = '';  // This will store the text of the current line
    let testLine = '';  // This will store the text when we add a word, to test if it's too long
    let lineArray = [];  // This is an array of lines, which the function will return

   // Iterating over each word
    for(var n = 0; n < words.length; n++) {
      // Create a test line, and measure it
        testLine += `${words[n]} `;
        let metrics = ctx.measureText(testLine);
        let testWidth = metrics.width;
      // If the width of this test line is more than the max width
        if (testWidth > maxWidth && n > 0) {
           // Then the line is finished, push the current line into "lineArray"
            lineArray.push([line, x, y]);
           // Increase the line height, so a new line is started
            y += lineHeight;
          // Update line and test line to use this word as the first word on the next line
            line = `${words[n]} `;
            testLine = `${words[n]} `;
        }
        else {
           // If the test line is still less than the max width, then add the word to the current line
            line += `${words[n]} `;
        }
      // If we never reach the full max width, then there is only one line.. so push it into the lineArray so we return something
        if(n === words.length - 1) {
            lineArray.push([line, x, y]);
        }
    }
   // Return the line array
    return lineArray;
}

// item[0] is the text
// item[1] is the x coordinate to fill the text at
// item[2] is the y coordinate to fill the text at

// Add top text
ctx.textBaseline = "top";
let wrappedToptext = wrapText(ctx, topText, 20, yOffset + 20, width,80);
wrappedToptext.forEach(function(item) {
  ctx.fillText(item[0], item[1], item[2]); 
});

// Add bottom text
ctx.textBaseline = "bottom";
let wrappedBottomtext = wrapText(ctx, bottomText, 20, height - (yOffset + 60), width,80);
wrappedBottomtext.forEach(function(item) {
  ctx.fillText(item[0], item[1], item[2]); 
});
}