package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/config"
	"github.com/songquanpeng/one-api/common/helper"
	"github.com/songquanpeng/one-api/model"
	"net/http"
	"strconv"
	"strings"
)

func GetAllChannels(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	channels, err := model.GetAllChannels(p*config.ItemsPerPage, config.ItemsPerPage, "limited")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    channels,
	})
	return
}

func SearchChannels(c *gin.Context) {
	keyword := c.Query("keyword")
	channels, err := model.SearchChannels(keyword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    channels,
	})
	return
}

func GetChannel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	channel, err := model.GetChannelById(id, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    channel,
	})
	return
}

func AddChannel(c *gin.Context) {
	channel := model.Channel{}
	err := c.ShouldBindJSON(&channel)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	channel.CreatedTime = helper.GetTimestamp()
	keys := strings.Split(channel.Key, "\n")
	channels := make([]model.Channel, 0, len(keys))
	for _, key := range keys {
		if key == "" {
			continue
		}
		localChannel := channel
		localChannel.Key = key
		channels = append(channels, localChannel)
	}
	err = model.BatchInsertChannels(channels)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

func DeleteChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	channel := model.Channel{Id: id}
	err := channel.Delete()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

func DeleteDisabledChannel(c *gin.Context) {
	rows, err := model.DeleteDisabledChannel()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    rows,
	})
	return
}

func UpdateChannel(c *gin.Context) {
	channel := model.Channel{}
	err := c.ShouldBindJSON(&channel)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	err = channel.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    channel,
	})
	return
}

type ExportChannel struct {
	Id           int     `json:"id"`
	Type         int     `json:"type"`
	Key          *string `json:"key,omitempty"`
	Name         string  `json:"name"`
	Status       int     `json:"status"`
	Weight       *uint   `json:"weight"`
	BaseURL      *string `json:"base_url"`
	Models       string  `json:"models"`
	Group        string  `json:"group"`
	ModelMapping *string `json:"model_mapping"`
	Priority     *int64  `json:"priority"`
	Config       string  `json:"config"`
	SystemPrompt *string `json:"system_prompt"`
}

func ExportChannels(c *gin.Context) {
	scope := c.Query("scope")
	includeKey := c.Query("include_key") == "true"

	channels, err := model.GetAllChannels(0, 0, scope)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	exportChannels := make([]ExportChannel, 0, len(channels))
	for _, ch := range channels {
		exportCh := ExportChannel{
			Id:           ch.Id,
			Type:         ch.Type,
			Name:         ch.Name,
			Status:       ch.Status,
			Weight:       ch.Weight,
			BaseURL:      ch.BaseURL,
			Models:       ch.Models,
			Group:        ch.Group,
			ModelMapping: ch.ModelMapping,
			Priority:     ch.Priority,
			Config:       ch.Config,
			SystemPrompt: ch.SystemPrompt,
		}
		if includeKey {
			exportCh.Key = &ch.Key
		}
		exportChannels = append(exportChannels, exportCh)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    exportChannels,
	})
	return
}

type ImportChannel struct {
	Id           int     `json:"id"`
	Type         int     `json:"type"`
	Key          string  `json:"key"`
	Name         string  `json:"name"`
	Status       int     `json:"status"`
	Weight       *uint   `json:"weight"`
	BaseURL      *string `json:"base_url"`
	Models       string  `json:"models"`
	Group        string  `json:"group"`
	ModelMapping *string `json:"model_mapping"`
	Priority     *int64  `json:"priority"`
	Config       string  `json:"config"`
	SystemPrompt *string `json:"system_prompt"`
}

func ImportChannels(c *gin.Context) {
	var channels []ImportChannel
	err := c.ShouldBindJSON(&channels)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	successCount := 0
	updateCount := 0
	for _, ch := range channels {
		if ch.Id > 0 {
			// Update existing channel
			existingChannel, err := model.GetChannelById(ch.Id, true)
			if err == nil && existingChannel != nil {
				// Only update non-empty fields
				if ch.Name != "" {
					existingChannel.Name = ch.Name
				}
				if ch.Status > 0 {
					existingChannel.Status = ch.Status
				}
				if ch.Weight != nil {
					existingChannel.Weight = ch.Weight
				}
				if ch.BaseURL != nil {
					existingChannel.BaseURL = ch.BaseURL
				}
				if ch.Models != "" {
					existingChannel.Models = ch.Models
				}
				if ch.Group != "" {
					existingChannel.Group = ch.Group
				}
				if ch.ModelMapping != nil {
					existingChannel.ModelMapping = ch.ModelMapping
				}
				if ch.Priority != nil {
					existingChannel.Priority = ch.Priority
				}
				if ch.Config != "" {
					existingChannel.Config = ch.Config
				}
				if ch.SystemPrompt != nil {
					existingChannel.SystemPrompt = ch.SystemPrompt
				}
				// Only update key if explicitly provided and not empty
				if ch.Key != "" {
					existingChannel.Key = ch.Key
				}
				err = existingChannel.Update()
				if err == nil {
					updateCount++
				}
			}
		} else if ch.Key != "" && ch.Type > 0 {
			// Create new channel
			newChannel := model.Channel{
				Type:         ch.Type,
				Key:          ch.Key,
				Name:         ch.Name,
				Status:       ch.Status,
				Weight:       ch.Weight,
				BaseURL:      ch.BaseURL,
				Models:       ch.Models,
				Group:        ch.Group,
				ModelMapping: ch.ModelMapping,
				Priority:     ch.Priority,
				Config:       ch.Config,
				SystemPrompt: ch.SystemPrompt,
			}
			if newChannel.Status == 0 {
				newChannel.Status = 1 // Default to enabled
			}
			newChannel.CreatedTime = helper.GetTimestamp()
			err = model.BatchInsertChannels([]model.Channel{newChannel})
			if err == nil {
				successCount++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"created": successCount,
			"updated": updateCount,
		},
	})
	return
}
