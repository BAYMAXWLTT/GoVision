var img_classify_src, img_left_src, img_right_src

function submitImageStyle() {
    // action for the submit button
    console.log("submitStyle");
    console.log(imagePreview2.src, ' ', imagePreview3.src)
    if (imagePreview2.src == blank || imagePreview3.src == blank) {
        window.alert("Please select both content image and style before submit.");
        return;
    }
    hide(imageDisplayStyle)
    show(loaderStyle)
    // call the predict function of the backend
    SendImages([img_left_src, img_right_src]);
}

function SendImages(images) {
    const formattedImages = {
        content: images[0],
        style: images[1]
    };
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/style");
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                imageDisplayStyle.src = "http://wlt.natapp1.cc/image?" + Date.now();
                show(imageDisplayStyle)
                hide(loaderStyle)
            } else {
                console.log("An error occured", xhr.statusText);
                window.alert("Oops! Something went wrong.");
            }
        }
    };
    xhr.send(JSON.stringify(formattedImages));
}

function displayStyleImage(data) {
    const img = document.getElementById('image-display-style');
    img.src = 'data:image/png;base64,' + data;
}

blank = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

function clearImageStyle() {
    // reset selected files
    fileSelect2.value = "";
    fileSelect3.value = "";

    // remove image sources and hide them
    imagePreview2.src = blank;
    imagePreview3.src = blank;


    imageDisplayStyle.src = "";

    hide(imagePreview2);
    hide(imagePreview3);
    hide(imageDisplayStyle);
    hide(loaderStyle);
    show(uploadCaption2);
    show(uploadCaption3);

    imageDisplayStyle.classList.remove("loading");
}