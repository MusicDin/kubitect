package actions

import (
	"fmt"
	"strings"
)

// listClusters lists all clusters located in the project clusters directory
// that contain terraform state file.
func ListClusters() error {
	clusters, err := ReadClustersInfo()

	if err != nil {
		return err
	}

	if len(clusters) == 0 {
		fmt.Println("No clusters initialized yet. Run 'kubitect apply' to create the cluster.")
		return nil
	}

	fmt.Println("Clusters:")

	for _, c := range clusters {
		var opt []string

		if c.Active() {
			opt = append(opt, "active")
		}

		if c.Local {
			opt = append(opt, "local")
		}

		if len(opt) > 0 {
			fmt.Printf("  - %s (%s)\n", c.Name, strings.Join(opt, ", "))
		} else {
			fmt.Printf("  - %s\n", c.Name)
		}
	}

	return nil
}
