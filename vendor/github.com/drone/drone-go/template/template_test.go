package template

import (
	"testing"

	"github.com/drone/drone-go/drone"
)

var tests = []struct {
	Payload *drone.Payload
	Input   string
	Output  string
}{
	{
		&drone.Payload{Build: &drone.Build{
			Commit: "0a266f42a9aef9db97a005ab46f6c53890339a9c"},
		},
		"{{ truncate build.commit 8 }}",
		"0a266f42",
	},
	{
		&drone.Payload{Build: &drone.Build{Number: 1}},
		"build #{{build.number}}",
		"build #1",
	},
	{
		&drone.Payload{Build: &drone.Build{Status: drone.StatusSuccess}},
		"{{uppercase build.status}}",
		"SUCCESS",
	},
	{
		&drone.Payload{Build: &drone.Build{Author: "Octocat"}},
		"{{lowercase build.author}}",
		"octocat",
	},
	{
		&drone.Payload{Build: &drone.Build{Status: drone.StatusSuccess}},
		"{{uppercasefirst build.status}}",
		"Success",
	},
	{
		&drone.Payload{Build: &drone.Build{
			Started:  1448127131,
			Finished: 1448127505},
		},
		"{{ duration build.started_at build.finished_at }}",
		"6m14s",
	},
	{
		&drone.Payload{Build: &drone.Build{Finished: 1448127505}},
		`finished at {{ datetime build.finished_at "3:04PM" "UTC" }}`,
		"finished at 5:38PM",
	},
	// verify the success if / else block works
	{
		&drone.Payload{Build: &drone.Build{Status: drone.StatusSuccess}},
		"{{#success build.status}}SUCCESS{{/success}}",
		"SUCCESS",
	},
	{
		&drone.Payload{Build: &drone.Build{Status: drone.StatusFailure}},
		"{{#success build.status}}SUCCESS{{/success}}",
		"",
	},
	{
		&drone.Payload{Build: &drone.Build{Status: drone.StatusFailure}},
		"{{#success build.status}}SUCCESS{{else}}NOT SUCCESS{{/success}}",
		"NOT SUCCESS",
	},
	// verify the failure if / else block works
	{
		&drone.Payload{Build: &drone.Build{Status: drone.StatusFailure}},
		"{{#failure build.status}}FAILURE{{/failure}}",
		"FAILURE",
	},
	{
		&drone.Payload{Build: &drone.Build{Status: drone.StatusSuccess}},
		"{{#failure build.status}}FAILURE{{/failure}}",
		"",
	},
	{
		&drone.Payload{Build: &drone.Build{Status: drone.StatusSuccess}},
		"{{#failure build.status}}FAILURE{{else}}NOT FAILURE{{/failure}}",
		"NOT FAILURE",
	},
	{
		&drone.Payload{Build: &drone.Build{Author: "url&unsafe=author!"}},
		"{{#urlencode}}google https://www.google.co.jp/ {{{build.author}}}{{/urlencode}}",
		"google+https%3A%2F%2Fwww.google.co.jp%2F+url%26unsafe%3Dauthor%21",
	},
}

func TestTemplate(t *testing.T) {

	for _, test := range tests {
		got, err := RenderTrim(test.Input, test.Payload)
		if err != nil {
			t.Errorf("Failed rendering template %q, got error %s.", test.Input, err)
		}
		if got != test.Output {
			t.Errorf("Wanted rendered template %q, got %q", test.Output, got)
		}
	}
}
