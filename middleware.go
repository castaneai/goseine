package goseine

import "github.com/sirupsen/logrus"

// adapter pattern
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
type PacketMiddleware func(PacketHandler) PacketHandler

func RequestPacketLogger(logger *logrus.Logger) PacketMiddleware {
	return func(h PacketHandler) PacketHandler {
		return PacketHandlerFunc(func(w *PacketWriter, req *Packet) {
			logger.Println(req)
			h.Handle(w, req)
		})
	}
}

func WithMiddleware(h PacketHandler, middlewares ...PacketMiddleware) PacketHandler {
	for _, mw := range middlewares {
		h = mw(h)
	}
	return h
}
