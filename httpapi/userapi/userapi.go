package userapi

import (
	"github.com/centrifuge/go-centrifuge/documents"
	"github.com/centrifuge/go-centrifuge/extensions/transferdetails"
	"github.com/centrifuge/go-centrifuge/nft"
	"github.com/go-chi/chi"
)

const (
	documentIDParam = "document_id"
	transferIDParam = "transfer_id"
)

// Register registers the core apis to the router.
func Register(r chi.Router,
	docSrv documents.Service,
	nftSrv nft.Service,
	transferSrv transferdetails.Service) {
	h := handler{
		tokenRegistry: nftSrv.(documents.TokenRegistry),
		srv:           Service{docSrv: docSrv, transferDetailsService: transferSrv},
	}
	r.Post("/documents/{"+documentIDParam+"}/transfer_details", h.CreateTransferDetail)
	r.Put("/documents/{"+documentIDParam+"}/transfer_details/{"+transferIDParam+"}", h.UpdateTransferDetail)
	r.Get("/documents/{"+documentIDParam+"}/transfer_details", h.GetTransferDetailList)
	r.Get("/documents/{"+documentIDParam+"}/transfer_details/{"+transferIDParam+"}", h.GetTransferDetail)
}
