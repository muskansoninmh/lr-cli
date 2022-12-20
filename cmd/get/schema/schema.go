package schema

import (
	"os"
	"sort"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

func NewschemaCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Gets schema",
		Long:  `Use this command to get the list of configured registration schema fields.`,
		Example: heredoc.Doc(`$ lr get schema
+-----------+---------------+----------+---------+
|   NAME    |    DISPLAY    |   TYPE   | ENABLED |
+-----------+---------------+----------+---------+
| password  | Password      | password | true    |
| emailid   | Email Id      | email    | true    |
| lastname  | Last Name     | string   | false   |
| birthdate | Date of Birth | string   | false   |
| country   | Country       | string   | false   |
| firstname | First Name    | string   | false   |
+-----------+---------------+----------+---------+
+---------------+----------+---------+
| CUSTOM FIELDS |   TYPE   | ENABLED |
+---------------+----------+---------+
| MyCF          | string   | false   |
+---------------+----------+---------+
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return get()
		},
	}

	return cmd
}

func get() error {

	features, err := api.GetSiteFeatures()
	if err != nil {
		return err
	}

	regFields, err := api.GetAllRegistrationFields()
	activeRegField, err := api.GetRegistrationFields()
	
	var data [][]string
	for k, v := range regFields {
		if k == "phoneid" && !api.IsPhoneLoginEnabled(*features) {
			continue
		}
		enabled := "false"
		_, ok := activeRegField[k]
		if ok {
			enabled = "true"
		}
		Type := v.Type
		if Type == "multi" {
			Type = "checkbox"
		}
		data = append(data, []string{k, v.Display, Type, enabled})
	}
	sort.SliceStable(data, func(i, j int) bool {
		return data[i][3] == "true"
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Display", "Type", "Enabled"})
	table.AppendBulk(data)
	table.Render()

	customFields, err := api.GetAllCustomFields()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	cfTable := tablewriter.NewWriter(os.Stdout)
	if len(customFields.Data) > 0 {
		for _, v := range customFields.Data {
			enabled := "false"
		_, ok := activeRegField["cf_" + v.Display]
		if ok {
			enabled = "true"
		}
		Type := activeRegField["cf_" + v.Display].Type
			cfTable.Append([]string{v.Display, Type,enabled})
		}
	} else {
		cfTable.Append([]string{"No Custom Fields"})
		cfTable.SetCaption(true, "Use command `lr add custom-field` to add the Custom Field")
	}
	cfTable.SetHeader([]string{"Custom Fields", "Type", "Enabled"})
	cfTable.Render()

	return nil
}
