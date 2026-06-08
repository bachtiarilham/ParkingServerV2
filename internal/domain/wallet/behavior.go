package wallet

import "errors"

func (w *Wallet) Freeze(amount int64) error {
	if w.AvailableAmount < amount {
		return errors.New("saldo tidak cukup untuk dibekukan")
	}
	w.AvailableAmount -= amount
	w.FrozenAmount += amount
	return nil
}

func (w *Wallet) Unfreeze(amount int64) error {
	if w.FrozenAmount < amount {
		return errors.New("jumlah frozen tidak cukup")
	}
	w.FrozenAmount -= amount
	w.AvailableAmount += amount
	return nil
}

func (w *Wallet) Credit(amount int64) {
	w.BalanceAmount += amount
	w.AvailableAmount += amount
}

func (w *Wallet) Debit(amount int64) error {
	if w.AvailableAmount < amount {
		return errors.New("saldo tidak mencukupi")
	}
	w.BalanceAmount -= amount
	w.AvailableAmount -= amount
	return nil
}
