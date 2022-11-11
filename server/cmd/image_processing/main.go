package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gera9/image_processing/pkg/helpers"
	"github.com/rs/cors"
)

// Controller.
func processImage(w http.ResponseWriter, r *http.Request) {
	// Get image from the request body.
	multipartFile, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// ReadAll encodes the file object to []byte.
	img, err := io.ReadAll(multipartFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// ImageProcessing compresses the image.
	compressedImg, err := helpers.CompressImage(img, 40)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// GenerateNewImageName generates an unique image name.
	imageName := helpers.GenerateNewImageName(helpers.GetImageExtension(compressedImg))

	// WriteImage writes the compressed image in the given directory.
	helpers.WriteImage(compressedImg, fmt.Sprintf("pkg/uploads/%s", imageName))

	// DetectContentType returns the image extension.
	contentType := http.DetectContentType(img)

	// EncodeToB64 encodes the image to base 64.
	b64, err := helpers.EncodeToB64(compressedImg, contentType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// NewStorage creates a Singleton instance of the MongoDB client.
	storage, err := helpers.NewStorage("mongodb://root:example@mongo/", "image_processing")
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = storage.InsertImage(b64)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	res, err := helpers.BuidlResponse(compressedImg, imageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Respond successfully with a json.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func main() {
	// NewServeMux creates a new router.
	r := http.NewServeMux()

	// Serving Static Files (images).
	fs := http.FileServer(http.Dir("pkg/uploads"))
	r.Handle("/uploads/", http.StripPrefix("/uploads", fs))

	// Get a "Hello, world!".
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":    "Hello, world!",
			"statusCode": http.StatusOK,
		})
	})

	// /image compresses the image and encode it to b64.
	r.HandleFunc("/image", processImage)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	cors := cors.Default().Handler(r)

	// ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections.
	log.Panic(http.ListenAndServe(":3000", cors))
}
