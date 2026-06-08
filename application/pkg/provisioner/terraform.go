package provisioner

import (
	"fmt"
	"log"
	"os/exec"
)

func TriggerTerraform(serverCount int) {
	varFlag := fmt.Sprintf("-var=server_count=%d", serverCount)
	fmt.Println("Calculated servers:", serverCount)
	fmt.Println("Triggering terraform...")

	cmd := exec.Command("terraform", "plan", varFlag)
	//cmd := exec.Command("terraform", "apply", "-auto-approve", varFlag)
	output, err := cmd.CombinedOutput()

	fmt.Println(string(output))

	if err != nil {
		log.Fatal(err)
	}
}
