package dataplex_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	dataplex "github.com/hashicorp/terraform-provider-google/google/services/dataplex"
	"reflect"
	"strings"
	"testing"
)

func TestNumberOfEntryLinkAspectsValidation(t *testing.T) {
	fieldName := "aspects"
	numbers_2 := make([]interface{}, 2)
	for i := 0; i < 2; i++ {
		numbers_2[i] = i
	}
	numbers_1 := make([]interface{}, 1)
	for i := 0; i < 1; i++ {
		numbers_1[i] = i
	}
	numbers_empty := make([]interface{}, 0)
	map_2 := make(map[string]interface{}, 2)
	for i := 0; i < 2; i++ {
		key := fmt.Sprintf("key%d", i)
		map_2[key] = i
	}
	map_1 := make(map[string]interface{}, 1)
	for i := 0; i < 1; i++ {
		key := fmt.Sprintf("key%d", i)
		map_1[key] = i
	}
	map_empty := make(map[string]interface{}, 0)

	testCases := []struct {
		name        string
		input       interface{}
		expectError bool
		errorMsg    string
	}{
		{"too many aspects in a slice", numbers_2, true, "The number of required aspects can be atmost 1."},
		{"max number of aspects in a slice", numbers_1, false, ""},
		{"min number of aspects in a slice", numbers_empty, false, ""},
		{"too many aspects in a map", map_2, true, "The number of required aspects can be atmost 1."},
		{"max number of aspects in a map", map_1, false, ""},
		{"min number of aspects in a map", map_empty, false, ""},
		{"a string is not a valid input", "some-random-string", true, "to be array"},
		{"nil is not a valid input", nil, true, "to be array"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := dataplex.NumberOfEntryLinkAspectsValidation(tc.input, fieldName)
			hasError := len(errors) > 0

			if hasError != tc.expectError {
				t.Fatalf("%s: NumberOfEntryLinkAspectsValidation() error expectation mismatch: got error = %v (%v), want error = %v", tc.name, hasError, errors, tc.expectError)
			}

			if tc.expectError && tc.errorMsg != "" {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Error(), tc.errorMsg) { // Check if error message contains the expected substring
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: NumberOfEntryLinkAspectsValidation() expected error containing %q, but got: %v", tc.name, tc.errorMsg, errors)
				}
			}
		})
	}
}

func TestEntryLinkProjectNumberValidation(t *testing.T) {
	fieldName := "some_field"
	testCases := []struct {
		name        string
		input       interface{}
		expectError bool
		errorMsg    string
	}{
		{"valid input", "projects/1234567890/locations/us-central1", false, ""},
		{"valid input with only number", "projects/987/stuff", false, ""},
		{"valid input with trailing slash content", "projects/1/a/b/c", false, ""},
		{"valid input minimal", "projects/1/a", false, ""},
		{"invalid input trailing slash only", "projects/555/", true, "has an invalid format"},
		{"invalid type - int", 123, true, `to be string, but got int`},
		{"invalid type - nil", nil, true, `to be string, but got <nil>`},
		{"invalid format - missing 'projects/' prefix", "12345/locations/us", true, "has an invalid format"},
		{"invalid format - project number starts with 0", "projects/0123/data", true, "has an invalid format"},
		{"invalid format - no project number", "projects//data", true, "has an invalid format"},
		{"invalid format - letters instead of number", "projects/abc/data", true, "has an invalid format"},
		{"invalid format - missing content after number/", "projects/123", true, "has an invalid format"},
		{"invalid format - empty string", "", true, "has an invalid format"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := dataplex.EntryLinkProjectNumberValidation(tc.input, fieldName)
			hasError := len(errors) > 0

			if hasError != tc.expectError {
				t.Fatalf("%s: EntryLinkProjectNumberValidation() error expectation mismatch: got error = %v (%v), want error = %v", tc.name, hasError, errors, tc.expectError)
			}

			if tc.expectError && tc.errorMsg != "" {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Error(), tc.errorMsg) { // Check if error message contains the expected substring
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: EntryLinkProjectNumberValidation() expected error containing %q, but got: %v", tc.name, tc.errorMsg, errors)
				}
			}
		})
	}
}

