package main

import (
	"lynisreport/lynis"
	"strings"
	"testing"
)

const (
	testParse1 string = `report_version_major=1
report_version_minor=0
report_datetime_start=2022-04-05 13:36:19
auditor=[Not Specified]
lynis_version=3.0.7
os=Linux
container=0
systemd=1


plugins_enabled=0
hostid=37feb2a24d03136df71ae200121805f5f4d526aa
hostid2=6773c6c8cc73bc7d9dda14ceb321bef3a373faa3e6662ecdbdf9eae6520fb454
running_service_tool=systemctl

running_service[]=auditd
running_service[]=bluetooth
running_service[]=cups
running_service[]=dbus
running_service[]=getty@tty1
warning[]=NETW-2706|Couldn't find 2 responsive nameservers|-|-|
#warning[]=NETW-2705|Couldn't find 2 responsive nameservers|-|-|
warning[]=NETW-2707|Couldn't find 2 responsive nameservers|-|-|
running_service[]=getty@tty2
running_service[]=getty@tty3
suggestion[]=NETW-3202|Determine if protocol 'dccp' is really needed on this system|-|-|
running_service[]=getty@tty4
running_service[]=netdata
running_service[]=polkit
#running_service[]=prometheus-node-exporter
warning[]=NETW-2708|Couldn't find 2 responsive nameservers|-|-|
 # warning[]=NETW-2705|Couldn't find 2 responsive nameservers|-|-|
warning[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
suggestion[]=NETW-3200|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3201|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3203|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
`
)

// Test parsing report
func TestReportParse(t *testing.T) {
	report, err := lynis.CreateReport(strings.NewReader(testParse1))
	if err != nil {
		t.Errorf("error parsing report: %s", err)
	}

	// Check correct amount of tests were parsed
	if len(report.Tests) != 8 {
		t.Errorf("parsed %d tests wanted %d", len(report.Tests), 8)
	}
}

// Test parsing report for Elasticsearch
func TestReportParseElastic(t *testing.T) {
	report, err := lynis.CreateReport(strings.NewReader(testParse1))
	if err != nil {
		t.Errorf("error parsing report: %s", err)
	}

	tees, _ := report.CreateTestElementElastics()
	if len(tees) != 9 {
		t.Errorf("parsed %d test elements wanted %d", len(tees), 9)
	}
}

const (
	testParse2 string = `report_version_major=1
report_version_minor=0
report_datetime_start=2022-04-05 13:36:19
auditor=[Not Specified]
lynis_version=3.0.6
warning[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
suggestion[]=NETW-3200|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3201|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3203|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
`
	testParse3 string = `report_version_major=1
report_version_minor=0
report_datetime_start=2022-04-05 13:36:19
auditor=[Not Specified]
lynis_version=2.0.7
warning[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
suggestion[]=NETW-3200|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3201|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3203|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
`
	testParse4 string = `report_version_major=1
report_version_minor=0
report_datetime_start=2022-04-05 13:36:19
auditor=[Not Specified]
lynis_version=2.7
warning[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
suggestion[]=NETW-3200|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3201|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3203|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
`
	testParse5 string = `report_version_major=1
report_version_minor=0
report_datetime_start=2022-04-05 13:36:19
auditor=[Not Specified]
lynis_version=3.a.7
warning[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
suggestion[]=NETW-3200|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3201|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3203|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
`
)

// test with invalid lynis versions
func TestReportParseInvalidLynisVersion(t *testing.T) {
	_, err := lynis.CreateReport(strings.NewReader(testParse2))
	if err == nil {
		t.Errorf("expected error generated from parsing incorrect version")
	} else {
		want := "line 5: Lynis version is not compatable"
		if err.Error() != want {
			t.Errorf("expected error message %s got %s", want, err.Error())
		}
	}
	_, err = lynis.CreateReport(strings.NewReader(testParse3))
	if err == nil {
		t.Errorf("expected error generated from parsing incorrect version")
	} else {
		want := "line 5: Lynis version is not compatable"
		if err.Error() != want {
			t.Errorf("expected error message %s got %s", want, err.Error())
		}
	}
	_, err = lynis.CreateReport(strings.NewReader(testParse4))
	if err == nil {
		t.Errorf("expected error generated from parsing incorrect version")
	} else {
		want := "line 5: Lynis version string invalid format"
		if err.Error() != want {
			t.Errorf("expected error message %s got %s", want, err.Error())
		}
	}
	_, err = lynis.CreateReport(strings.NewReader(testParse5))
	if err == nil {
		t.Errorf("expected error generated from parsing incorrect version")
	} else {
		want := "line 5: Lynis version string contains non number"
		if err.Error() != want {
			t.Errorf("expected error message %s got %s", want, err.Error())
		}
	}
}

const(
        testParse6 string = `report_version_major=1
report_version_minor=0
report_datetime_start=2022-04-05T13:36:19
auditor=[Not Specified]
lynis_version=3.0.7
warning[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
suggestion[]=NETW-3200|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3201|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3203|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
`
)
// test with invalid timestamp
func TestReportParseInvalidTimeStamp(t *testing.T) {
	_, err := lynis.CreateReport(strings.NewReader(testParse6))
	if err == nil {
		t.Errorf("expected error generated from parsing incorrect version")
	} else {
                want := `line 3: parsing time "2022-04-05T13:36:19-0600" as "2006-01-02 15:04:05-0700": cannot parse "T13:36:19-0600" as " "`
		if err.Error() != want {
			t.Errorf("expected error message %s got %s", want, err.Error())
		}
	}
}

const(
        testParse7 string = `report_version_major=1
report_version_minor=0
report_datetime_start=2022-04-05 13:36:19
auditor=[Not Specified]
lynis_version=3.0.7
warning[]=
suggestion[]=NETW-3200|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3201|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-3203|Determine if protocol 'dccp' is really needed on this system|-|-|
suggestion[]=NETW-2709|Couldn't find 2 responsive nameservers|-|-|
`
)
// test with invalid amount of test fields
func TestReportParseInvalidAmountTestFields(t *testing.T) {
	_, err := lynis.CreateReport(strings.NewReader(testParse7))
	if err == nil {
		t.Errorf("expected error generated from parsing incorrect version")
	} else {
                want := `line 6: malformed line no test name or test is missing info`
		if err.Error() != want {
			t.Errorf("expected error message %s got %s", want, err.Error())
		}
	}
}


