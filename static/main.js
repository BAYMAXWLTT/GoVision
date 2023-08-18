//========================================================================
// Drag and drop image handling
//========================================================================

var fileDrag = document.getElementById("file-drag");
var fileSelect1 = document.getElementById("file-upload1");
var fileSelect2 = document.getElementById("file-upload2");
var fileSelect3 = document.getElementById("file-upload3");

// Add event listeners
fileDrag.addEventListener("dragover", fileDragHover, false);
fileDrag.addEventListener("dragleave", fileDragHover, false);
fileDrag.addEventListener("drop", fileSelectHandler, false);

fileSelect1.addEventListener("change", fileSelectHandler, false);
fileSelect2.addEventListener("change", fileSelectHandler, false);
fileSelect3.addEventListener("change", fileSelectHandler, false);

function fileDragHover(e) {
  // prevent default behaviour
  e.preventDefault();
  e.stopPropagation();

  fileDrag.className = e.type === "dragover" ? "upload-box dragover" : "upload-box";
}

function fileSelectHandler(e) {
  // handle file selecting
  var files = e.target.files || e.dataTransfer.files;
  var id = e.target.id
  if (id == "file-upload1") {
    // reset
    predResult.innerHTML = "";
    imageDisplayGradCam.src = blank
  }
  fileDragHover(e);
  for (var i = 0, f; (f = files[i]); i++) {
    previewFile(f, id);
  }
}

//========================================================================
// Web page elements for functions to use
//========================================================================

var imagePreview1 = document.getElementById("image-preview1");
var imagePreview2 = document.getElementById("image-preview2");
var imagePreview3 = document.getElementById("image-preview3");

var imageDisplayClassify = document.getElementById("image-display-classify");
var imageDisplayStyle = document.getElementById("image-display-style");
var imageDisplayGradCam = document.getElementById("image-display-gradcam")

var uploadCaption1 = document.getElementById("upload-caption1");
var uploadCaption2 = document.getElementById("upload-caption2");
var uploadCaption3 = document.getElementById("upload-caption3");

var predResult = document.getElementById("pred-result");
var loaderClassify = document.getElementById("loaderClassify");
var loaderStyle = document.getElementById("loaderStyle");

//========================================================================
// Main button events
//========================================================================

function submitImageClassify() {
  // action for the submit button
  console.log("submit");

  if (!imageDisplayClassify.src || !imageDisplayClassify.src.startsWith("data")) {
    window.alert("Please select an image before submit.");
    return;
  }

  // loaderClassify.classList.remove("hidden");
  show(loaderClassify)
  imageDisplayClassify.classList.add("loading");

  // call the predict function of the backend
  predictImage(img_classify_src);
}

function clearImageClassify() {
  // reset selected files
  fileSelect1.value = "";

  // remove image sources and hide them
  imagePreview1.src = "";
  imageDisplayClassify.src = blank;
  imageDisplayGradCam.src = blank;
  predResult.innerHTML = "";

  hide(imagePreview1);
  hide(imageDisplayClassify);
  hide(imageDisplayGradCam);
  hide(loaderClassify);
  hide(predResult);
  show(uploadCaption1);

  img_classify_src = ""

  imageDisplayClassify.classList.remove("loading");
}

function previewFile(file, id) {
  // show the preview of the image
  console.log(file.name);
  var fileName = encodeURI(file.name);

  var reader = new FileReader();
  reader.readAsDataURL(file);
  reader.onloadend = () => {

    if (id == "file-upload1") {
      imagePreview1.src = URL.createObjectURL(file);
      show(imagePreview1)
      hide(uploadCaption1);
      displayImage(reader.result, "image-display-classify");
      img_classify_src = reader.result
      imageDisplayClassify.classList.remove("loading");
      // // reset
      // predResult.innerHTML = "";
    } else if (id == "file-upload2") {
      imagePreview2.src = URL.createObjectURL(file);
      show(imagePreview2)
      hide(uploadCaption2);
      img_left_src = reader.result
      imageDisplayStyle.classList.remove("loading");
    } else if (id = "file-upload3") {
      imagePreview3.src = URL.createObjectURL(file);
      show(imagePreview3)
      hide(uploadCaption3);
      img_right_src = reader.result
      imageDisplayStyle.classList.remove("loading");
    }
  };
}

//========================================================================
// Helper functions
//========================================================================

function predictImage(image) {
  fetch("/predict", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(image)
  })
    .then(resp => {
      if (resp.ok)
        resp.json().then(data => {
          displayResult(data);
        });
      imageDisplayGradCam.src = "http://wlt.natapp1.cc/gradcam?" + Date.now();
      show(imageDisplayGradCam)
    })
    .catch(err => {
      console.log("An error occured", err.message);
      window.alert("Oops! Something went wrong.");
    });
}

function displayImage(image, id) {
  // display image on given id <img> element
  let display = document.getElementById(id);
  display.src = image;
  show(display);
}

function displayResult(data) {
  // display the result
  hide(loaderClassify);
  predResult.innerHTML = data.result;
  show(predResult);
}

function hide(el) {
  // hide an element
  el.classList.add("hidden");
}

function show(el) {
  // show an element
  el.classList.remove("hidden");
}
