const imageFileInput = document.querySelector("#imageFileInput");
const canvas = document.querySelector("#meme");
const topTextInput = document.querySelector("#topTextInput");
const bottomTextInput = document.querySelector("#bottomTextInput");
const textColorInput = document.querySelector("#textColorInput");

let image;

imageFileInput.addEventListener("change", (e) => {
  const imageDataUrl = URL.createObjectURL(e.target.files[0]);

  image = new Image();
  image.src = imageDataUrl;
  image.style.width = "400px";
  image.style.maxHeight = "500px";
  canvas.style.display = "block";
  $(canvas).attr("width", "400px");
  $(canvas).attr("height", image.style.getPropertyValue("height"));

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
  console.log(image);
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
  const fontSize = Math.floor(width / 15);
  const yOffset = height / 25;

  // Update canvas background
  canvas.width = width;
  canvas.height = height;
  ctx.drawImage(image, 0, 0);

  // Prepare text
  ctx.strokeStyle = "black";
  ctx.lineWidth = Math.floor(fontSize / 10);
  ctx.fillStyle = textColor;
  ctx.textAlign = "center";
  ctx.lineJoin = "round";
  ctx.font = `${fontSize}px sans-serif`;

  // Add top text
  ctx.textBaseline = "top";
  ctx.strokeText(topText, width / 2, yOffset);
  ctx.fillText(topText, width / 2, yOffset);

  // Add bottom text
  ctx.textBaseline = "bottom";
  ctx.strokeText(bottomText, width / 2, height - yOffset);
  ctx.fillText(bottomText, width / 2, height - yOffset);
}
