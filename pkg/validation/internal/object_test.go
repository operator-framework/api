package internal

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func TestValidateObject(t *testing.T) {
	var table = []struct {
		description string
		path        string
		error       bool
		warning     bool
		detail      string
	}{
		{
			description: "valid PDB",
			path:        "./testdata/objects/valid_pdb.yaml",
		},
		{
			description: "invalid PDB - minAvailable set to 100%",
			path:        "./testdata/objects/invalid_pdb_minAvailable.yaml",
			error:       true,
			detail:      "minAvailable field cannot be set to 100%",
		},
		{
			description: "invalid PDB - maxUnavailable set to 0",
			path:        "./testdata/objects/invalid_pdb_maxUnavailable.yaml",
			error:       true,
			detail:      "maxUnavailable field cannot be set to 0 or 0%",
		},
		{
			description: "valid priorityclass",
			path:        "./testdata/objects/valid_priorityclass.yaml",
		},
		{
			description: "invalid priorityclass - global default set to true",
			path:        "./testdata/objects/invalid_priorityclass.yaml",
			error:       true,
			detail:      "globalDefault field cannot be set to true",
		},
		{
			description: "valid pdb role",
			path:        "./testdata/objects/valid_role_get_pdb.yaml",
		},
		{
			description: "invalid role - modify pdb",
			path:        "./testdata/objects/invalid_role_create_pdb.yaml",
			warning:     true,
			detail:      "RBAC includes permission to create/update poddisruptionbudgets, which could impact cluster stability",
		},
		{
			description: "valid scc role",
			path:        "./testdata/objects/valid_role_get_scc.yaml",
		},
		{
			description: "invalid scc role - modify default scc",
			path:        "./testdata/objects/invalid_role_modify_scc.yaml",
			error:       true,
			detail:      "RBAC includes permission to modify default securitycontextconstraints, which could impact cluster stability",
		},
	}

	for _, tt := range table {
		t.Log(tt.description)

		u := unstructured.Unstructured{}
		o, err := ioutil.ReadFile(tt.path)
		if err != nil {
			t.Fatalf("reading yaml object file: %s", err)
		}
		if err := yaml.Unmarshal(o, &u); err != nil {
			t.Fatalf("unmarshalling object at path %s: %v", tt.path, err)
		}

		results := ObjectValidator.Validate(&u)

		// check errors
		if len(results[0].Errors) > 0 && tt.error == false {
			t.Fatalf("received errors %#v when no validation error expected for %s", results, tt.path)
		}
		if len(results[0].Errors) == 0 && tt.error == true {
			t.Fatalf("received no errors when validation error expected for %s", tt.path)
		}
		if len(results[0].Errors) > 0 {
			if results[0].Errors[0].Detail != tt.detail {
				t.Fatalf("expected validation error detail %s, got %s", tt.detail, results[0].Errors[0].Detail)
			}
		}

		// check warnings
		if len(results[0].Warnings) > 0 && tt.warning == false {
			t.Fatalf("received errors %#v when no validation warning expected for %s", results, tt.path)
		}
		if len(results[0].Warnings) == 0 && tt.warning == true {
			t.Fatalf("received no errors when validation warning expected for %s", tt.path)
		}
		if len(results[0].Warnings) > 0 {
			if results[0].Warnings[0].Detail != tt.detail {
				t.Fatalf("expected validation warning detail %s, got %s", tt.detail, results[0].Warnings[0].Detail)
			}
		}
	}

}
