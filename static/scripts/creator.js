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

topTextInput.addEventListener("input", () => {
  updateMemeCanvas(canvas, image, topTextInput.value, bottomTextInput.value, textColorInput.value);
});

bottomTextInput.addEventListener("input", () => {
  updateMemeCanvas(canvas, image, topTextInput.value, bottomTextInput.value, textColorInput.value);
});

textColorInput.addEventListener("input", ()=>{
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

    let words = text.split(' ');
    let line = '';
    let testLine = '';
    let lineArray = []; 

    for(var n = 0; n < words.length; n++) {
        testLine += `${words[n]} `;
        let metrics = ctx.measureText(testLine);
        let testWidth = metrics.width;
        if (testWidth > maxWidth && n > 0) {
            lineArray.push([line, x, y]);
            y += lineHeight;
            line = `${words[n]} `;
            testLine = `${words[n]} `;
        }
        else {
            line += `${words[n]} `;
        }
        if(n === words.length - 1) {
            lineArray.push([line, x, y]);
        }
    }
    return lineArray;
}


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