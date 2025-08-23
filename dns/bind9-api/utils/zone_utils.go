package utils

import (
	"bind9-api/config"
	"bind9-api/models"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/miekg/dns"
)

func ParseZoneFile(filename string, origin string) ([]models.ResourceRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open zone file: %w", err)
	}
	defer file.Close()

	var records []models.ResourceRecord
	reader := bufio.NewReader(file)

	zp := dns.NewZoneParser(reader, origin, filename)

	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		hdr := rr.Header()
		record := models.ResourceRecord{
			Name:   hdr.Name,
			TTL:    hdr.Ttl,
			Class:  dns.ClassToString[hdr.Class],
			Type:   dns.TypeToString[hdr.Rrtype],
			Fields: make(map[string]string),
		}

		switch v := rr.(type) {
		case *dns.A:
			record.Fields["address"] = v.A.String()
		case *dns.AAAA:
			record.Fields["address"] = v.AAAA.String()
		case *dns.CNAME:
			record.Fields["target"] = v.Target
		case *dns.NS:
			record.Fields["ns"] = v.Ns
		case *dns.MX:
			record.Fields["preference"] = fmt.Sprint(v.Preference)
			record.Fields["exchange"] = v.Mx
		case *dns.TXT:
			record.Fields["text"] = fmt.Sprint(v.Txt)
		case *dns.SOA:
			record.Fields["ns"] = v.Ns
			record.Fields["mbox"] = v.Mbox
			record.Fields["serial"] = fmt.Sprint(v.Serial)
			record.Fields["refresh"] = fmt.Sprint(v.Refresh)
			record.Fields["retry"] = fmt.Sprint(v.Retry)
			record.Fields["expire"] = fmt.Sprint(v.Expire)
			record.Fields["minttl"] = fmt.Sprint(v.Minttl)
		case *dns.SRV:
			record.Fields["priority"] = fmt.Sprint(v.Priority)
			record.Fields["weight"] = fmt.Sprint(v.Weight)
			record.Fields["port"] = fmt.Sprint(v.Port)
			record.Fields["target"] = v.Target
		case *dns.PTR:
			record.Fields["ptr"] = v.Ptr
		default:
			record.Fields["raw"] = v.String()
		}

		records = append(records, record)
	}

	if err := zp.Err(); err != nil {
		return nil, fmt.Errorf("zone parse error: %w", err)
	}

	return records, nil
}

func WriteZoneFile(zone models.Zone, conf *config.Config) error {
	zoneFileDir := conf.Bind9.ZoneFileDir
	origin := zone.Name
	ttl := zone.TTL
	if ttl == 0 {
		ttl = conf.Bind9.DefaultTTL
	}
	records := zone.Records

	// Ensure the zone directory exists
	if err := os.MkdirAll(conf.Bind9.ZoneFileDir, 0755); err != nil {
		return err
	}

	// Create the zone file
	filePath := filepath.Join(zoneFileDir, origin+".zone")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read old content to revert in case of syntax error
	oldContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Write $ORIGIN and $TTL
	if origin != "" {
		_, _ = fmt.Fprintf(file, "$ORIGIN %s.\n", strings.TrimSuffix(origin, "."))
	}
	if ttl > 0 {
		_, _ = fmt.Fprintf(file, "$TTL %d\n\n", ttl)
	}

	// Write each record
	for _, rr := range records {

		rrStr := MakeRRString(rr)
		record, err := dns.NewRR(rrStr)
		if err != nil {
			if len(oldContent) == 0 {
				os.Remove(filePath)
			} else {
				os.WriteFile(filePath, oldContent, 0644) // Revert to old content
			}
			return fmt.Errorf("zone syntax error: %s", err.Error())
		}

		// Handle @ shorthand for origin
		if rr.Name == origin || rr.Name == "" {
			rr.Name = origin
		}

		_, _ = fmt.Fprintf(file, "%s\n", record.String())
	}

	// Syntax check the zone file
	success, _, err := ExecCmd("named-checkzone", zone.Name, filePath)
	if !success {
		if len(oldContent) == 0 {
			os.Remove(filePath)
		} else {
			os.WriteFile(filePath, oldContent, 0644) // Revert to old content
		}
		return fmt.Errorf("zone syntax error: %s", err.Error())
	}

	return nil
}

func MakeRRString(record models.ResourceRecord) string {
	switch record.Type {
	case "A":
		return fmt.Sprintf("%s %d IN A %s", record.Name, record.TTL, record.Fields["address"])
	case "AAAA":
		return fmt.Sprintf("%s %d IN AAAA %s", record.Name, record.TTL, record.Fields["address"])
	case "CNAME":
		return fmt.Sprintf("%s %d IN CNAME %s", record.Name, record.TTL, record.Fields["target"])
	case "NS":
		return fmt.Sprintf("%s %d IN NS %s", record.Name, record.TTL, record.Fields["ns"])
	case "MX":
		return fmt.Sprintf("%s %d IN MX %s %s", record.Name, record.TTL, record.Fields["preference"], record.Fields["exchange"])
	case "TXT":
		return fmt.Sprintf("%s %d IN TXT \"%s\"", record.Name, record.TTL, record.Fields["text"])
	case "SOA":
		return fmt.Sprintf("%s %d IN SOA %s %s (%s %s %s %s %s)", record.Name, record.TTL, record.Fields["ns"], record.Fields["mbox"], record.Fields["serial"], record.Fields["refresh"], record.Fields["retry"], record.Fields["expire"], record.Fields["minttl"])
	case "SRV":
		return fmt.Sprintf("%s %d IN SRV %s %s %s %s", record.Name, record.TTL, record.Fields["priority"], record.Fields["weight"], record.Fields["port"], record.Fields["target"])
	case "PTR":
		return fmt.Sprintf("%s %d IN PTR %s", record.Name, record.TTL, record.Fields["ptr"])
	default:
		return fmt.Sprintf("%s %d IN %s %s", record.Name, record.TTL, record.Type, record.Fields["raw"])
	}
}
