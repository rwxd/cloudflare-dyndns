package cmd

import (
	"fmt"
	"os"

	"github.com/rwxd/cloudflare-dyndns/cloudflare"
	"github.com/rwxd/cloudflare-dyndns/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update DynDNS Entry",
	Run: func(cmd *cobra.Command, args []string) {
		record := cmd.Flag("record").Value.String()
		zone := cmd.Flag("zone").Value.String()
		apiToken := cmd.Flag("api-token").Value.String()
		recordTTL, _ := cmd.Flags().GetInt("ttl")
		logLevel, err := utils.GetLogrusLogLevelFromString(cmd.Flag("log-level").Value.String())
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}
		logrus.SetLevel(logLevel)

		domain := utils.CombineRecordAndZone(record, zone)
		logrus.Infof("Update DynDNS Entry for %s", domain)

		cf := cloudflare.NewCloudFlareClient(apiToken)

		tokenTest, err := cf.VerifyToken()
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		if !tokenTest.Success {
			logrus.Error("Invalid API key")
			logrus.Infof("CloudFlare token verification: %+v", tokenTest)
			logrus.Error(tokenTest.Errors)
			os.Exit(1)
		}

		logrus.Debug("API key verification was successful")

		cfZone, err := cf.GetZone(zone)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		ipClient, err := utils.NewIPChecker()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, ipVersion := range []int{4, 6} {
			ip, err := ipClient.GetIPAddress(ipVersion)
			if err != nil {
				logrus.Error(err)
				continue
			}

			if ip == "" {
				logrus.Errorf("No IPv%v found", ipVersion)
				continue
			}

			dnsType, err := utils.GetDNSTypeForIPVersion(ipVersion)
			if err != nil {
				logrus.Error(err)
				continue
			}

			record := &cloudflare.CloudFlareDNSRecordBody{
				Name:    domain,
				Type:    dnsType,
				Content: ip,
				Ttl:     recordTTL,
			}

			recordExists, err := cf.CheckDNSRecordAlreadyExists(cfZone.ID, dnsType, domain)
			if err != nil {
				continue
			}

			if recordExists {
				logrus.Printf("%s Record already exists\n", dnsType)
				cfRecord, err := cf.GetDNSRecord(cfZone.ID, dnsType, domain)
				if err != nil {
					logrus.Error(err)
					os.Exit(1)
				}

				if cfRecord.Content != record.Content || cfRecord.TTL != record.Ttl {
					fmt.Printf("Updating %s Record with content \"%s\" & ttl %v\n", dnsType, ip, recordTTL)
					err = cf.UpdateDNSRecord(cfZone.ID, cfRecord.ID, record)
					if err != nil {
						logrus.Error(err)
						os.Exit(1)
					}
				} else {
					fmt.Printf("No update needed for %s Record\n", dnsType)
				}

			} else {
				fmt.Printf("Creating %s Record with content \"%s\" & ttl %v\n", dnsType, ip, recordTTL)
				err = cf.CreateDNSRecord(cfZone.ID, record)
				if err != nil {
					logrus.Error(err)
					os.Exit(1)
				}

			}

		}

	},
}

func init() {
	updateCmd.Flags().StringP("record", "r", "", "dns record to change")
	updateCmd.Flags().StringP("zone", "z", "", "zone name")
	updateCmd.Flags().StringP("api-token", "t", "", "cloudflare api token")
	updateCmd.Flags().Int("ttl", 1, "ttl for record")
	updateCmd.Flags().String("log-level", "warning", "log level")
}
