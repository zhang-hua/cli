package securitygroup

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry/cli/cf/api/security_groups"
	"github.com/cloudfoundry/cli/cf/command_metadata"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/flag_helpers"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type CreateSecurityGroup struct {
	ui                terminal.UI
	securityGroupRepo security_groups.SecurityGroupRepo
	configRepo        configuration.Reader
}

func NewCreateSecurityGroup(ui terminal.UI, configRepo configuration.Reader, securityGroupRepo security_groups.SecurityGroupRepo) CreateSecurityGroup {
	return CreateSecurityGroup{
		ui:                ui,
		configRepo:        configRepo,
		securityGroupRepo: securityGroupRepo,
	}
}

func (cmd CreateSecurityGroup) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "create-security-group",
		Description: "Create a security group",
		Usage:       "CF_NAME create-security-group NAME [--json PATH_TO_JSON_FILE]",
		Flags: []cli.Flag{
			flag_helpers.NewStringFlag("json", "Path to a file containing rules in JSON format"),
			flag_helpers.NewStringFlag("csv", "Path to a file containing rules in CSV format"),
		},
	}
}

func (cmd CreateSecurityGroup) GetRequirements(requirementsFactory requirements.Factory, context *cli.Context) ([]requirements.Requirement, error) {
	if len(context.Args()) != 1 {
		cmd.ui.FailWithUsage(context)
	}

	requirements := []requirements.Requirement{requirementsFactory.NewLoginRequirement()}
	return requirements, nil
}

func (cmd CreateSecurityGroup) Run(context *cli.Context) {
	name := context.Args()[0]
	pathToJSONFile := context.String("json")
	rules, err := cmd.parseJSON(pathToJSONFile)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	pathToCSVFile := context.String("csv")
	csvRules, err := cmd.parseCSV(pathToCSVFile)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	rules = append(rules, csvRules...)

	cmd.ui.Say("Creating security group %s as %s",
		terminal.EntityNameColor(name),
		terminal.EntityNameColor(cmd.configRepo.Username()))

	err = cmd.securityGroupRepo.Create(name, rules)

	httpErr, ok := err.(errors.HttpError)
	if ok && httpErr.ErrorCode() == errors.SECURITY_GROUP_EXISTS {
		cmd.ui.Ok()
		cmd.ui.Warn("Security group %s %s",
			terminal.EntityNameColor(name),
			terminal.WarningColor("already exists"))
		return
	}

	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Ok()
}

func (cmd CreateSecurityGroup) parseJSON(path string) ([]map[string]string, error) {
	if path == "" {
		return []map[string]string{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	ruleMaps := []map[string]string{}
	err = json.Unmarshal(bytes, &ruleMaps)
	if err != nil {
		cmd.ui.Failed("Incorrect json format: %s", err.Error())
	}

	return ruleMaps, nil
}

func (cmd CreateSecurityGroup) parseCSV(path string) ([]map[string]string, error) {
	if path == "" {
		return []map[string]string{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	rules := []map[string]string{}
	for _, record := range records {
		protocol := record[0]

		switch protocol {
		case "tcp", "udp":
			rules = append(rules, map[string]string{
				"protocol":    protocol,
				"destination": record[1],
				"ports":       record[2],
			})
		case "icmp":
			rules = append(rules, map[string]string{
				"protocol":    protocol,
				"destination": record[1],
				"type":        record[2],
				"code":        record[3],
			})
		case "all":
			rules = append(rules, map[string]string{
				"protocol":    protocol,
				"destination": record[1],
			})
		default:
			return nil, errors.New("Unknown protocol: " + protocol)
		}
	}

	return rules, nil
}
