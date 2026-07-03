package ledger

type JournalService interface {
	Record(
		transactionID string,
		entries []Entry,
	) error
}
