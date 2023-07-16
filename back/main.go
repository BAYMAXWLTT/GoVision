package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	classificationAddress, styleAddress string
}

func NewController(classificationAddress, styleAddress string) *Controller {
	return &Controller{
		classificationAddress: classificationAddress,
		styleAddress:          styleAddress,
	}
}

func (con *Controller) Home(c *gin.Context) {
	c.HTML(http.StatusOK, `index.html`, nil)
}

func (con *Controller) StylePost(c *gin.Context) {
	type formattedImgs struct {
		Content string `json:"content"`
		Style   string `json:"style"`
	}
	data := formattedImgs{}
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{`err`: err.Error()})
		return
	}
	data_byte, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{`err`: err.Error()})
		return
	}
	_, err = sendPostRequest(fmt.Sprintf("%s/style", con.styleAddress), data_byte, `application/json`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{`err`: err.Error()})
		return
	}
}

func (con *Controller) Predict(c *gin.Context) {
	img_data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{`err`: err.Error()})
		return
	}
	// 访问python的推理服务
	print(string(img_data)[:100])
	resp, err := sendPostRequest(fmt.Sprintf("%s/predict", con.classificationAddress), img_data, `application/octet-stream`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{`err`: err.Error()})
		return
	}
	type Response struct {
		Result      string `json:"result"`
		Probability string `json:"probability"`
	}
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{`err`: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

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
		c.File("saved_image.jpg")
	})

	router.GET("/brand", func(c *gin.Context) {
		c.File("brand.jpg")
	})

	router.GET("/gradcam", func(c *gin.Context) {
		c.File("cam.jpg")
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

func sendPostRequest(url string, body []byte, contentType string) (*http.Response, error) {
	client := http.Client{}
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func main() {
	// frontAddress := `localhost:8000`
	frontAddress := `localhost:8000`
	classificationAddress := `http://localhost:4000`
	styleAddress := `http://localhost:3000`
	app := NewApplication(frontAddress, classificationAddress, styleAddress)
	log.Fatal(app.Run())
}
