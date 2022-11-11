"use strict";

const submitBtn = document.getElementById("submitBtn");
const imgInput = document.getElementById("imageInput");
const imgTag = document.getElementById("imgTag");
const imgSizeField = document.getElementById("imgSizeField");
const imgTypeField = document.getElementById("imgTypeField");
const imgLinkField = document.getElementById("imgLinkField");
const imgRateField = document.getElementById("imgRateField");

async function postImage(url, data = {}) {
  const response = await fetch(url, {
    method: "POST",
    mode: "cors",
    headers: {
      "Content-Type": "multipart/form-data",
    },
    body: data,
  });

  return response.json();
}

submitBtn.addEventListener("click", (e) => {
  e.preventDefault();

  if (!imgInput.files.length) {
    alert("¡No has seleccionado una imagen!");
    return;
  }

  const img = imgInput.files[0];

  if (
    img.type != "image/png" &&
    img.type != "image/jpeg" &&
    img.type != "image/jpg"
  ) {
    alert("¡Archivo no válido!");
    return;
  }

  let bodyContent = new FormData();
  bodyContent.append("image", img);

  fetch("http://localhost:3000/image", {
    method: "POST",
    body: bodyContent,
  })
    .then((res) => {
      res.json().then((v) => {
        imgTag.src = v["image"];
        imgTypeField.innerHTML = `<b>Formato:</b> ${
          img.type.split("/")[1]
        } -> ${v["type"]}`;
        imgSizeField.innerHTML = `<b>Tamaño:</b> ${img.size / 1000} kB -> ${
          v["size"]
        }`;
        imgLinkField.innerHTML = `<b>Descargar:</b> <a target="_blank" href="${v["image"]}">${v["image"]}</a>`;

        let compressRate = Math.round(
          ((parseFloat(img.size / 1000) - parseFloat(v["size"])) /
            parseFloat(img.size / 1000)) *
            100
        );

        imgRateField.innerHTML = `<b>Tasa de reducción del tamaño ≈</b> ${compressRate} %`;
      });
    })
    .catch((e) => {
      alert(e);
    });
});
