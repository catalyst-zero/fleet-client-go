package client

import (
	"testing"
)

const (
	OUTPUT_1 = `UNIT						STATE		LOAD	ACTIVE	SUB	DESC				MACHINE
13d1395a-65fd-47c8-8d80-4d31da4499f6.service	inactive	-	-	-	Presence redis Service		-
0c4f0407-b5f8-4a6d-9b8b-0172db5b105f.service	launched	loaded	active	running	Ambassador api:redis Ambassador	26f4636c.../172.17.8.101
2703648d-6b45-4d74-9e15-7c0de5c63ed1.service	launched	loaded	active	running	Presence api Service		26f4636c.../172.17.8.101
3d29c681-014e-433c-952c-590a83c1d1e7.service	launched	loaded	active	exited	Storage redis Storage		26f4636c.../172.17.8.101
6e301e0f-75c7-44d3-b0fe-e7f6787ef4b3.service	launched	loaded	active	running	Presence redis Service		26f4636c.../172.17.8.101
de3f0ecd-68b8-47ad-9265-96275e25fd1e.service	launched	loaded	active	running	User api Service		26f4636c.../172.17.8.101
df7df57f-f6ca-4c10-bcfa-b5164baff5a3.service	launched	loaded	active	running	User redis Service		26f4636c.../172.17.8.101
`
)

func AssertStatusParsedAs(t *testing.T, status UnitStatus,
	expectedUnitName, expectedState, expectedLoad, expectedActive,
	expectedSub, expectedDesc, expectedMachine, expectedMachineIP string) {
	if status.Unit != expectedUnitName {
		t.Fatalf("Unexpected unit name: %s", status.Unit)
	}
	if status.State != expectedState {
		t.Fatalf("Unexpected state: %s, expected: %s", status.State, expectedState)
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

func TestParser__OUTPUT_1(t *testing.T) {
	status, err := parseFleetStatusOutput(OUTPUT_1)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(status) != 7 {
		t.Fatalf("Invalid number of status objects returned, expected 7, got: %d", len(status))
	}
	AssertStatusParsedAs(t,
		status[0],
		"13d1395a-65fd-47c8-8d80-4d31da4499f6.service",
		"inactive",
		"-",
		"-",
		"-",
		"Presence redis Service",
		"-",
		"",
	)
	AssertStatusParsedAs(t,
		status[2],
		"2703648d-6b45-4d74-9e15-7c0de5c63ed1.service",
		"launched",
		"loaded",
		"active",
		"running",
		"Presence api Service",
		"26f4636c.../172.17.8.101",
		"172.17.8.101",
	)
	AssertStatusParsedAs(t,
		status[3],
		"3d29c681-014e-433c-952c-590a83c1d1e7.service",
		"launched",
		"loaded",
		"active",
		"exited",
		"Storage redis Storage",
		"26f4636c.../172.17.8.101",
		"172.17.8.101",
	)
}
