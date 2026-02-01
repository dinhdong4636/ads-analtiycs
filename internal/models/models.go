package models

import "time"

type CampaignRecord struct {
	CampaignID  string    `json:"campaign_id"`
	Date        time.Time `json:"date"`
	Impressions int       `json:"impressions"`
	Clicks      int       `json:"clicks"`
	Spend       float64   `json:"spend"`
	Conversions int       `json:"conversions"`
}

type CampaignMetrics struct {
	CampaignID       string  `json:"campaign_id"`
	TotalImpressions int64   `json:"total_impressions"`
	TotalClicks      int64   `json:"total_clicks"`
	TotalSpend       float64 `json:"total_spend"`
	TotalConversions int64   `json:"total_conversions"`
	CTR              float64 `json:"ctr"`
	CPA              *float64 `json:"cpa,omitempty"`
}

type AnalysisResult struct {
	TopCTR []CampaignMetrics `json:"top_ctr"`
	TopCPA []CampaignMetrics `json:"top_cpa"`
}
