 package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTempCSV(tb testing.TB, content string) string {
	tb.Helper()
	dir := tb.TempDir()
	path := filepath.Join(dir, "test_data.csv")
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		tb.Fatalf("failed to create temp csv: %v", err)
	}
	return path
}

func TestAnalyzeCSV(t *testing.T) {
	csvContent := `campaign_id,date,impressions,clicks,spend,conversions
CMP1,2023-01-01,1000,10,100.0,2
CMP1,2023-01-02,1000,20,100.0,3
CMP2,2023-01-01,500,50,50.0,1
CMP3,2023-01-01,100,5,20.0,0
CMP4,2023-01-01,2000,200,500.0,10
CMP4,2023-01-02,2000,200,500.0,15
`
	// Expected calculations:
	// CMP1: 
	//   Impressions: 2000, Clicks: 30, Spend: 200.0, Conversions: 5
	//   CTR: 30/2000 = 0.015
	//   CPA: 200/5 = 40.0

	// CMP2:
	//   Impressions: 500, Clicks: 50, Spend: 50.0, Conversions: 1
	//   CTR: 50/500 = 0.10
	//   CPA: 50/1 = 50.0

	// CMP3 (Zero Conversions):
	//   Impressions: 100, Clicks: 5, Spend: 20.0, Conversions: 0
	//   CTR: 5/100 = 0.05
	//   CPA: NULL (should be filtered out of Top CPA or treated as N/A)

	// CMP4:
	//   Impressions: 4000, Clicks: 400, Spend: 1000.0, Conversions: 25
	//   CTR: 400/4000 = 0.10
	//   CPA: 1000/25 = 40.0

	path := createTempCSV(t, csvContent)

	a, err := New()
	if err != nil {
		t.Fatalf("failed to create analyzer: %v", err)
	}
	defer a.Close()

	result, err := a.AnalyzeCSV(path)
	if err != nil {
		t.Fatalf("AnalyzeCSV failed: %v", err)
	}

	// 1. Verify Top CTR 
	// Expected Order (Desc CTR): 
	// 1. CMP2 (0.10) - Same as CMP4 but let's check values provided by DuckDB (stable sort isn't guaranteed without tiebreaker, but let's check existence)
	// 2. CMP4 (0.10)
	// 3. CMP3 (0.05)
	// 4. CMP1 (0.015)
	
	if len(result.TopCTR) != 4 {
		t.Errorf("Expected 4 campaigns in TopCTR, got %d", len(result.TopCTR))
	} else {
		// Checking top values. Since CMP2 and CMP4 have same CTR, order might vary unless we added secondary sort.
		// Let's just verify the top 2 are CMP2/CMP4 with 0.10
		if result.TopCTR[0].CTR != 0.10 {
			t.Errorf("Expected TopCTR[0] to be 0.10, got %f", result.TopCTR[0].CTR)
		}
		if result.TopCTR[1].CTR != 0.10 {
			t.Errorf("Expected TopCTR[1] to be 0.10, got %f", result.TopCTR[1].CTR)
		}
		
		// CMP3
		if result.TopCTR[2].CampaignID != "CMP3" {
			t.Errorf("Expected TopCTR[2] to be CMP3, got %s", result.TopCTR[2].CampaignID)
		}
		if result.TopCTR[2].CTR != 0.05 {
			t.Errorf("Expected TopCTR[2] to be 0.05, got %f", result.TopCTR[2].CTR)
		}

		// CMP1
		if result.TopCTR[3].CampaignID != "CMP1" {
			t.Errorf("Expected TopCTR[3] to be CMP1, got %s", result.TopCTR[3].CampaignID)
		}
		if result.TopCTR[3].CTR != 0.015 {
			t.Errorf("Expected TopCTR[3] to be 0.015, got %f", result.TopCTR[3].CTR)
		}
	}

	// 2. Verify Top CPA (Ascending CPA)
	// CMP3 has 0 conversions, so it should NOT be in the Top CPA list (filtered by HAVING SUM(conversions) > 0)
	// Remaining:
	// CMP1: 40.0
	// CMP4: 40.0
	// CMP2: 50.0
	
	if len(result.TopCPA) != 3 {
		t.Errorf("Expected 3 campaigns in TopCPA (CMP3 excluded), got %d", len(result.TopCPA))
	} else {
		// First two should be CMP1/CMP4 with 40.0
		if *result.TopCPA[0].CPA != 40.0 {
			t.Errorf("Expected TopCPA[0] to be 40.0, got %f", *result.TopCPA[0].CPA)
		}
		if *result.TopCPA[1].CPA != 40.0 {
			t.Errorf("Expected TopCPA[1] to be 40.0, got %f", *result.TopCPA[1].CPA)
		}
		
		// Last one CMP2 with 50.0
		if result.TopCPA[2].CampaignID != "CMP2" {
			t.Errorf("Expected TopCPA[2] to be CMP2, got %s", result.TopCPA[2].CampaignID)
		}
		if *result.TopCPA[2].CPA != 50.0 {
			t.Errorf("Expected TopCPA[2] to be 50.0, got %f", *result.TopCPA[2].CPA)
		}
	}
}

func BenchmarkAnalyzeCSV(b *testing.B) {
	// Define decreasing dataset sizes
	sizes := []int{10000, 1000, 100, 10}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Rows_%d", size), func(b *testing.B) {
			// Generate CSV content
			var sb strings.Builder
			sb.WriteString("campaign_id,date,impressions,clicks,spend,conversions\n")
			
			for i := 0; i < size; i++ {
				// Cycle through 10 campaigns
				campaignID := fmt.Sprintf("CMP%d", i%10)
				sb.WriteString(fmt.Sprintf("%s,2023-01-01,1000,10,10.0,1\n", campaignID))
			}
			
			path := createTempCSV(b, sb.String())
			
			a, err := New()
			if err != nil {
				b.Fatalf("failed to create analyzer: %v", err)
			}
			defer a.Close()
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := a.AnalyzeCSV(path)
				if err != nil {
					b.Fatalf("AnalyzeCSV failed: %v", err)
				}
			}
		})
	}
}