func TestEntryLinkAspectProjectNumberValidation(t *testing.T) {
	fieldName := "some_field"
	testCases := []struct {
		name        string
		input       interface{}
		expectError bool
		errorMsg    string
	}{
		{"valid input", "655216118709.global.some-one-p-aspect", false, ""},
		{"valid input minimal", "655216118709.global.a", false, ""},
		{"invalid input trailing dot only", "655216118709.global.", true, "has an invalid format"},
		{"invalid type - int", 456, true, `to be string, but got int`},
		{"invalid type - nil", nil, true, `to be string, but got <nil>`},
		{"invalid format - missing project number", ".global.some-one-p-aspect", true, "has an invalid format"},
		{"invalid format - wrong project number", "655123118709.global.some-one-p-aspect", true, "has an invalid format"},
		{"invalid format - missing dot", "655216118709globalsome-one-p-aspect", true, "has an invalid format"},
		{"invalid format - project id instead of number", "one-p-aspect-project.global.some-one-p-aspect", true, "has an invalid format"},
		{"invalid format - empty string", "", true, "has an invalid format"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := dataplex.EntryLinkAspectProjectNumberValidation(tc.input, fieldName)
			hasError := len(errors) > 0

			if hasError != tc.expectError {
				t.Fatalf("%s: EntryLinkAspectProjectNumberValidation() error expectation mismatch: got error = %v (%v), want error = %v", tc.name, hasError, errors, tc.expectError)
			}

			if tc.expectError && tc.errorMsg != "" {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Error(), tc.errorMsg) { // Check if error message contains the expected substring
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: EntryLinkAspectProjectNumberValidation() expected error containing %q, but got: %v", tc.name, tc.errorMsg, errors)
				}
			}
		})
	}
}

func TestAddEntryLinkAspectsToSet(t *testing.T) {
	testCases := []struct {
		name         string
		initialSet   map[string]struct{}
		aspectsInput interface{}
		expectedSet  map[string]struct{}
		expectError  bool
		errorMsg     string
	}{
		{"add to empty set", map[string]struct{}{}, []interface{}{map[string]interface{}{"aspect_type": "key1"}, map[string]interface{}{"aspect_type": "key2"}}, map[string]struct{}{"key1": {}, "key2": {}}, false, ""},
		{"add to existing set", map[string]struct{}{"existing": {}}, []interface{}{map[string]interface{}{"aspect_type": "key1"}}, map[string]struct{}{"existing": {}, "key1": {}}, false, ""},
		{"add duplicate keys", map[string]struct{}{}, []interface{}{map[string]interface{}{"aspect_type": "key1"}, map[string]interface{}{"aspect_type": "key1"}, map[string]interface{}{"aspect_type": "key2"}}, map[string]struct{}{"key1": {}, "key2": {}}, false, ""},
		{"input aspects is empty slice", map[string]struct{}{"existing": {}}, []interface{}{}, map[string]struct{}{"existing": {}}, false, ""},
		{"input aspects is nil", map[string]struct{}{"original": {}}, nil, map[string]struct{}{"original": {}}, false, ""},
		{"input aspects is wrong type", map[string]struct{}{}, "not a slice", map[string]struct{}{}, true, "AddEntryLinkAspectsToSet: input 'aspects' is not a []interface{}, got string"},
		{"item in slice is not a map", map[string]struct{}{}, []interface{}{"not a map"}, map[string]struct{}{}, true, "AddEntryLinkAspectsToSet: item at index 0 is not a map[string]interface{}, got string"},
		{"item map missing aspect_type", map[string]struct{}{}, []interface{}{map[string]interface{}{"wrong_key": "key1"}}, map[string]struct{}{}, true, "AddEntryLinkAspectsToSet: 'aspect_type' not found in aspect item at index 0"},
		{"aspect_type is not a string", map[string]struct{}{}, []interface{}{map[string]interface{}{"aspect_type": 123}}, map[string]struct{}{}, true, "AddEntryLinkAspectsToSet: 'aspect_type' in item at index 0 is not a string, got int"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			currentSet := make(map[string]struct{})
			for k, v := range tc.initialSet {
				currentSet[k] = v
			}

			err := dataplex.AddEntryLinkAspectsToSet(currentSet, tc.aspectsInput)

			if tc.expectError {
				if err == nil {
					t.Fatalf("%s: Expected an error, but got nil", tc.name)
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("%s: Expected error message containing %q, got %q", tc.name, tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("%s: Did not expect an error, but got: %v", tc.name, err)
				}
				if !reflect.DeepEqual(currentSet, tc.expectedSet) {
					t.Errorf("%s: AddEntryLinkAspectsToSet() result mismatch:\ngot:  %v\nwant: %v", tc.name, currentSet, tc.expectedSet)
				}
			}
		})
	}
}

