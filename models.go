package twitter_ads

import "time"

type ActiveEntitiesResponse struct {
	Request struct {
		Params struct {
			AccountID string    `json:"account_id"`
			Entity    string    `json:"entity"`
			StartTime time.Time `json:"start_time"`
			EndTime   time.Time `json:"end_time"`
		} `json:"params"`
	} `json:"request"`
	Data []struct {
		EntityID          string    `json:"entity_id"`
		ActivityStartTime time.Time `json:"activity_start_time"`
		ActivityEndTime   time.Time `json:"activity_end_time"`
		Placements        []string  `json:"placements"`
	} `json:"data"`
	ErrorResponse *ErrorResponse
}

type ErrorResponse struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Request struct {
		Params struct {
		} `json:"params"`
	} `json:"request"`
}

type SpendStats struct {
	DataType string `json:"data_type"`
	Data     []struct {
		ID     string `json:"id"`
		IDData []struct {
			Metrics struct {
				BilledChargeLocalMicro []int `json:"billed_charge_local_micro"`
			} `json:"metrics"`
		} `json:"id_data"`
	} `json:"data"`
	Request struct {
		Params struct {
			Country     interface{} `json:"country"`
			Placement   string      `json:"placement"`
			Granularity string      `json:"granularity"`
			Platform    interface{} `json:"platform"`
		} `json:"params"`
	} `json:"request"`
}

type CampaignInfo struct {
	Data struct {
		Name                        string      `json:"name"`
		StartTime                   time.Time   `json:"start_time"`
		Servable                    bool        `json:"servable"`
		EffectiveStatus             string      `json:"effective_status"`
		DailyBudgetAmountLocalMicro int         `json:"daily_budget_amount_local_micro"`
		EndTime                     interface{} `json:"end_time"`
		FundingInstrumentId         string      `json:"funding_instrument_id"`
		Id                          string      `json:"id"`
		EntityStatus                string      `json:"entity_status"`
		Currency                    string      `json:"currency"`
		CreatedAt                   time.Time   `json:"created_at"`
		UpdatedAt                   time.Time   `json:"updated_at"`
	} `json:"data"`
}
