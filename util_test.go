package client

import (
	"testing"
)

func TestGetMachineIp(t *testing.T) {
	test := func(input, expectedOutput string) {
		ip := GetMachineIP(input)

		if ip != expectedOutput {
			t.Fatalf("Expected '%s', but got '%s' for GetMachineIp('%s')", expectedOutput, ip, input)
		}
	}

	test("f715c360.../172.17.8.103", "172.17.8.103")
	test("-", "")
	test("F8C1750E-116F-455A-AAEE-6EF24B439687/255.255.255.255", "255.255.255.255")
}
