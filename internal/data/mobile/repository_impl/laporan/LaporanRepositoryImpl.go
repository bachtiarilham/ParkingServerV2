package laporan

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	model "modulegue/internal/domain/mobile/model/laporan"
	"modulegue/internal/domain/mobile/repository"
)

type LaporanRepositoryImpl struct {
	db *sql.DB
}

func NewLaporanRepositoryImpl(db *sql.DB) repository.LaporanRepository {
	return &LaporanRepositoryImpl{db: db}
}

func (r *LaporanRepositoryImpl) GetLaporan(ctx context.Context, filter model.LaporanRequestModel) (*model.LaporanModel, error) {
	startDate := strings.TrimSpace(filter.StartDate)
	endDate := strings.TrimSpace(filter.EndDate)
	lokasi := strings.TrimSpace(filter.Lokasi)

	if startDate == "" {
		startDate = currentDateString()
	}
	if endDate == "" {
		endDate = startDate
	}

	summary, err := r.getSummary(ctx, filter.UserID, filter.RoleID, startDate, endDate, lokasi)
	if err != nil {
		return nil, err
	}

	chartBars, err := r.getChartBars(ctx, filter.UserID, filter.RoleID, startDate, endDate, lokasi)
	if err != nil {
		return nil, err
	}

	paymentSummaries, err := r.getPaymentSummaries(ctx, filter.UserID, filter.RoleID, startDate, endDate, lokasi)
	if err != nil {
		return nil, err
	}

	recentTransactions, err := r.getRecentTransactions(ctx, filter.UserID, filter.RoleID, startDate, endDate, lokasi)
	if err != nil {
		return nil, err
	}

	tanggalTerpilih := endDate
	label := fmt.Sprintf("%s s/d %s", startDate, endDate)

	return &model.LaporanModel{
		TanggalTerpilih: &tanggalTerpilih,
		Periode: &model.LaporanDateRangeModel{
			StartDate: &startDate,
			EndDate:   &endDate,
			Label:     &label,
		},
		Summary:            summary,
		ChartBars:          chartBars,
		PaymentSummaries:   paymentSummaries,
		RecentTransactions: recentTransactions,
	}, nil
}

func (r *LaporanRepositoryImpl) getSummary(ctx context.Context, userID, roleID int64, startDate, endDate, lokasi string) (*model.LaporanSummaryModel, error) {
	query := `
		SELECT
			COUNT(*) AS total_transaksi,
			COALESCE(SUM(fpt.final_amount), 0) AS total_pendapatan
		FROM financial_parking_transaction fpt
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		WHERE DATE(fpt.occurred_at) BETWEEN ? AND ?
		  AND (? = '' OR pl.location_name = ?)
		  AND (? = 0 OR fpt.customer_user_id = ? OR fpt.jukir_user_id = ? OR fpt.officer_user_id = ?)
		  AND fpt.transaction_status <> 'void'
	`

	var totalTransaksi int
	var totalPendapatan int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		startDate, endDate,
		lokasi, lokasi,
		userID, userID, userID, userID,
	).Scan(&totalTransaksi, &totalPendapatan)
	if err != nil {
		return nil, fmt.Errorf("get laporan summary: %w", err)
	}

	var rataRata int64
	if totalTransaksi > 0 {
		rataRata = totalPendapatan / int64(totalTransaksi)
	}

	return &model.LaporanSummaryModel{
		TotalTransaksi:    &totalTransaksi,
		TotalPendapatan:   &totalPendapatan,
		RataRataTransaksi: &rataRata,
	}, nil
}

