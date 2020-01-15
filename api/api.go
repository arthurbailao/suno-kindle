package api

import (
	"net/http"
	"strings"

	"github.com/arthurbailao/suno-kindle/database"
	"github.com/arthurbailao/suno-kindle/domain"
	"github.com/arthurbailao/suno-kindle/worker"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	DB     *database.DB
	Worker *worker.Worker
}

func (ctrl Controller) CreateOrUpdateDevice(c *gin.Context) {
	var d domain.Device
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": strings.Split(err.Error(), "\n")})
		return
	}

	err := ctrl.DB.CreateOrUpdateDevice(d)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": strings.Split(err.Error(), "\n")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"device": d})
}

func (ctrl Controller) ListDevices(c *gin.Context) {
	devices, err := ctrl.DB.ListDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": strings.Split(err.Error(), "\n")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

func (ctrl Controller) Process(c *gin.Context) {
	ch := make(chan worker.Response)
	ctrl.Worker.Call(ch)

	resp := <-ch
	if resp.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": strings.Split(resp.Err.Error(), "\n")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": resp.Reports})

}
