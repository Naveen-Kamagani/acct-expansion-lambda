package slog

import (
	"log/slog"
	"os"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	slogformatter "github.com/samber/slog-formatter"
)

var programLevel = new(slog.LevelVar)

var AnyFormatter slogformatter.Formatter
var ErrorFormatter slogformatter.Formatter
var TimeFormatter slogformatter.Formatter

var HTTPRequestFormatter slogformatter.Formatter
var HTTPResponseFormatter slogformatter.Formatter

var UpdateItemInputFormatter slogformatter.Formatter
var UpdateItemOutputFormatter slogformatter.Formatter
var QueryInputFormatter slogformatter.Formatter
var QueryOutputFormatter slogformatter.Formatter
var PutItemInputFormatter slogformatter.Formatter
var PutItemOutputFormatter slogformatter.Formatter
var GetItemInputFormatter slogformatter.Formatter
var GetItemOutputFormatter slogformatter.Formatter
var TransactWriteItemsInputFormatter slogformatter.Formatter
var TransactWriteItemsOutputFormatter slogformatter.Formatter
var DeleteItemInputFormatter slogformatter.Formatter
var DeleteItemOutputFormatter slogformatter.Formatter
var PutFormatter slogformatter.Formatter

var GetQueueURLInputFormatter slogformatter.Formatter
var GetQueueURLOutputFormatter slogformatter.Formatter
var GetQueueAttributesInputFormatter slogformatter.Formatter
var GetQueueAttributesOutputFormatter slogformatter.Formatter
var SendMessageInputFormatter slogformatter.Formatter
var SendMessageOutputFormatter slogformatter.Formatter
var SendMessageBatchInputFormatter slogformatter.Formatter
var SendMessageBatchOutputFormatter slogformatter.Formatter

var JSONLogger *slog.Logger
var TextLogger *slog.Logger

func getGroupValueForObj(objType string, obj interface{}) slog.Value {
	return slog.GroupValue(
		slog.String("Object Type", objType),
		slog.Any(reflect.TypeOf(obj).String(), obj),
	)
}

