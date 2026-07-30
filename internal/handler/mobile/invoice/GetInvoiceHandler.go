package invoice

import (
	"net/http"

	"modulegue/core/response"
	usecase "modulegue/internal/domain/mobile/usecase/invoice"
)

type GetInvoiceHandler struct {
	getInvoiceUc *usecase.GetInvoiceUseCase
}

func NewGetInvoiceHandler(getInvoiceUc *usecase.GetInvoiceUseCase) *GetInvoiceHandler {
	return &GetInvoiceHandler{getInvoiceUc: getInvoiceUc}
}

func (h *GetInvoiceHandler) Execute(w http.ResponseWriter, r *http.Request) {
	invoiceNumber := r.PathValue("invoice_number")
	if invoiceNumber == "" {
		response.Error(w, http.StatusBadRequest, "invoice code wajib diisi")
		return
	}

	result, err := h.getInvoiceUc.Execute(r.Context(), invoiceNumber)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil invoice")
		return
	}

	response.Success(w, http.StatusOK, "ambil invoice berhasil", result)
}
