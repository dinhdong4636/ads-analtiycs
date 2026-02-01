package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/dinhdong4636/ads-analytics/internal/analyzer"
	"github.com/dinhdong4636/ads-analytics/internal/models"
)

const defaultDataPath = "ad_data.csv.zip"

func main() {
	csvPath := flag.String("csv", "", "Path to CSV or ZIP file (default: ad_data.csv.zip)")
	flag.Parse()

	if *csvPath == "" {
		*csvPath = defaultDataPath
		log.Printf("No -csv flag provided, using default: %s", *csvPath)
	}

	if _, err := os.Stat(*csvPath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", *csvPath)
	}

	finalCSVPath, cleanup, err := prepareCSVFile(*csvPath)
	if err != nil {
		log.Fatalf("Failed to prepare CSV file: %v", err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	a, err := analyzer.New()
	if err != nil {
		log.Fatalf("Failed to create analyzer: %v", err)
	}
	defer a.Close()

	log.Printf("Analyzing CSV file: %s", finalCSVPath)
	start := time.Now()
	
	result, err := a.AnalyzeCSV(finalCSVPath)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}
	
	duration := time.Since(start)
	log.Printf("Analysis completed in %v", duration)

	if err := exportToCSV(result); err != nil {
		log.Fatalf("Failed to export CSV: %v", err)
	}

	outputTextResults(result)
}

func prepareCSVFile(path string) (string, func(), error) {
	ext := strings.ToLower(filepath.Ext(path))
	
	if ext == ".csv" {
		return path, nil, nil
	}
	
	if ext == ".zip" {
		log.Printf("Detected ZIP file, extracting...")
		csvPath, err := extractZipToTemp(path)
		if err != nil {
			return "", nil, fmt.Errorf("failed to extract ZIP: %w", err)
		}
		
		cleanup := func() {
			os.Remove(csvPath)
			log.Printf("Cleaned up temporary file: %s", csvPath)
		}
		
		return csvPath, cleanup, nil
	}
	
	return "", nil, fmt.Errorf("unsupported file type: %s (only .csv and .zip supported)", ext)
}

func extractZipToTemp(zipPath string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	if len(r.File) == 0 {
		return "", fmt.Errorf("ZIP file is empty")
	}

	var csvFile *zip.File
	for _, f := range r.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".csv") {
			csvFile = f
			break
		}
	}

	if csvFile == nil {
		return "", fmt.Errorf("no CSV file found in ZIP archive")
	}

	log.Printf("Found CSV in ZIP: %s", csvFile.Name)

	rc, err := csvFile.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	tmpFile, err := os.CreateTemp("", "ad_data_*.csv")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, rc)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	log.Printf("Extracted to: %s", tmpFile.Name())
	return tmpFile.Name(), nil
}

func exportToCSV(result *models.AnalysisResult) error {
	if err := os.MkdirAll("export", 0755); err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")

	ctrFile := fmt.Sprintf("export/top_ctr_%s.csv", timestamp)
	if err := writeCSV(ctrFile, result.TopCTR, "Top CTR"); err != nil {
		return err
	}
	log.Printf("Exported Top CTR to: %s", ctrFile)

	cpaFile := fmt.Sprintf("export/top_cpa_%s.csv", timestamp)
	if err := writeCSV(cpaFile, result.TopCPA, "Top CPA"); err != nil {
		return err
	}
	log.Printf("Exported Top CPA to: %s", cpaFile)

	return nil
}

func writeCSV(filename string, metrics []models.CampaignMetrics, title string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	_, err = file.WriteString("campaign_id,total_impressions,total_clicks,total_spend,total_conversions,ctr,cpa\n")
	if err != nil {
		return err
	}

	for _, m := range metrics {
		cpa := ""
		if m.CPA != nil {
			cpa = fmt.Sprintf("%.2f", *m.CPA)
		}
		line := fmt.Sprintf("%s,%d,%d,%.2f,%d,%.6f,%s\n",
			m.CampaignID,
			m.TotalImpressions,
			m.TotalClicks,
			m.TotalSpend,
			m.TotalConversions,
			m.CTR,
			cpa,
		)
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}


func outputTextResults(result *models.AnalysisResult) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TOP 10 CAMPAIGNS BY CTR (Click-Through Rate)")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-15s %15s %12s %15s %15s %10s %10s\n",
		"Campaign ID", "Impressions", "Clicks", "Spend", "Conversions", "CTR", "CPA")
	fmt.Println(strings.Repeat("-", 80))
	
	for i, m := range result.TopCTR {
		cpa := "N/A"
		if m.CPA != nil {
			cpa = fmt.Sprintf("$%.2f", *m.CPA)
		}
		fmt.Printf("%-15s %15d %12d $%14.2f %15d %9.4f%% %10s\n",
			m.CampaignID, m.TotalImpressions, m.TotalClicks, 
			m.TotalSpend, m.TotalConversions, m.CTR*100, cpa)
		
		if i < len(result.TopCTR)-1 {
			fmt.Println()
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TOP 10 CAMPAIGNS BY CPA (Cost Per Acquisition)")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-15s %15s %12s %15s %15s %10s %10s\n",
		"Campaign ID", "Impressions", "Clicks", "Spend", "Conversions", "CTR", "CPA")
	fmt.Println(strings.Repeat("-", 80))
	
	for i, m := range result.TopCPA {
		cpa := "N/A"
		if m.CPA != nil {
			cpa = fmt.Sprintf("$%.2f", *m.CPA)
		}
		fmt.Printf("%-15s %15d %12d $%14.2f %15d %9.4f%% %10s\n",
			m.CampaignID, m.TotalImpressions, m.TotalClicks, 
			m.TotalSpend, m.TotalConversions, m.CTR*100, cpa)
		
		if i < len(result.TopCPA)-1 {
			fmt.Println()
		}
	}
	fmt.Println(strings.Repeat("=", 80))
}
