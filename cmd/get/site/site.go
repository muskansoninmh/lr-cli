package site

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var all *bool
var active *bool
var appid *int64

func NewSiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Shows Current/All sites",
		Long: heredoc.Doc(`
		Use this command to get the information about the:
			- Current site/app (--active)
			- All sites/app (--all) 
			- Specific site based on the appid (--appid)
		`),
		Example: heredoc.Doc(`
			$ lr get site --all
			All sites: 
			+--------+-----------------+-------------------------+
			|   ID   |      NAME       |         DOMAIN          |
			+--------+-----------------+-------------------------+
			| 111111 | new-test1       | https://mail7.io        |
			| 122222 | my-app-final    | loginradius.com         | 
			| 142670 | trail-pro       | https://loginradius.com | 

			$ lr get site --active
			Current site: 
			....

			$ lr get site --appid <appid>
			....

		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getSite()
		},
	}
	fl := cmd.Flags()
	all = fl.Bool("all", false, "Lists all sites")
	active = fl.Bool("active", false, "Shows active site")
	appid = fl.Int64P("appid", "i", -1, "Filters sites based on ID")
	return cmd
}

func getSite() error {
	AppInfo,SharedAppInfo, err := api.GetAppsInfo()
	if err != nil {
		return err
	}

	if *active && (!*all && *appid == -1) {
		currentID, err := api.CurrentID()
		if err != nil {
			return err
		}
		fmt.Println("Active site: ")
		if len(SharedAppInfo) != 0 {
			val, _ := SharedAppInfo[currentID] 
			Output(val.Appname, val.Appid,val.Domain)
		} else {	
			vals, _ := AppInfo[currentID] 
			Output(vals.Appname, vals.Appid,vals.Domain)
		}
	} else if *all && (!*active && *appid == -1) {
		var data [][]string
		var sharedAppdata [][]string
		fmt.Println("All sites: ")
		if len(AppInfo) != 0 {
			fmt.Println("Your sites: ")
			for _, site := range AppInfo {
				data = append(data, []string{strconv.FormatInt(site.Appid, 10), site.Appname, site.Domain})
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Domain"})
			table.AppendBulk(data)
			table.Render()
		} 
		if len(SharedAppInfo) != 0 {
			fmt.Println("Shared sites: ")
			for _, site := range SharedAppInfo {
				sharedAppdata = append(sharedAppdata, []string{strconv.FormatInt(site.Appid, 10), site.Appname, site.Domain})
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Domain"})
			table.AppendBulk(sharedAppdata)
			table.Render()
		}
	} else if *appid != -1 && (!*active && !*all) {
		site, ok := AppInfo[*appid]
		sharedsite, sharedAppstatus := SharedAppInfo[*appid]

		if !ok && !sharedAppstatus {
			return errors.New("There is no site with this App ID")
		}
		if ok {
			Output(site.Appname, site.Appid,site.Domain)
		} else {
			Output(sharedsite.Appname, sharedsite.Appid,sharedsite.Domain)
			
		}

	} else {
		fmt.Println("Use exactly one of the following flags: ")
		fmt.Println("--all: Displays all sites ")
		fmt.Println("--active: Displays active site: ")
		fmt.Println("--appid: Displays site with entered appid")

	}

	return nil
}

func Output(AppName string, Appid int64, Domain string) {
	fmt.Println("------------------------------")
	fmt.Println("App Name: ", AppName)
	fmt.Println("App ID: ", Appid)
	fmt.Println("Domain: ", Domain)
}
