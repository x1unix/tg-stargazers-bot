package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	_ EventHandler = (*Router)(nil)
)

type ChatID = int64

type userEventType = int

const (
	userJoinEvent userEventType = iota
	userLeftEvent
)

type RouteEventResult struct {
	Message            tgbotapi.Chattable
	NextMessageHandler RoutedEventHandler
}

type RoutedEvent struct {
	ChatID ChatID

	*tgbotapi.Update
}

// ChatLifecycleHandler handles bot start and stop events.
type ChatLifecycleHandler interface {
	// HandleUserJoin triggered when user creates bot chat.
	HandleUserJoin(ctx context.Context, e RoutedEvent) (*RouteEventResult, error)

	// HandleUserLeave triggered when user blocks the bot and leaves the chat.
	HandleUserLeave(ctx context.Context, e RoutedEvent) error
}

// RoutedEventHandler handles messages from Telegram bot user
type RoutedEventHandler interface {
	HandleBotEvent(ctx context.Context, e RoutedEvent) (*RouteEventResult, error)
}

type CommandHandlers = map[string]CommandHandler

// CommandHandler implements bot command handler.
type CommandHandler interface {
	RoutedEventHandler

	// CommandDescription returns bot command description for help.
	CommandDescription() string
}

// Handlers contains Telegram bot events and commands handling config.
type Handlers struct {
	// Commands is key-value pair of command name and handler.
	Commands CommandHandlers

	// Help is help command handler.
	//
	// Use this to override default help message.
	Help RoutedEventHandler

	// Start is start command handler.
	//
	// This is a mandatory field.
	Start RoutedEventHandler

	// Default is default event handler.
	//
	// Handler is called when no matching command found.
	Default RoutedEventHandler

	// LifecycleHandler handles chat lifecycle.
	LifecycleHandler ChatLifecycleHandler
}

type PendingEvent struct {
	PreviousEvent *tgbotapi.Update
	NextHandler   RoutedEventHandler
}

type Router struct {
	handlers      Handlers
	pendingEvents Map[int64, *PendingEvent]
}

func NewRouter(handlers Handlers) *Router {
	return &Router{
		handlers:      applyDefaults(handlers),
		pendingEvents: NewMap[int64, *PendingEvent](),
	}
}

func (r Router) HandleBotEvent(ctx context.Context, e *tgbotapi.Update) (tgbotapi.Chattable, error) {
	event := RoutedEvent{
		Update: e,
	}

	if chat := e.FromChat(); chat != nil {
		event.ChatID = chat.ID
	}

	if chatEvent, ok := isChatLifecycleEvent(e); ok {
		switch chatEvent {
		case userJoinEvent:
			result, err := r.handlers.LifecycleHandler.HandleUserJoin(ctx, event)
			r.handleResult(result, e)
			return result.Message, err
		case userLeftEvent:
			r.removePendingEvents(e)
			return nil, r.handlers.LifecycleHandler.HandleUserLeave(ctx, event)
		}
	}

	handler, err := r.getHandler(e)
	if err != nil {
		return nil, err
	}

	result, err := handler.HandleBotEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	r.handleResult(result, e)
	return result.Message, nil
}

func (r Router) handleResult(result *RouteEventResult, u *tgbotapi.Update) {
	if result != nil && result.NextMessageHandler != nil {
		r.setPendingEvent(u, &PendingEvent{
			PreviousEvent: u,
			NextHandler:   result.NextMessageHandler,
		})
		return
	}

	r.removePendingEvents(u)
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

	switch cmd {
	case "":
		return r.handlers.Default, nil
	case "start":
		return r.handlers.Start, nil
	case "help":
		return r.handlers.Help, nil
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

func (r Router) getHelp(e *RoutedEvent) (*RouteEventResult, error) {
	sb := strings.Builder{}
	sb.WriteString("Here is a list of available commands:\n\n")
	for cmdName, cmd := range r.handlers.Commands {
		if cmdName == "start" {
			continue
		}

		sb.WriteString("* ")
		sb.WriteRune('/')
		sb.WriteString(cmdName)
		sb.WriteString(" - ")
		sb.WriteString(cmd.CommandDescription())
		sb.WriteRune('\n')
	}

	msg := tgbotapi.NewMessage(e.ChatID, sb.String())
	msg.ParseMode = ParseModeMarkdown
	return &RouteEventResult{
		Message: msg,
	}, nil
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

func isChatLifecycleEvent(u *tgbotapi.Update) (userEventType, bool) {
	if u.Message == nil {
		return -1, false
	}

	if u.Message.NewChatMembers != nil {
		return userJoinEvent, true
	}

	if u.Message.LeftChatMember != nil {
		return userLeftEvent, true
	}

	return -1, false
}
