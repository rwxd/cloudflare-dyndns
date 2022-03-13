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
		domain := cmd.Flag("domain").Value.String()
		zone := cmd.Flag("zone").Value.String()
		apiToken := cmd.Flag("api-token").Value.String()

		domain = utils.RemoveZoneFromDomainName(domain, zone)

		logrus.Infof("Update DynDNS Entry for %s", domain)
		ipClient, err := utils.NewIPChecker()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ipv4, err := ipClient.GetPublicIPv4Address()
		if err != nil {
			logrus.Warn(err)
			os.Exit(1)
		}
		logrus.Infof("IPv4: %s\n", ipv4)

		ipv6, err := ipClient.GetPublicIPv6Address()
		if err != nil {
			logrus.Warn(err)
		}
		logrus.Infof("IPv6: %s\n", ipv6)

		if ipv4 == "" && ipv6 == "" {
			logrus.Error("No public IP address found")
			os.Exit(1)
		}

		cf := cloudflare.NewCloudFlareClient(apiToken)

		tokenTest, err := cf.VerifyToken()
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		if !tokenTest.Sucess {
			logrus.Error("Invalid API key")
			logrus.Error(tokenTest.Errors)
			os.Exit(1)
		}

		cfZone, err := cf.GetZone(zone)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		logrus.Info("Zone: ", cfZone)
	},
}

func init() {
	updateCmd.Flags().StringP("domain", "d", "", "Domain Name")
	updateCmd.Flags().StringP("zone", "z", "", "Zone Name")
	updateCmd.Flags().StringP("api-token", "t", "", "CloudFlare API Token")
}