func TestInverseTransformEntryLinkAspects(t *testing.T) {
	testCases := []struct {
		name             string
		resInput         map[string]interface{}
		expectedAspects  []interface{}
		expectNilAspects bool
		expectError      bool
		errorMsg         string
	}{
		{"aspects key is absent", map[string]interface{}{"otherKey": "value"}, nil, true, false, ""},
        {"aspects value is nil", map[string]interface{}{"aspects": nil}, nil, true, false, ""},
        {"aspects is empty map", map[string]interface{}{"aspects": map[string]interface{}{}}, []interface{}{}, false, false, ""},
        {"aspects with one entry", map[string]interface{}{"aspects": map[string]interface{}{"key1": map[string]interface{}{"data": map[string]interface{}{"foo": "bar"}}}}, []interface{}{map[string]interface{}{"aspectType": "key1", "data": "{\"foo\":\"bar\"}"}}, false, false, ""},
        {"aspects is wrong type (not map)", map[string]interface{}{"aspects": "not a map"}, nil, false, true, "InverseTransformEntryLinkAspects: 'aspects' field is not a map[string]interface{}, got string"},
        {"aspect value does not contain data", map[string]interface{}{"aspects": map[string]interface{}{"key1": map[string]interface{}{"wrong": "value"}}}, nil, false, true, "InverseTransformEntryLinkAspects: value for key 'key1' does not contain 'data' field"},
        {"aspect value is not a map", map[string]interface{}{"aspects": map[string]interface{}{"key1": "not a map value"}}, nil, false, true, "InverseTransformEntryLinkAspects: value for key 'key1' is not a map[string]interface{}, got string"},
        {"data in aspect value is not a map", map[string]interface{}{"aspects": map[string]interface{}{"key1": map[string]interface{}{"data": "not-a-map"}}}, nil, false, true, "InverseTransformEntryLinkAspects: 'data' field for key 'key1' is not a map[string]interface{}, got string"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resCopy := deepCopyMap(tc.resInput)
			originalAspectsBeforeCall := deepCopyValue(resCopy["aspects"])

			err := dataplex.InverseTransformEntryLinkAspects(resCopy)

			if tc.expectError {
				if err == nil {
					t.Fatalf("%s: Expected an error, but got nil", tc.name)
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("%s: Expected error message containing %q, got %q", tc.name, tc.errorMsg, err.Error())
				}
				if !reflect.DeepEqual(resCopy["aspects"], originalAspectsBeforeCall) {
					t.Errorf("%s: resCopy['aspects'] was modified during error case.\nBefore: %#v\nAfter: %#v", tc.name, originalAspectsBeforeCall, resCopy["aspects"])
				}
				return
			}

			if err != nil {
				t.Fatalf("%s: Did not expect an error, but got: %v", tc.name, err)
			}

			actualAspectsRaw, aspectsKeyExists := resCopy["aspects"]

			if tc.expectNilAspects {
				if aspectsKeyExists && actualAspectsRaw != nil {
					t.Errorf("%s: Expected 'aspects' to be nil or absent, but got: %#v", tc.name, actualAspectsRaw)
				}
				return
			}

			if !aspectsKeyExists {
				t.Fatalf("%s: Expected 'aspects' key in result map, but it was missing. Expected value: %#v", tc.name, tc.expectedAspects)
			}
			if actualAspectsRaw == nil && tc.expectedAspects != nil {
				t.Fatalf("%s: Expected 'aspects' to be non-nil, but got nil. Expected value: %#v", tc.name, tc.expectedAspects)
			}

			actualAspectsSlice, ok := actualAspectsRaw.([]interface{})
			if !ok {
				if tc.expectedAspects != nil || actualAspectsRaw != nil {
					t.Fatalf("%s: Expected 'aspects' to be []interface{}, but got %T. Value: %#v", tc.name, actualAspectsRaw, actualAspectsRaw)
				}
			}

			if actualAspectsSlice != nil {
				sortAspectSlice(actualAspectsSlice)
			}
			if tc.expectedAspects != nil {
				sortAspectSlice(tc.expectedAspects)
			}

			if !reflect.DeepEqual(actualAspectsSlice, tc.expectedAspects) {
				t.Errorf("%s: InverseTransformEntryLinkAspects() result mismatch:\ngot:  %#v\nwant: %#v", tc.name, actualAspectsSlice, tc.expectedAspects)
			}
		})
	}
}

