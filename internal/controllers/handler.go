package controllers

type handler struct {
	consumer EventWriter
}

func New(consumer EventWriter) handler {
	return handler{
		consumer: consumer,
	}
}