func (r *LaporanRepositoryImpl) getChartBars(ctx context.Context, userID, roleID int64, startDate, endDate, lokasi string) ([]model.LaporanChartBarModel, error) {
	query := `
		SELECT
			DATE_FORMAT(fpt.occurred_at, '%Y-%m-%d') AS tanggal,
			COALESCE(SUM(fpt.final_amount), 0) AS amount
		FROM financial_parking_transaction fpt
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		WHERE DATE(fpt.occurred_at) BETWEEN ? AND ?
		  AND (? = '' OR pl.location_name = ?)
		  AND (? = 0 OR fpt.customer_user_id = ? OR fpt.jukir_user_id = ? OR fpt.officer_user_id = ?)
		  AND fpt.transaction_status <> 'void'
		GROUP BY DATE(fpt.occurred_at)
		ORDER BY DATE(fpt.occurred_at) ASC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		startDate, endDate,
		lokasi, lokasi,
		userID, userID, userID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get laporan chart bars: %w", err)
	}
	defer rows.Close()

	result := []model.LaporanChartBarModel{}
	var maxAmount int64
	type rawRow struct {
		tanggal string
		amount  int64
	}
	rawRows := []rawRow{}

	for rows.Next() {
		var item rawRow
		if err := rows.Scan(&item.tanggal, &item.amount); err != nil {
			return nil, fmt.Errorf("scan laporan chart bar: %w", err)
		}
		if item.amount > maxAmount {
			maxAmount = item.amount
		}
		rawRows = append(rawRows, item)
	}

	for _, item := range rawRows {
		tanggal := item.tanggal
		amount := item.amount
		value := 0.0
		if maxAmount > 0 {
			value = float64(amount) / float64(maxAmount)
		}
		periodLabel := tanggal
		periodStartDate := startDate
		periodEndDate := endDate

		result = append(result, model.LaporanChartBarModel{
			Tanggal:         &tanggal,
			Amount:          &amount,
			Value:           &value,
			PeriodLabel:     &periodLabel,
			PeriodStartDate: &periodStartDate,
			PeriodEndDate:   &periodEndDate,
		})
	}

	return result, nil
}

func (r *LaporanRepositoryImpl) getPaymentSummaries(ctx context.Context, userID, roleID int64, startDate, endDate, lokasi string) ([]model.LaporanPaymentSummaryModel, error) {
	totalAmountQuery := `
		SELECT COALESCE(SUM(fpt.final_amount), 0)
		FROM financial_parking_transaction fpt
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		WHERE DATE(fpt.occurred_at) BETWEEN ? AND ?
		  AND (? = '' OR pl.location_name = ?)
		  AND (? = 0 OR fpt.customer_user_id = ? OR fpt.jukir_user_id = ? OR fpt.officer_user_id = ?)
		  AND fpt.transaction_status <> 'void'
	`

	var totalAmount int64
	if err := r.db.QueryRowContext(
		ctx,
		totalAmountQuery,
		startDate, endDate,
		lokasi, lokasi,
		userID, userID, userID, userID,
	).Scan(&totalAmount); err != nil {
		return nil, fmt.Errorf("get laporan payment total: %w", err)
	}

	query := `
		SELECT
			COALESCE(pm.payment_method_name, COALESCE(fpt.payment_method, 'LAINNYA')) AS label,
			COALESCE(SUM(fpt.final_amount), 0) AS amount
		FROM financial_parking_transaction fpt
		LEFT JOIN payment_method pm ON pm.id = fpt.payment_method_id
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		WHERE DATE(fpt.occurred_at) BETWEEN ? AND ?
		  AND (? = '' OR pl.location_name = ?)
		  AND (? = 0 OR fpt.customer_user_id = ? OR fpt.jukir_user_id = ? OR fpt.officer_user_id = ?)
		  AND fpt.transaction_status <> 'void'
		GROUP BY COALESCE(pm.payment_method_name, COALESCE(fpt.payment_method, 'LAINNYA'))
		ORDER BY amount DESC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		startDate, endDate,
		lokasi, lokasi,
		userID, userID, userID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get laporan payment summaries: %w", err)
	}
	defer rows.Close()

	result := []model.LaporanPaymentSummaryModel{}
	for rows.Next() {
		var label string
		var amount int64
		if err := rows.Scan(&label, &amount); err != nil {
			return nil, fmt.Errorf("scan laporan payment summary: %w", err)
		}

		percentage := 0
		if totalAmount > 0 {
			percentage = int((amount * 100) / totalAmount)
		}

		result = append(result, model.LaporanPaymentSummaryModel{
			Label:      &label,
			Amount:     &amount,
			Percentage: &percentage,
		})
	}

	return result, nil
}

func (r *LaporanRepositoryImpl) getRecentTransactions(ctx context.Context, userID, roleID int64, startDate, endDate, lokasi string) ([]model.LaporanRecentTransactionModel, error) {
	query := `
		SELECT
			fpt.transaction_code,
			DATE_FORMAT(COALESCE(fpt.paid_at, fpt.occurred_at), '%Y-%m-%d %H:%i:%s') AS waktu,
			fpt.final_amount,
			COALESCE(pm.payment_method_name, COALESCE(fpt.payment_method, 'LAINNYA')) AS payment_tag
		FROM financial_parking_transaction fpt
		LEFT JOIN payment_method pm ON pm.id = fpt.payment_method_id
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		WHERE DATE(fpt.occurred_at) BETWEEN ? AND ?
		  AND (? = '' OR pl.location_name = ?)
		  AND (? = 0 OR fpt.customer_user_id = ? OR fpt.jukir_user_id = ? OR fpt.officer_user_id = ?)
		  AND fpt.transaction_status <> 'void'
		ORDER BY COALESCE(fpt.paid_at, fpt.occurred_at) DESC
		LIMIT 10
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		startDate, endDate,
		lokasi, lokasi,
		userID, userID, userID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get laporan recent transactions: %w", err)
	}
	defer rows.Close()

	result := []model.LaporanRecentTransactionModel{}
	for rows.Next() {
		var code string
		var timeText string
		var total int64
		var paymentTag string
		if err := rows.Scan(&code, &timeText, &total, &paymentTag); err != nil {
			return nil, fmt.Errorf("scan laporan recent transaction: %w", err)
		}

		result = append(result, model.LaporanRecentTransactionModel{
			Code:       &code,
			Time:       &timeText,
			Total:      &total,
			PaymentTag: &paymentTag,
		})
	}

	return result, nil
}

func currentDateString() string {
	return time.Now().Format("2006-01-02")
}