func TestTransformEntryLinkAspects(t *testing.T) {
	testCases := []struct {
		name             string
		objInput         map[string]interface{}
		expectedAspects  map[string]interface{}
		expectNilAspects bool
		expectError      bool
		errorMsg         string
	}{
		{"aspects key is absent", map[string]interface{}{"otherKey": "value"}, nil, true, false, ""},
        {"aspects value is nil", map[string]interface{}{"aspects": nil}, nil, true, false, ""},
        {"aspects is empty slice", map[string]interface{}{"aspects": []interface{}{}}, map[string]interface{}{}, false, false, ""},
        {"aspects with one item", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectType": "key1", "data": map[string]interface{}{"foo": "bar"}}}}, map[string]interface{}{"key1": map[string]interface{}{"data": map[string]interface{}{"foo": "bar"}}}, false, false, ""},
        {"aspects with multiple items", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectType": "key1", "data": map[string]interface{}{"foo": "bar"}}, map[string]interface{}{"aspectType": "key2", "data": map[string]interface{}{"abc": "xyz"}}}}, map[string]interface{}{"key1": map[string]interface{}{"data": map[string]interface{}{"foo": "bar"}}, "key2": map[string]interface{}{"data": map[string]interface{}{"abc": "xyz"}}}, false, false, ""},
        {"aspects with duplicate aspectType", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectType": "key1", "data": map[string]interface{}{"foo": "bar"}}, map[string]interface{}{"aspectType": "key2", "data": map[string]interface{}{"abc": "xyz"}}, map[string]interface{}{"aspectType": "key1", "data": map[string]interface{}{"new": "baz"}}}}, map[string]interface{}{"key1": map[string]interface{}{"data": map[string]interface{}{"new": "baz"}}, "key2": map[string]interface{}{"data": map[string]interface{}{"abc": "xyz"}}}, false, false, ""},
        {"aspects is wrong type (not slice)", map[string]interface{}{"aspects": "not a slice"}, nil, false, true, "TransformAspects: 'aspects' field is not a []interface{}, got string"},
        {"item in slice is not a map", map[string]interface{}{"aspects": []interface{}{"not a map"}}, nil, false, true, "TransformEntryLinkAspects: item in 'aspects' slice at index 0 is not a map[string]interface{}, got string"},
        {"item map missing aspectType", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"wrongKey": "k1", "data": map[string]interface{}{}}}}, nil, false, true, "TransformEntryLinkAspects: 'aspectType' not found in aspect item at index 0"},
        {"aspectType is not a string", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectType": 123, "data": map[string]interface{}{}}}}, nil, false, true, "TransformEntryLinkAspects: 'aspectType' in item at index 0 is not a string, got int"},
        {"item map missing data", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectType": "key1"}}}, nil, false, true, "TransformEntryLinkAspects: 'data' not found in aspect item at index 0"},
        {"data is not a map", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectType": "key1", "data": "not-a-map"}}}, nil, false, true, "TransformEntryLinkAspects: 'data' in item at index 0 is not a map[string]interface{}, got string"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			objCopy := deepCopyMap(tc.objInput)
			originalAspectsBeforeCall := deepCopyValue(objCopy["aspects"])

			err := dataplex.TransformEntryLinkAspects(objCopy)

			if tc.expectError {
				if err == nil {
					t.Fatalf("%s: Expected an error, but got nil", tc.name)
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("%s: Expected error message containing %q, got %q", tc.name, tc.errorMsg, err.Error())
				}
				if !reflect.DeepEqual(objCopy["aspects"], originalAspectsBeforeCall) {
					t.Errorf("%s: objCopy['aspects'] was modified during error case.\nBefore: %#v\nAfter: %#v", tc.name, originalAspectsBeforeCall, objCopy["aspects"])
				}
				return
			}

			if err != nil {
				t.Fatalf("%s: Did not expect an error, but got: %v", tc.name, err)
			}

			actualAspectsRaw, aspectsKeyExists := objCopy["aspects"]

			if tc.expectNilAspects {
				if aspectsKeyExists && actualAspectsRaw != nil {
					t.Errorf("%s: Expected 'aspects' to be nil or absent, but got: %#v", tc.name, actualAspectsRaw)
				}
				return
			}

			if !aspectsKeyExists {
				t.Fatalf("%s: Expected 'aspects' key in result map, but it was missing. Expected value: %#v", tc.name, tc.expectedAspects)
			}
			if actualAspectsRaw == nil && tc.expectedAspects != nil {
				t.Fatalf("%s: Expected 'aspects' to be non-nil, but got nil. Expected value: %#v", tc.name, tc.expectedAspects)
			}

			actualAspectsMap, ok := actualAspectsRaw.(map[string]interface{})
			if !ok {
				if tc.expectedAspects != nil || actualAspectsRaw != nil {
					t.Fatalf("%s: Expected 'aspects' to be map[string]interface{}, but got %T. Value: %#v", tc.name, actualAspectsRaw, actualAspectsRaw)
				}
			}

			if !reflect.DeepEqual(actualAspectsMap, tc.expectedAspects) {
				t.Errorf("%s: TransformAspects() result mismatch:\ngot:  %#v\nwant: %#v", tc.name, actualAspectsMap, tc.expectedAspects)
			}
		})
	}
}

func TestAccDataplexEntryLink_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
    ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexEntryLink_dataplexEntryLinkUpdate(context),
			},
			{
				ResourceName:            "google_dataplex_entry_link.basic_entry_link",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entry_group_id", "entry_link_id", "location"},
			},
			{
				ResourceName:            "google_dataplex_entry_link.full_entry_link",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entry_group_id", "entry_link_id", "location"},
			},
		},
	})
}

