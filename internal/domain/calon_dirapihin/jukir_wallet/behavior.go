package jukir_wallet

// import "errors"

// func (w *JukirWallet) Freeze(amount int64) error {
// 	if w.AvailableAmount < amount {
// 		return errors.New("saldo tidak cukup untuk dibekukan")
// 	}
// 	w.AvailableAmount -= amount
// 	w.FrozenAmount += amount
// 	return nil
// }

// func (w *JukirWallet) Unfreeze(amount int64) error {
// 	if w.FrozenAmount < amount {
// 		return errors.New("jumlah frozen tidak cukup")
// 	}
// 	w.FrozenAmount -= amount
// 	w.AvailableAmount += amount
// 	return nil
// }

// func (w *JukirWallet) Credit(amount int64) {
// 	w.BalanceAmount += amount
// 	w.AvailableAmount += amount
// }

// func (w *JukirWallet) Debit(amount int64) error {
// 	if w.AvailableAmount < amount {
// 		return errors.New("saldo tidak mencukupi")
// 	}
// 	w.BalanceAmount -= amount
// 	w.AvailableAmount -= amount
// 	return nil
// }
