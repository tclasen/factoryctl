package web

import "github.com/tclasen/factoryctl/internal/delivery"

func (h *Handler) documentDelivery(d Document) any {
	for _, p := range h.projects {
		if documentID(p.Checkout, d.Path) == d.ID {
			return delivery.Resolve(p.Root, d.Path)
		}
	}
	return nil
}