func testAccDataplexEntryLink_dataplexEntryLinkUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_entry_group" "entry-group-basic" {
  location = "us-central1"
  entry_group_id = "tf-test-entry-group%{random_suffix}"
  project = "%{project_number}"
}
resource "google_dataplex_entry" "source" {
  location = "us-central1"
  entry_group_id = google_dataplex_entry_group.entry-group-basic.entry_group_id
  entry_id = "tf-test-source-entry%{random_suffix}"
  entry_type = google_dataplex_entry_type.entry-type-basic.name
  project = "%{project_number}"
}
resource "google_dataplex_entry_type" "entry-type-basic" {
  entry_type_id = "tf-test-entry-type%{random_suffix}"
  location = "us-central1"
  project = "%{project_number}"
}
resource "google_dataplex_glossary" "term_test_id_full" {
  glossary_id = "tf-test-glossary%{random_suffix}"
  location    = "us-central1"
}
resource "google_dataplex_glossary_term" "term_test_id_full" {
  parent = "projects/${google_dataplex_glossary.term_test_id_full.project}/locations/us-central1/glossaries/${google_dataplex_glossary.term_test_id_full.glossary_id}"
  glossary_id = google_dataplex_glossary.term_test_id_full.glossary_id
  location = "us-central1"
  term_id = "tf-test-term-full%{random_suffix}"
  labels = { "tag": "test-tf" }
  display_name = "terraform term"
  description = "term created by Terraform"
}
resource "google_dataplex_entry_link" "basic_entry_link" {
  project = "%{project_number}"
  location = "us-central1"
  entry_group_id = google_dataplex_entry_group.entry-group-basic.entry_group_id
  entry_link_id = "tf-test-entry-link%{random_suffix}"
  entry_link_type = "projects/655216118709/locations/global/entryLinkTypes/definition"
  entry_references {
    name = google_dataplex_entry.source.name
	type = "SOURCE"
  }
  entry_references {
    name = "projects/${google_dataplex_entry_group.entry-group-basic.project}/locations/us-central1/entryGroups/@dataplex/entries/projects/${google_dataplex_entry_group.entry-group-basic.project}/locations/us-central1/glossaries/${google_dataplex_glossary.term_test_id_full.glossary_id}/terms/${google_dataplex_glossary_term.term_test_id_full.term_id}"
	type = "TARGET"
  }
}
resource "google_bigquery_dataset" "bq_dataset" {
  dataset_id = "tf_test_dataset_%{random_suffix}"
  project    = "%{project_number}"
  location   = "us-central1"
}

