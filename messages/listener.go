package messages

import (
	"errors"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/skyterra/clog"
	"github.com/skyterra/elastic-embedding-searcher/elastic"
	"github.com/skyterra/elastic-embedding-searcher/helper"
	"golang.org/x/net/context"
	"strings"
	"sync"
	"time"
)

const (
	DefaultCapOfCacheMessage = 1024
	SyncIntervalSec          = 5 // second

	moduleName = "listener"
)

type IMessage interface {
	GetValue() []byte
}

type IConsumer interface {
	FetchMessage(ctx context.Context) (IMessage, error)
	CommitMessages(ctx context.Context, msgs ...IMessage) error
	Close() error
}

type MessageParseFunc func([]byte) (elastic.IDocument, error)

type MessageListener struct {
	IndexName      string
	DimensionPart1 int32
	DimensionPart2 int32

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mutex    sync.Mutex
	consumer IConsumer
	messages []IMessage
	parse    MessageParseFunc
}

// Start begins listening to the message queue.
func (l *MessageListener) Start() error {
	// check to see if index exists. if NOT exist, return err
	exist, err := elastic.ExistIndex(l.ctx, l.IndexName)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("the index \"%s\" NOT exist", l.IndexName)
	}

	l.wg.Add(2)

	// worker for fetching message.
	go func() {
		defer func() {
			l.wg.Done()
			clog.Info(l.ctx, "exit fetch worker normally.")
		}()

		clog.Info(l.ctx, "fetch worker is ready.")

		for {
			select {
			case <-l.ctx.Done():
				return
			default:
				l.fetch()
			}
		}
	}()

	// worker for sync message to elastic.
	go func() {
		defer func() {
			l.wg.Done()
			clog.Info(l.ctx, "exit sync worker normally.")
		}()
		clog.Info(l.ctx, "sync worker is ready.")

		ticker := time.NewTicker(SyncIntervalSec * time.Second)
		for {
			select {
			case <-l.ctx.Done():
				return
			case <-ticker.C:
				l.sync()
			}
		}
	}()

	return nil
}

// Stop stops listen to kafka.
func (l *MessageListener) Stop() error {
	l.cancel()
	l.wg.Wait()
	return l.consumer.Close()
}

// fetch fetch message from message queue BUT not commit.
func (l *MessageListener) fetch() {
	defer func() {
		if err := recover(); err != nil {
			clog.Error(l.ctx, "panic:%v", err)
		}
	}()

	// read message from kafka.
	message, err := l.consumer.FetchMessage(l.ctx)
	if errors.Is(err, context.Canceled) {
		return
	}

	if err != nil || message == nil {
		return
	}

	// push message to queue.
	l.mutex.Lock()
	l.messages = append(l.messages, message)
	l.mutex.Unlock()
}

// sync parse message and sync message data to elasticsearch.
func (l *MessageListener) sync() {
	defer func() {
		if err := recover(); err != nil {
			clog.Error(l.ctx, "panic:%v", err)
		}
	}()

	// cut all message to local variable.
	var messages []IMessage
	if len(l.messages) > 0 {
		l.mutex.Lock()
		messages = l.messages
		l.messages = make([]IMessage, 0, DefaultCapOfCacheMessage)
		l.mutex.Unlock()
	}

	if len(messages) == 0 {
		return
	}

	// commit messages.
	err := l.consumer.CommitMessages(l.ctx, messages...)
	if err != nil {
		clog.Warn(l.ctx, "fail to commit messages in message listener. err:%s", err.Error())
	}

	handledMessage := make(map[uint64]struct{})
	documents := make([]elastic.IDocument, 0, len(messages))
	for i := len(messages) - 1; i >= 0; i-- {
		// remove duplicate message.
		sign := xxhash.Sum64(messages[i].GetValue())
		if _, exist := handledMessage[sign]; exist {
			continue
		}

		handledMessage[sign] = struct{}{}
		document, err := l.parse(messages[i].GetValue())
		if err != nil {
			clog.Warn(l.ctx, "fail to parse message to document. err:%s", err.Error())
			continue
		}

		documents = append(documents, document)
	}

	if len(documents) == 0 {
		return
	}

	err = elastic.IndexDocuments(l.ctx, l.IndexName, documents)
	if err != nil {
		clog.Warn(l.ctx, "fail to index document to elastic. err:%s", err.Error())
	}
}

// NewMessageListener create message listener.
func NewMessageListener(indexName string, consumer IConsumer, parse MessageParseFunc) (*MessageListener, error) {
	if strings.TrimSpace(indexName) == "" || consumer == nil || parse == nil {
		return nil, errors.New("invalid parameters")
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = helper.WithContextTrace(ctx)
	ctx = helper.WithContextModule(ctx, moduleName)

	listener := &MessageListener{
		IndexName: indexName,
		consumer:  consumer,
		parse:     parse,
		messages:  make([]IMessage, 0, DefaultCapOfCacheMessage),
		ctx:       ctx,
		cancel:    cancel,
	}

	return listener, nil
}
