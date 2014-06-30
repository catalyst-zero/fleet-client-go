package client

import (
	"testing"
)

func TestParser(t *testing.T) {
	output := `UNIT						LOAD	ACTIVE	SUB	DESC							MACHINE
1de9403c-2592-4b01-af17-6e7f06ae05de.service	-	-	-	stm-session-service Sidekick				-
2ecf3fc0-b541-41ce-85b9-4958e5e30e1c.service	loaded	active	running	stm-session-redis Service				f715c360.../172.17.8.103
3eb35412-7999-48b7-a545-3f28e5245105.service	loaded	active	running	stm-api-service:stm-session-service Ambassador		f715c360.../172.17.8.103
52468430-9c81-4bec-9cc0-e144b6377d59.service	loaded	active	running	stm-api-redis Sidekick					f715c360.../172.17.8.103
5a42cf23-773d-4303-932f-a226a0f6b3a4.service	loaded	active	running	stm-api-redis Service					f715c360.../172.17.8.103
71f3a18b-316b-475d-bc28-49e374c255cf.service	loaded	active	running	stm-api-redis Sidekick					f715c360.../172.17.8.103
7b1a73fa-5991-4671-be8a-662ce9c5c211.service	loaded	active	running	stm-api-service Service					f715c360.../172.17.8.103
8451b50e-4cda-4baf-8b33-f9f579a44912.service	loaded	active	running	stm-api-redis Service					f715c360.../172.17.8.103
b99917a1-48ac-4d16-8fc1-45127d323cfe.service	loaded	active	running	stm-session-service:stm-session-redis Ambassador	f715c360.../172.17.8.103
cca374f0-6eef-4547-a80b-1762fc3e3118.service	loaded	active	running	stm-session-redis Sidekick				f715c360.../172.17.8.103
e82a5ab7-e388-4394-b51f-88c118e89780.service	loaded	failed	failed	stm-session-service Service				f715c360.../172.17.8.103
fc80556b-0f2a-46bf-b580-76d2228cb5bd.service	loaded	active	running	stm-api-service:stm-api-redis Ambassador		f715c360.../172.17.8.103
`
	status, err := parseFleetStatusOutput(output)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(status) != 12 {
		t.Fatalf("Invalid number of status objects returned, expected 12, got: %d", len(status))
	}

	testUnit := func(status UnitStatus, expectedUnitName, expectedLoad, expectedActive, expectedSub, expectedDesc, expectedMachine, expectedMachineIP string) {
		if status.Unit != expectedUnitName {
			t.Fatalf("Unexpected unit name: %s", status.Unit)
		}
		if status.Load != expectedLoad {
			t.Fatalf("Unexpected unit name: %s", status.Load)
		}
		if status.Active != expectedActive {
			t.Fatalf("Unexpected unit name: %s", status.Active)
		}
		if status.Sub != expectedSub {
			t.Fatalf("Unexpected unit name: %s", status.Sub)
		}
		if status.Description != expectedDesc {
			t.Fatalf("Unexpected description: %s", status.Description)
		}
		if status.Machine != expectedMachine {
			t.Fatalf("Unexpected machine: %s", status.Machine)
		}
		if status.MachineIP() != expectedMachineIP {
			t.Fatalf("Unexpected IP: %s", status.MachineIP())
		}
	}

	testUnit(
		status[0],
		"1de9403c-2592-4b01-af17-6e7f06ae05de.service",
		"-",
		"-",
		"-",
		"stm-session-service Sidekick",
		"-",
		"",
	)
	testUnit(
		status[6],
		"7b1a73fa-5991-4671-be8a-662ce9c5c211.service",
		"loaded",
		"active",
		"running",
		"stm-api-service Service",
		"f715c360.../172.17.8.103",
		"172.17.8.103",
	)
}