resource "google_bigquery_table" "table1" {
  dataset_id = google_bigquery_dataset.bq_dataset.dataset_id
  table_id   = "table1_%{random_suffix}"
  project    = "%{project_number}"
  schema     = <<EOF
[
  {
    "name": "col1",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "Column 1"
  }
]
EOF
}

resource "google_bigquery_table" "table2" {
  dataset_id = google_bigquery_dataset.bq_dataset.dataset_id
  table_id   = "table2_%{random_suffix}"
  project    = "%{project_number}"
  schema     = <<EOF
[
  {
    "name": "colA",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "Column A"
  }
]
EOF
}
resource "time_sleep" "wait_180s_for_dataplex_ingestion" {
  depends_on = [
    google_bigquery_table.table1,
    google_bigquery_table.table2,
  ]
  create_duration = "180s"
}
resource "google_dataplex_entry_link" "full_entry_link" {
  depends_on = [time_sleep.wait_180s_for_dataplex_ingestion]
  project = "%{project_number}"
  location = "us-central1"
  entry_group_id = "@bigquery"
  entry_link_id = "tf-test-full-entry-link%{random_suffix}"
  entry_link_type = "projects/655216118709/locations/global/entryLinkTypes/schema-join"
  entry_references {
    name = "projects/%{project_number}/locations/us-central1/entryGroups/@bigquery/entries/bigquery.googleapis.com/projects/%{project_number}/datasets/${google_bigquery_dataset.bq_dataset.dataset_id}/tables/${google_bigquery_table.table1.table_id}"
    type = ""
  }
  entry_references {
    name = "projects/%{project_number}/locations/us-central1/entryGroups/@bigquery/entries/bigquery.googleapis.com/projects/%{project_number}/datasets/${google_bigquery_dataset.bq_dataset.dataset_id}/tables/${google_bigquery_table.table2.table_id}"
    type = ""
  }
  aspects {
    aspect_type = "655216118709.global.schema-join"
    data = <<EOF
{
  "joins": [],
  "userManaged": true
}
EOF
  }
}
`, context)
}
