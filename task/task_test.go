package task

import (
	"testing"
)

func TestParseEntrypointWithSpace(t *testing.T) {
	input := `test1 test2 test3`
	result, err := parseEntrypoint(input)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 3 {
		t.Error("len(result) expect = 3, but actual = %v", len(result))
	}

	if *result[0] != "test1" ||
		*result[1] != "test2" ||
		*result[2] != "test3" {
		t.Error("parse result is not expected")
	}
}

func TestParseEntrypointWithSpaceAndQuote(t *testing.T) {

	input := `nginx -g "daemon off;" -c /etc/nginx/nginx.conf`

	result, err := parseEntrypoint(input)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 5 {
		t.Error("len(result) expect = 5, but actual = %v", len(result))
	}

	if *result[0] != "nginx" ||
		*result[1] != "-g" ||
		*result[2] != "daemon off;" ||
		*result[3] != "-c" ||
		*result[4] != "/etc/nginx/nginx.conf" {
		t.Error("parse result is not expected")
	}
}

func TestCreateContainerDefinition(t *testing.T) {

	input := ContainerDefinition{
		Name:  "test-container",
		Image: "stormcat24/test:latest",
		Ports: []string{
			"80:8080",
		},
		Environment: map[string]string{
			"PARAM1": "VALUE1",
		},
		Links: []string{
			"mysql",
		},
		Volumes: []string{
			"/var/log/container/test:/var/log/test",
		},
		VolumesFrom: []string{
			"api",
		},
		Memory:            int64(256),
		CpuUnits:          int64(1024),
		Essential:         true,
		EntryPoint:        "entrypoint value",
		Command:           "command value",
		DisableNetworking: true,
		DnsSearchDomains: []string{
			"test.dns.domain",
		},
		DnsServers: []string{
			"test.dns.server",
		},
		DockerLabels: map[string]string{
			"LABEL1": "VALUE1",
		},
		DockerSecurityOptions: []string{
			"ECS_SELINUX_CAPABLE=true",
		},
		ExtraHosts: []HostEntry{
			HostEntry{
				Hostname:  "host1",
				IpAddress: "192.168.1.100",
			},
		},
	}

	con, volumes, _ := createContainerDefinition(&input)

	if input.Name != *con.Name {
		t.Errorf("Name: expect = %v, but actual = %v", input.Name, *con.Name)
	}

	if input.Image != *con.Image {
		t.Errorf("Image: expect = %v, but actual = %v", input.Image, *con.Image)
	}

	if len(input.Ports) != len(con.PortMappings) {
		t.Fatalf("len(PortMappings): expect = %v, but actual = %v", len(input.Ports), len(con.PortMappings))
	}

	if int64(8080) != *con.PortMappings[0].ContainerPort {
		t.Errorf("ContainerPort: expect = %v, but actual = %v", "8080", *con.PortMappings[0].ContainerPort)
	}

	if int64(80) != *con.PortMappings[0].HostPort {
		t.Errorf("HostPort: expect = %v, but actual = %v", "8080", *con.PortMappings[0].HostPort)
	}

	if len(input.Environment) != len(con.Environment) {
		t.Fatalf("len(Environment): expect = %v, but actual = %v", len(input.Environment), len(con.Environment))
	}

	if "PARAM1" != *con.Environment[0].Name {
		t.Errorf("Environment.Name: expect = %v, but actual = %v", "PARAM1", *con.Environment[0].Name)
	}

	if "VALUE1" != *con.Environment[0].Value {
		t.Errorf("Environment.Value: expect = %v, but actual = %v", "VALUE1", *con.Environment[0].Value)
	}

	if len(input.Links) != len(con.Links) {
		t.Fatalf("len(Links): expect = %v, but actual = %v", len(input.Links), len(con.Links))
	}

	if input.Links[0] != *con.Links[0] {
		t.Errorf("Link: expect = %v, but actual = %v", input.Links[0], *con.Links[0])
	}

	if len(input.Volumes) != len(volumes) {
		t.Fatalf("len(volumes): expect = %v, but actual = %v", len(input.Volumes), len(volumes))
	}

	if "VarLogContainerTest" != *volumes[0].Name {
		t.Errorf("Volumes.Name: expect = %v, but actual = %v", "VarLogContainerTest", *volumes[0].Name)
	}

	if "/var/log/container/test" != *volumes[0].Host.SourcePath {
		t.Errorf("Volumes.Host.SourcePath: expect = %v, but actual = %v", "/var/log/container/test", *volumes[0].Host.SourcePath)
	}

	if len(input.VolumesFrom) != len(con.VolumesFrom) {
		t.Fatalf("len(VolumesFrom): expect = %v, but actual = %v", len(input.VolumesFrom), len(con.VolumesFrom))
	}

	if input.VolumesFrom[0] != *con.VolumesFrom[0].SourceContainer {
		t.Errorf("VolumesFrom.SourceContainer: expect = %v, but actual = %v", input.VolumesFrom[0], *con.VolumesFrom[0].SourceContainer)
	}

	if input.Memory != *con.Memory {
		t.Errorf("Memory: expect = %v, but actual = %v", input.Memory, *con.Memory)
	}

	if input.CpuUnits != *con.Cpu {
		t.Errorf("CpuUnits: expect = %v, but actual = %v", input.CpuUnits, *con.Cpu)
	}

	if input.Essential != *con.Essential {
		t.Errorf("Essential: expect = %v, but actual = %v", input.Essential, *con.Essential)
	}

	if 2 != len(con.EntryPoint) {
		t.Fatalf("len(EntryPoint): expect = %v, but actual = %v", 1, len(con.EntryPoint))
	}

	if "entrypoint" != *con.EntryPoint[0] {
		t.Errorf("EntryPoint[0]: expect = %v, but actual = %v", "entrypoint", *con.EntryPoint[0])
	}

	if "value" != *con.EntryPoint[1] {
		t.Errorf("EntryPoint[1]: expect = %v, but actual = %v", "value", *con.EntryPoint[1])
	}

	if 2 != len(con.Command) {
		t.Fatalf("len(Command): expect = %v, but actual = %v", 1, len(con.Command))
	}

	if "command" != *con.Command[0] {
		t.Errorf("Command[0]: expect = %v, but actual = %v", "command", *con.Command[0])
	}

	if "value" != *con.Command[1] {
		t.Errorf("Command[1]: expect = %v, but actual = %v", "value", *con.Command[1])
	}

	if input.DisableNetworking != *con.DisableNetworking {
		t.Errorf("DisableNetworking: expect = %v, but actual = %v", input.DisableNetworking, *con.DisableNetworking)
	}

	if 1 != len(con.DnsSearchDomains) {
		t.Fatalf("len(DnsSearchDomains): expect = %v, but actual = %v", 1, len(con.DnsSearchDomains))
	}

	if input.DnsSearchDomains[0] != *con.DnsSearchDomains[0] {
		t.Errorf("DnsSearchDomains[0]: expect = %v, but actual = %v", input.DnsSearchDomains[0], *con.DnsSearchDomains[0])
	}

	if 1 != len(con.DnsServers) {
		t.Fatalf("len(DnsServers): expect = %v, but actual = %v", 1, len(con.DnsServers))
	}

	if input.DnsServers[0] != *con.DnsServers[0] {
		t.Errorf("DnsServers[0]: expect = %v, but actual = %v", input.DnsServers[0], *con.DnsServers[0])
	}

	if len(input.DockerLabels) != len(con.DockerLabels) {
		t.Fatalf("len(DockerLabels): expect = %v, but actual = %v", len(input.DockerLabels), len(con.DockerLabels))
	}

	if val, ok := con.DockerLabels["LABEL1"]; ok {
		if "VALUE1" != *val {
			t.Errorf("DockerLabels.LABEL1: expect = %v, but actual = %v", "VALUE1", val)
		}
	} else {
		t.Errorf("DockerLabels.LABEL1: not found")
	}

	if 1 != len(con.DockerSecurityOptions) {
		t.Fatalf("len(DockerSecurityOptions): expect = %v, but actual = %v", 1, len(con.DockerSecurityOptions))
	}

	if input.DockerSecurityOptions[0] != *con.DockerSecurityOptions[0] {
		t.Errorf("DockerSecurityOptions[0]: expect = %v, but actual = %v", input.DockerSecurityOptions[0], *con.DockerSecurityOptions[0])
	}

	if len(input.ExtraHosts) != len(con.ExtraHosts) {
		t.Fatalf("len(ExtraHosts): expect = %v, but actual = %v", len(input.ExtraHosts), len(con.ExtraHosts))
	}

	if input.ExtraHosts[0].Hostname != *con.ExtraHosts[0].Hostname {
		t.Errorf("ExtraHosts[0].Hostname: expect = %v, but actual = %v", input.ExtraHosts[0].Hostname, *con.ExtraHosts[0].Hostname)
	}

	if input.ExtraHosts[0].IpAddress != *con.ExtraHosts[0].IpAddress {
		t.Errorf("ExtraHosts[0].IpAddress: expect = %v, but actual = %v", input.ExtraHosts[0].IpAddress, *con.ExtraHosts[0].IpAddress)
	}

	if input.Hostname != *con.Hostname {
		t.Errorf("Hostname: expect = %v, but actual = %v", input.Hostname, *con.Hostname)
	}
}
