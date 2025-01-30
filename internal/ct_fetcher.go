package internal

import (
	"context"
	"strings"

	"github.com/Diniboy1123/transparnsee/config"
	ctg "github.com/google/certificate-transparency-go"
	ct "github.com/google/certificate-transparency-go/client"
)

func ProcessEntries(client *ct.LogClient, startIndex, endIndex int64, resultCh chan<- string) int {
	entries, err := client.GetEntries(context.Background(), startIndex, endIndex)
	if err != nil {
		return ProcessEntries(client, startIndex, endIndex, resultCh)
	}

	var processed int
	for _, entry := range entries {
		processed++
		domain, organizations := extractDomainAndIssuerOrganizations(entry)
		if domain == "" || len(organizations) == 0 {
			continue
		}

		if strings.HasSuffix(domain, config.AppConfig.CommonNameSuffix) && isRootDomain(domain) && containsIssuerOrganization(organizations) {
			if CheckNSRecordsCloudflareDoH(domain) {
				resultCh <- domain
			}
		}
	}
	return processed
}

func extractDomainAndIssuerOrganizations(entry ctg.LogEntry) (string, []string) {
	if entry.Leaf.TimestampedEntry.EntryType != ctg.X509LogEntryType {
		return "", nil
	}
	return entry.X509Cert.Subject.CommonName, entry.X509Cert.Issuer.Organization
}

func isRootDomain(domain string) bool {
	// ugly, but it works for my use case
	return strings.Count(domain, ".") == 1
}

func containsIssuerOrganization(organizations []string) bool {
	for _, org := range organizations {
		if listContains(org, config.AppConfig.TrustedIssuers) {
			return true
		}
	}
	return false
}

func listContains(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
