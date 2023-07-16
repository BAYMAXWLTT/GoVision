package main

import "github.com/gin-gonic/gin"

type Application struct {
	frontAddress, classificationAddress, styleAddress string
	*Controller
	*gin.Engine
}

func NewApplication(frontAddress, classificationAddress, styleAddress string) *Application {
	router := gin.Default()
	router.LoadHTMLGlob(`./templates/*.html`)
	router.Static(`/static`, `./static`)

	router.GET("/image", func(c *gin.Context) {
		c.File("tmp_images/saved_image.jpg")
	})

	router.GET("/brand", func(c *gin.Context) {
		c.File("tmp_images/brand.jpg")
	})

	router.GET("/gradcam", func(c *gin.Context) {
		c.File("tmp_images/cam.jpg")
	})

	controller := NewController(classificationAddress, styleAddress)
	router.GET(`/`, controller.Home)
	router.POST(`/predict`, controller.Predict)
	router.POST(`/style`, controller.StylePost)
	return &Application{
		frontAddress:          frontAddress,
		classificationAddress: classificationAddress,
		styleAddress:          styleAddress,
		Engine:                router}
}

func (app *Application) Run() error {
	return app.Engine.Run(app.frontAddress)
}
