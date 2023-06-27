package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	_ EventHandler = (*Router)(nil)
)

type RouteEventResult struct {
	Message            tgbotapi.Chattable
	NextMessageHandler RoutedEventHandler
}

type RoutedEvent struct {
	ChatID int64

	*tgbotapi.Update
}

// RoutedEventHandler handles messages from Telegram bot user
type RoutedEventHandler interface {
	HandleBotEvent(ctx context.Context, e RoutedEvent) (*RouteEventResult, error)
}

type Handlers struct {
	Commands map[string]RoutedEventHandler
	Default  RoutedEventHandler
}

type PendingEvent struct {
	PreviousEvent *tgbotapi.Update
	NextHandler   RoutedEventHandler
}

type Router struct {
	handlers      Handlers
	pendingEvents Map[int64, *PendingEvent]
}

func (r Router) HandleBotEvent(ctx context.Context, e *tgbotapi.Update) (tgbotapi.Chattable, error) {
	event := RoutedEvent{
		Update: e,
	}

	if chat := e.FromChat(); chat != nil {
		event.ChatID = chat.ID
	}

	handler, err := r.getHandler(e)
	if err != nil {
		return nil, err
	}

	result, err := handler.HandleBotEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	if result.NextMessageHandler != nil {
		r.setPendingEvent(e, &PendingEvent{
			PreviousEvent: e,
			NextHandler:   result.NextMessageHandler,
		})
	} else {
		r.removePendingEvents(e)
	}

	return result.Message, nil
}

func (r Router) getHandler(e *tgbotapi.Update) (RoutedEventHandler, error) {
	cmd, ok := commandFromMessage(e)
	if !ok {
		pendingEvent := r.getPendingEvent(e)
		if pendingEvent != nil {
			return pendingEvent.NextHandler, nil
		}

		return r.handlers.Default, nil
	}

	if cmd == "" {
		return r.handlers.Default, nil
	}

	handler, ok := r.handlers.Commands[cmd]
	if ok {
		return handler, nil
	}

	return r.handlers.Default, nil
}

func (r Router) getPendingEvent(e *tgbotapi.Update) *PendingEvent {
	chat := e.FromChat()
	if chat == nil {
		return nil
	}

	event, ok := r.pendingEvents.Get(chat.ID)
	if !ok {
		return nil
	}

	return event
}

func (r Router) setPendingEvent(e *tgbotapi.Update, nextEvent *PendingEvent) {
	chat := e.FromChat()
	if chat == nil {
		return
	}

	r.pendingEvents.Set(chat.ID, nextEvent)
}

func (r Router) removePendingEvents(e *tgbotapi.Update) {
	chat := e.FromChat()
	if chat == nil {
		return
	}

	r.pendingEvents.Delete(chat.ID)
}

func NewRouter(handlers Handlers) *Router {
	return &Router{
		handlers:      handlers,
		pendingEvents: NewMap[int64, *PendingEvent](),
	}
}

func commandFromMessage(u *tgbotapi.Update) (string, bool) {
	if u.Message == nil {
		return "", false
	}

	msg := strings.TrimSpace(u.Message.Text)
	if msg == "" {
		return "", false
	}

	return msg[1:], msg[0] == '/'
}
