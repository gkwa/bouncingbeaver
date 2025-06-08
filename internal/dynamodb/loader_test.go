package dynamodb

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gkwa/bouncingbeaver/internal/testutil"
)

func TestLoadData_FromString(t *testing.T) {
	client := NewClient()

	jsonInput := `{
		"Items": [
			{
				"id": {"S": "test-id"},
				"name": {"S": "Test Product"},
				"price": {"S": "$9.99"}
			}
		],
		"Count": 1,
		"ScannedCount": 1
	}`

	reader := strings.NewReader(jsonInput)
	items, err := client.LoadDataFromReader(reader)
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	if items[0]["id"].(*types.AttributeValueMemberS).Value != "test-id" {
		t.Errorf("Expected id 'test-id', got %v", items[0]["id"])
	}
}

func TestLoadData_FromFile(t *testing.T) {
	client := NewClient()

	items, err := client.LoadData("testdata/sample_input.json")
	if err != nil {
		t.Fatalf("Failed to load data from file: %v", err)
	}

	if len(items) == 0 {
		t.Error("Expected items to be loaded")
	}
}

func TestUnmarshalProducts_Golden(t *testing.T) {
	client := NewClient()

	// Load test data
	items, err := client.LoadData("testdata/sample_input.json")
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	// Unmarshal to products
	products, err := client.UnmarshalProducts(items)
	if err != nil {
		t.Fatalf("Failed to unmarshal products: %v", err)
	}

	// Convert to JSON for golden file comparison
	actual := testutil.ToJSON(t, products)

	// Try to read golden file, if it doesn't exist, create it
	expected, err := testutil.GoldenSafe(t, actual, "products_output.golden")
	if err != nil {
		t.Logf("Golden file created: %v", err)
		return
	}

	if string(actual) != string(expected) {
		testutil.AssertGoldenMatch(t, actual, expected, "products_output.golden")
	}
}
