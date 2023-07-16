package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	_, err = SendPostRequest(fmt.Sprintf("%s/style", con.styleAddress), data_byte, `application/json`)
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
	resp, err := SendPostRequest(fmt.Sprintf("%s/predict", con.classificationAddress), img_data, `application/octet-stream`)
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
