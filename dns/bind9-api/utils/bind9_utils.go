package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"bind9-api/config"
	"bind9-api/models"
)

const zoneFileTemplate = `$ORIGIN {{.Name}}.
$TTL 3600
@    IN    SOA    ns1.{{.Name}}. admin.{{.Name}}. (
            2023081801 ; Serial
            3600       ; Refresh
            1800       ; Retry
            604800     ; Expire
            86400      ; Minimum TTL
)

; Name servers
@    IN    NS    ns1.{{.Name}}.
@    IN    NS    ns2.{{.Name}}.

; A records
ns1    IN    A    {{getRecordValue .Records "A" "ns1"}}
ns2    IN    A    {{getRecordValue .Records "A" "ns2"}}

; Additional records
{{range .Records}}{{if and (ne .Name "@") (ne .Type "SOA") (ne .Name "ns1") (ne .Name "ns2")}}
{{.Name}}    IN    {{.Type}}    {{.Value}}{{end}}
{{end}}`

func getRecordValue(records []models.Record, rtype string, name string) string {
	for _, r := range records {
		if r.Name == name && r.Type == rtype {
			return r.Value
		}
	}
	return "127.0.0.1" // default value
}

func CreateZoneFile(zone models.Zone, zoneFileDir string) error {
	// Ensure the zone directory exists
	if err := os.MkdirAll(zoneFileDir, 0755); err != nil {
		return err
	}

	// Create the zone file
	filePath := filepath.Join(zoneFileDir, zone.Name+".zone")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add default NS records if not provided
	hasNS1, hasNS2 := false, false
	for _, r := range zone.Records {
		if r.Name == "ns1" && r.Type == "A" {
			hasNS1 = true
		}
		if r.Name == "ns2" && r.Type == "A" {
			hasNS2 = true
		}
	}

	if !hasNS1 {
		zone.Records = append(zone.Records, models.Record{
			Name:  "ns1",
			Type:  "A",
			Value: "127.0.0.1",
			TTL:   3600,
		})
	}

	if !hasNS2 {
		zone.Records = append(zone.Records, models.Record{
			Name:  "ns2",
			Type:  "A",
			Value: "127.0.0.1",
			TTL:   3600,
		})
	}

	// Parse and execute the template
	tmpl := template.New("zonefile").Funcs(template.FuncMap{
		"getRecordValue": getRecordValue,
	})

	tmpl, err = tmpl.Parse(zoneFileTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, zone)
}

func UpdateZoneConfig(configFile string, zoneName string, zoneFile string) error {
	// Check if the zone already exists in the config
	file, err := os.OpenFile(configFile, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	zoneExists := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), `zone "`+zoneName+`"`) {
			zoneExists = true
			break
		}
	}

	if zoneExists {
		return nil
	}

	// Append the new zone configuration
	zoneConfig := fmt.Sprintf(`
zone "%s" {
    type master;
    file "%s";
};`, zoneName, zoneFile)

	if _, err := file.WriteString(zoneConfig); err != nil {
		return err
	}

	return nil
}

func DeleteZoneFile(zoneFileDir string, zoneName string) error {
	filePath := filepath.Join(zoneFileDir, zoneName+".zone")
	return os.Remove(filePath)
}

func RemoveZoneFromConfig(configFile string, zoneName string) error {
	// This is a simplified implementation
	// In a real application, you'd want to properly parse the config file
	// and remove only the relevant zone block

	input, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	var newLines []string
	inZoneBlock := false

	for _, line := range lines {
		if strings.Contains(line, `zone "`+zoneName+`"`) {
			inZoneBlock = true
			continue
		}

		if inZoneBlock && strings.Contains(line, "};") {
			inZoneBlock = false
			continue
		}

		if !inZoneBlock {
			newLines = append(newLines, line)
		}
	}

	output := strings.Join(newLines, "\n")
	return os.WriteFile(configFile, []byte(output), 0644)
}

// Add this function to utils/bind9_utils.go
func ListZones(zoneFileDir string) ([]string, error) {
	var zones []string

	// Check if the directory exists
	if _, err := os.Stat(zoneFileDir); os.IsNotExist(err) {
		return zones, nil // Return empty slice if directory doesn't exist
	}

	files, err := os.ReadDir(zoneFileDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".zone") {
			zoneName := strings.TrimSuffix(file.Name(), ".zone")
			zones = append(zones, zoneName)
		}
	}

	return zones, nil
}

func ReloadBind9() (bool, string, error) {
	return ExecCmd("rndc", "reload")
}

func CheckAndReload(zoneName string, config *config.Config) (bool, string, error) {
	// Check zone syntax
	zoneFilePath := filepath.Join(config.Bind9.ZoneFileDir, zoneName+".zone")
	success, output, err := ExecCmd("named-checkzone", zoneName, zoneFilePath)
	if !success {
		return false, output, fmt.Errorf("zone syntax error: %s", err.Error())
	}

	// Check config syntax
	configFile := config.Bind9.ConfigFile
	success, output, err = ExecCmd("named-checkconf", configFile)
	if !success {
		return false, output, fmt.Errorf("config syntax error: %s", err.Error())
	}
	// Reload Bind9
	success, output, err = ExecCmd("rndc", "reload")
	if !success {
		return false, output, fmt.Errorf("failed to reload Bind9: %s", err.Error())
	}

	return success, output, err
}
