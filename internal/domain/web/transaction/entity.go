package transaction

type PaymentBreakdownItem struct {
	Label  string `json:"label"`
	Amount string `json:"amount"`
	Share  string `json:"share"`
	Tone   string `json:"tone"`
}

type TransactionIssueItem struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Action string `json:"action"`
	Tone   string `json:"tone"`
}

type DisputeCaseSummary struct {
	ID                  int64  `json:"id"`
	DisputeCaseCode     string `json:"disputeCaseCode"`
	ReferenceEntityType string `json:"referenceEntityType"`
	ReferenceEntityID   int64  `json:"referenceEntityId"`
	CaseType            string `json:"caseType"`
	Status              string `json:"status"`
	OpenedByUserID      int64  `json:"openedByUserId"`
	AssignedToUserID    int64  `json:"assignedToUserId"`
	OpenedAt            string `json:"openedAt"`
	ResolvedAt          string `json:"resolvedAt,omitempty"`
	ResolutionNote      string `json:"resolutionNote,omitempty"`
}