func initializeFormatters() {
	// Common formatters
	AnyFormatter = slogformatter.FormatByType(func(obj interface{}) slog.Value {
		return getGroupValueForObj(reflect.TypeOf(obj).String(), obj)
	})
	ErrorFormatter = slogformatter.ErrorFormatter("error")
	TimeFormatter = slogformatter.TimeFormatter(time.DateTime, time.UTC)

	// net/http request and response formatters
	HTTPRequestFormatter = slogformatter.HTTPRequestFormatter(false)
	HTTPResponseFormatter = slogformatter.HTTPResponseFormatter(false)

	UpdateItemInputFormatter = slogformatter.FormatByType(func(obj dynamodb.UpdateItemInput) slog.Value {
		return getGroupValueForObj("Request to update item to dynamodb", obj)
	})
	UpdateItemOutputFormatter = slogformatter.FormatByType(func(obj dynamodb.UpdateItemOutput) slog.Value {
		return getGroupValueForObj("Response from dynamodb to update item", obj)
	})
	QueryInputFormatter = slogformatter.FormatByType(func(obj dynamodb.QueryInput) slog.Value {
		return getGroupValueForObj("Request to query from dynamodb", obj)
	})
	QueryOutputFormatter = slogformatter.FormatByType(func(obj dynamodb.QueryOutput) slog.Value {
		return getGroupValueForObj("Response from dynamodb for query", obj)
	})
	PutItemInputFormatter = slogformatter.FormatByType(func(obj dynamodb.PutItemInput) slog.Value {
		return getGroupValueForObj("Request item to put to dynamodb", obj)
	})
	PutItemOutputFormatter = slogformatter.FormatByType(func(obj dynamodb.PutItemOutput) slog.Value {
		return getGroupValueForObj("Response from dynamodb for request to put item", obj)
	})
	GetItemInputFormatter = slogformatter.FormatByType(func(obj dynamodb.GetItemInput) slog.Value {
		return getGroupValueForObj("Request to retrieve single item from dynamodb", obj)
	})
	GetItemOutputFormatter = slogformatter.FormatByType(func(obj dynamodb.GetItemOutput) slog.Value {
		return getGroupValueForObj("Response from dynamodb for single item request", obj)
	})
	TransactWriteItemsInputFormatter = slogformatter.FormatByType(func(obj dynamodb.TransactWriteItemsInput) slog.Value {
		return getGroupValueForObj("Request items to insert to dynamodb", obj)
	})
	TransactWriteItemsOutputFormatter = slogformatter.FormatByType(func(obj dynamodb.TransactWriteItemsOutput) slog.Value {
		return getGroupValueForObj("Response from dynamodb to insert items", obj)
	})
	DeleteItemInputFormatter = slogformatter.FormatByType(func(obj dynamodb.DeleteItemInput) slog.Value {
		return getGroupValueForObj("Request to delete item from dynamodb", obj)
	})
	DeleteItemOutputFormatter = slogformatter.FormatByType(func(obj dynamodb.DeleteItemOutput) slog.Value {
		return getGroupValueForObj("Response from dynamodb to delete item", obj)
	})
	PutFormatter = slogformatter.FormatByType(func(obj types.Put) slog.Value {
		return getGroupValueForObj("Request items to put to dynamodb", obj)
	})

	// sqs formatters
	GetQueueURLInputFormatter = slogformatter.FormatByType(func(obj sqs.GetQueueUrlInput) slog.Value {
		return getGroupValueForObj("Request to retrieve sqs queue url", obj)
	})
	GetQueueURLOutputFormatter = slogformatter.FormatByType(func(obj sqs.GetQueueUrlOutput) slog.Value {
		return getGroupValueForObj("Response to query for queue url", obj)
	})
	GetQueueAttributesInputFormatter = slogformatter.FormatByType(func(obj sqs.GetQueueAttributesInput) slog.Value {
		return getGroupValueForObj("Request to retrieve sqs queue attributes", obj)
	})
	GetQueueAttributesOutputFormatter = slogformatter.FormatByType(func(obj sqs.GetQueueAttributesOutput) slog.Value {
		return getGroupValueForObj("Response to retrieve sqs queue attributes", obj)
	})
	SendMessageInputFormatter = slogformatter.FormatByType(func(obj sqs.SendMessageInput) slog.Value {
		return getGroupValueForObj("Request message to send to sqs queue", obj)
	})
	SendMessageOutputFormatter = slogformatter.FormatByType(func(obj sqs.SendMessageOutput) slog.Value {
		return getGroupValueForObj("Response to sending message to sqs queue", obj)
	})
	SendMessageBatchInputFormatter = slogformatter.FormatByType(func(obj sqs.SendMessageBatchInput) slog.Value {
		return getGroupValueForObj("Request batch message to send to sqs queue", obj)
	})
	SendMessageBatchOutputFormatter = slogformatter.FormatByType(func(obj sqs.SendMessageBatchOutput) slog.Value {
		return getGroupValueForObj("Response to sending batch message to sqs queue", obj)
	})
}

func getLogger(sink slog.Handler) *slog.Logger {
	initializeFormatters()
	return slog.New(
		slogformatter.NewFormatterHandler(
			// common formatters
			ErrorFormatter, TimeFormatter,
			// http req/resp formatters
			HTTPRequestFormatter, HTTPResponseFormatter,
			// dynamodb formatters
			UpdateItemInputFormatter, UpdateItemOutputFormatter, QueryInputFormatter, QueryOutputFormatter,
			PutItemInputFormatter, PutItemOutputFormatter, GetItemInputFormatter, GetItemOutputFormatter,
			TransactWriteItemsInputFormatter, TransactWriteItemsOutputFormatter, DeleteItemInputFormatter, DeleteItemOutputFormatter, PutFormatter,
			// sqs formatters
			GetQueueURLInputFormatter, GetQueueURLOutputFormatter, GetQueueAttributesInputFormatter, GetQueueAttributesOutputFormatter,
			SendMessageInputFormatter, SendMessageOutputFormatter, SendMessageBatchInputFormatter, SendMessageBatchOutputFormatter,
			// Any other formatter
			AnyFormatter,
		)(sink),
	)
}

func InitializeLoggers() {
	jsonSink := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	textSink := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	JSONLogger = getLogger(jsonSink)
	TextLogger = getLogger(textSink)
}

func SetInfoLevel() {
	programLevel.Set(slog.LevelInfo)
}

func SetDebugLevel() {
	programLevel.Set(slog.LevelDebug)
}

func SetErrorLevel() {
	programLevel.Set(slog.LevelError)
}

func SetWarnLevel() {
	programLevel.Set(slog.LevelWarn)
}
