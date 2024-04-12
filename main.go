package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func contractName(addr common.Address) string {
	switch addr.String() {
	case ValidatorContract:
		return "ValidatorContract"
	case SlashContract:
		return "SlashContract"
	case SystemRewardContract:
		return "SystemRewardContract"
	case LightClientContract:
		return "LightClientContract"
	case TokenHubContract:
		return "TokenHubContract"
	case RelayerIncentivizeContract:
		return "RelayerIncentivizeContract"
	case RelayerHubContract:
		return "RelayerHubContract"
	case GovHubContract:
		return "GovHubContract"
	case TokenManagerContract:
		return "TokenManagerContract"
	case CrossChainContract:
		return "CrossChainContract"
	case StakingContract:
		return "StakingContract"
	case StakeHubContract:
		return "StakeHubContract"
	case StakeCreditContract:
		return "StakeCreditContract"
	case GovernorContract:
		return "GovernorContract"
	case GovTokenContract:
		return "GovTokenContract"
	case TimelockContract:
		return "TimelockContract"
	case TokenRecoverPortalContract:
		return "TokenRecoverPortalContract"
	default:
		return "Unknown"
	}
}

func writeTypes(fileName string, content []byte) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		panic(err)
	}
}

func extract(upgradeConfig map[string]*Upgrade, upgradeName string, replacement string) string {
	err := os.RemoveAll(upgradeName)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(fmt.Sprintf("%s/", upgradeName), os.ModePerm)
	if err != nil {
		panic(err)
	}

	writeTypes(upgradeName+"/types.go", []byte(fmt.Sprintf(`package %s
	
import _ "embed"
	`, upgradeName)))
	for network, upgrade := range upgradeConfig {
		writeTypes(upgradeName+"/types.go", []byte(fmt.Sprintf(`
// contract codes for %s upgrade
var (
`, network)))
		err := os.MkdirAll(fmt.Sprintf("%s/%s/", upgradeName, strings.ToLower(network)), os.ModePerm)
		if err != nil {
			panic(err)
		}

		for _, config := range upgrade.Configs {
			contractName := contractName(config.ContractAddr)
			err = os.WriteFile(fmt.Sprintf("%s/%s/%s", upgradeName, strings.ToLower(network), contractName),
				[]byte(config.Code), 0644)
			if err != nil {
				panic(err)
			}
			replacement = strings.ReplaceAll(replacement, "\""+config.Code+"\"", fmt.Sprintf("%s.%s%s", upgradeName, strings.Title(network), contractName))
			writeTypes(upgradeName+"/types.go", []byte(fmt.Sprintf(`	//go:embed %s/%s
	%s%s string
`, strings.ToLower(network), contractName, strings.Title(network), contractName)))

		}

		writeTypes(upgradeName+"/types.go", []byte(`)
`))
	}

	return replacement
}

func main() {
	err := os.RemoveAll("output")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll("output", os.ModePerm)
	if err != nil {
		panic(err)
	}
	writeTypes("output/upgrade.go", nil)

	content, err := os.ReadFile("upgrade.go")
	if err != nil {
		panic(err)
	}

	replacement := string(content)
	replacement = extract(ramanujanUpgrade, "ramanujan", replacement)
	replacement = extract(nielsUpgrade, "niels", replacement)
	replacement = extract(mirrorUpgrade, "mirror", replacement)
	replacement = extract(brunoUpgrade, "bruno", replacement)
	replacement = extract(eulerUpgrade, "euler", replacement)
	replacement = extract(gibbsUpgrade, "gibbs", replacement)
	replacement = extract(moranUpgrade, "moran", replacement)
	replacement = extract(planckUpgrade, "planck", replacement)
	replacement = extract(lubanUpgrade, "luban", replacement)
	replacement = extract(platoUpgrade, "plato", replacement)
	replacement = extract(keplerUpgrade, "kepler", replacement)
	replacement = extract(feynmanUpgrade, "feynman", replacement)
	replacement = extract(feynmanFixUpgrade, "feynmanFix", replacement)

	writeTypes("output/upgrade.go", []byte(replacement))
}
