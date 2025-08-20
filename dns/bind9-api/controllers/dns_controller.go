package controllers

import (
	"bufio"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"bind9-api/config"
	"bind9-api/models"
	"bind9-api/utils"

	"github.com/gin-gonic/gin"
)

type DNSController struct {
	Config *config.Config
}

func NewDNSController(cfg *config.Config) *DNSController {
	return &DNSController{Config: cfg}
}

func (dc *DNSController) CreateZone(c *gin.Context) {
	var zone models.Zone
	if err := c.ShouldBindJSON(&zone); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Create zone file
	if err := utils.CreateZoneFile(zone, dc.Config.Bind9.ZoneFileDir); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Update Bind9 config
	zoneFilePath := filepath.Join(dc.Config.Bind9.ZoneFileDir, zone.Name+".zone")
	if err := utils.UpdateZoneConfig(dc.Config.Bind9.ConfigFile, zone.Name, zoneFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Reload Bind9 to apply changes
	stdOut, err := utils.ReloadBind9()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to reload Bind9: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Zone created successfully",
		Details: stdOut,
	})
}

func (dc *DNSController) GetZone(c *gin.Context) {
	zoneName := c.Param("zone")
	zoneFilePath := filepath.Join(dc.Config.Bind9.ZoneFileDir, zoneName+".zone")

	file, err := os.Open(zoneFilePath)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Zone not found",
		})
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var records []models.Record

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "$") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		record := models.Record{
			Name:  parts[0],
			TTL:   3600, // Default TTL
			Type:  parts[2],
			Value: parts[3],
		}

		// Handle @ for zone name
		if record.Name == "@" {
			record.Name = zoneName
		}

		records = append(records, record)
	}

	c.JSON(http.StatusOK, models.Zone{
		Name:    zoneName,
		Records: records,
	})
}

func (dc *DNSController) UpdateZone(c *gin.Context) {
	zoneName := c.Param("zone")
	var update models.ZoneUpdate

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Get existing zone to preserve SOA and NS records
	zoneFilePath := filepath.Join(dc.Config.Bind9.ZoneFileDir, zoneName+".zone")
	existingZone, err := dc.getExistingZone(zoneName, zoneFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Preserve SOA and NS records from the existing zone
	var preservedRecords []models.Record
	for _, r := range existingZone.Records {
		if r.Type == "SOA" || (r.Name == "@" && r.Type == "NS") {
			preservedRecords = append(preservedRecords, r)
		}
	}

	// Create a new zone with the preserved records and updated records
	updatedZone := models.Zone{
		Name:    zoneName,
		Records: append(preservedRecords, update.Records...),
	}

	// Recreate the zone file
	if err := utils.CreateZoneFile(updatedZone, dc.Config.Bind9.ZoneFileDir); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Reload Bind9 to apply changes
	stdOut, err := utils.ReloadBind9()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to reload Bind9: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Zone updated successfully",
		Details: stdOut,
	})
}

func (dc *DNSController) getExistingZone(zoneName, zoneFilePath string) (models.Zone, error) {
	file, err := os.Open(zoneFilePath)
	if err != nil {
		return models.Zone{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var records []models.Record

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "$") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		record := models.Record{
			Name:  parts[0],
			TTL:   3600, // Default TTL
			Type:  parts[2],
			Value: parts[3],
		}

		if record.Name == "@" {
			record.Name = zoneName
		}

		records = append(records, record)
	}

	return models.Zone{
		Name:    zoneName,
		Records: records,
	}, nil
}

func (dc *DNSController) DeleteZone(c *gin.Context) {
	zoneName := c.Param("zone")

	// Delete zone file
	if err := utils.DeleteZoneFile(dc.Config.Bind9.ZoneFileDir, zoneName); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Remove zone from config
	if err := utils.RemoveZoneFromConfig(dc.Config.Bind9.ConfigFile, zoneName); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Reload Bind9 to apply changes
	stdOut, err := utils.ReloadBind9()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to reload Bind9: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Zone deleted successfully",
		Details: stdOut,
	})
}

func (dc *DNSController) AddRecord(c *gin.Context) {
	zoneName := c.Param("zone")
	var record models.Record

	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Get existing zone
	zoneFilePath := filepath.Join(dc.Config.Bind9.ZoneFileDir, zoneName+".zone")
	zone, err := dc.getExistingZone(zoneName, zoneFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Add new record
	zone.Records = append(zone.Records, record)

	// Recreate the zone file
	if err := utils.CreateZoneFile(zone, dc.Config.Bind9.ZoneFileDir); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Reload Bind9 to apply changes
	stdOut, err := utils.ReloadBind9()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to reload Bind9: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Record added successfully",
		Details: stdOut,
	})
}

func (dc *DNSController) DeleteRecord(c *gin.Context) {
	zoneName := c.Param("zone")
	recordName := c.Param("record")
	recordType := c.Param("type")

	// Get existing zone
	zoneFilePath := filepath.Join(dc.Config.Bind9.ZoneFileDir, zoneName+".zone")
	zone, err := dc.getExistingZone(zoneName, zoneFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Filter out the record to delete
	var updatedRecords []models.Record
	for _, r := range zone.Records {
		if r.Name == recordName && r.Type == recordType {
			continue
		}
		updatedRecords = append(updatedRecords, r)
	}

	if len(updatedRecords) == len(zone.Records) {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Record not found",
		})
		return
	}

	zone.Records = updatedRecords

	// Recreate the zone file
	if err := utils.CreateZoneFile(zone, dc.Config.Bind9.ZoneFileDir); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Reload Bind9 to apply changes
	stdOut, err := utils.ReloadBind9()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to reload Bind9: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Record deleted successfully",
		Details: stdOut,
	})
}

// Add this method to the DNSController in controllers/dns_controller.go
func (dc *DNSController) ListZones(c *gin.Context) {
	zones, err := utils.ListZones(dc.Config.Bind9.ZoneFileDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"zones":   zones,
		"count":   len(zones),
	})
}
