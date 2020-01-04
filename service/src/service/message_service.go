package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/callstats-io/ai-decision/service/gen/protos"
	"github.com/callstats-io/ai-decision/service/src/flowdock"
	"github.com/callstats-io/ai-decision/service/src/grpc"
	"github.com/callstats-io/ai-decision/service/src/message"
	"github.com/callstats-io/ai-decision/service/src/storage"
	"github.com/callstats-io/go-common/log"
	"github.com/golang/protobuf/ptypes"
)

// MessageStorage defines the interface the service expects of any message storage backend
type MessageStorage interface {
	FetchMessageTemplates(ctx context.Context, messageType string, maxVersion int32) ([]*storage.MessageTemplate, error)
	CreateMessage(ctx context.Context, msg *storage.Message) error
	ListMessages(ctx context.Context, appID int32, messageType string, minVersion, maxVersion int32, from, to *time.Time) ([]*storage.Message, error)
}

// AIDecisionMessageService implements the protos AIDecisionMessageServiceServer
type AIDecisionMessageService struct {
	messageStorage MessageStorage
	flowdockClient *flowdock.Client
}

var _ = protos.AIDecisionMessageServiceServer(&AIDecisionMessageService{})

//NewAIDecisionMessageService returns a new AIDecisionMessageService or an error if initialization fails
func NewAIDecisionMessageService(ms MessageStorage, flowdockClient *flowdock.Client) (*AIDecisionMessageService, error) {
	s := &AIDecisionMessageService{
		messageStorage: ms,
		flowdockClient: flowdockClient,
	}
	return s, nil
}

// Create stores a new message based on a pre-existing template.
func (s *AIDecisionMessageService) Create(ctx context.Context, req *protos.MessageCreateRequest) (*protos.Message, error) {
	genTime, _ := ptypes.Timestamp(req.GenerationTime)
	logger := log.FromContext(ctx).With(
		log.Int(LogKeyAppID, int(req.AppId)),
		log.String(LogKeyTemplateType, req.Type),
		log.Int(LogKeyTemplateVersion, int(req.Version)),
		log.Time(LogKeyGenerationTime, genTime),
	)
	ctx = log.WithLogger(ctx, logger)
	if err := s.validateCreateRequest(ctx, req); err != nil {
		return nil, err
	}

	templateData, err := message.UnmarshalTemplateData(req.Data)
	if err != nil {
		return nil, grpc.ErrInvalidArgument(ctx, fmt.Errorf("data: %s", err))
	}

	templates, err := s.messageStorage.FetchMessageTemplates(ctx, req.Type, req.Version)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, grpc.ErrNotFound(ctx, err)
		}
		return nil, grpc.ErrUnavailable(ctx, err)
	}

	var renderedMsg string
	var template *storage.MessageTemplate
	for _, t := range templates {
		mt, err := message.NewTemplate(t)
		if err != nil {
			return nil, grpc.ErrFailedPrecondition(ctx, err)
		}
		if m, err := mt.RenderString(templateData); err != nil {
			return nil, grpc.ErrInvalidArgument(ctx, err)
		} else if mt.Version() == req.Version {
			template = t
			renderedMsg = m // keep the rendered message for return value
		}
	}

	// a template must always exist for us to end up here (otherwise db should return a not found error)
	// so we ignore nil-validations. Furthermore, validations should account for data validity so timestamp error is ignored.
	msg := &storage.Message{
		AppID:       req.AppId,
		TemplateID:  template.ID,
		Template:    template,
		GeneratedAt: genTime,
		Data:        req.Data,
	}

	if err := s.messageStorage.CreateMessage(ctx, msg); err != nil {
		if err == storage.ErrNotFound {
			return nil, grpc.ErrNotFound(ctx, err)
		}
		if strings.Contains(err.Error(), "violates unique constraint") {
			return nil, grpc.ErrFailedPrecondition(ctx, err)
		}
		return nil, grpc.ErrUnavailable(ctx, err)
	}

	if err := s.flowdockClient.SendAiNotificationMessage(req.AppId, req.Type, renderedMsg); err != nil {
		logger.Warn("Error in Flowdock Send: ", log.Error(err))
	}

	return &protos.Message{
		AppId:          req.AppId,
		Type:           req.Type,
		Version:        req.Version,
		GenerationTime: req.GenerationTime,
		Data:           req.Data,
		Message:        renderedMsg,
	}, nil
}

func (s *AIDecisionMessageService) validateCreateRequest(ctx context.Context, req *protos.MessageCreateRequest) error {
	return validate(ctx,
		validatePositiveInt("app_id", req.AppId),
		validateNonEmptyString("type", req.Type),
		validatePositiveInt("version", req.Version),
		validateNonEmptyBytes("data", req.Data),
		validateTimestamp("generation_time", req.GenerationTime),
	)
}

// List renders a list of messages as a stream based on the provided filter criteria.
func (s *AIDecisionMessageService) List(req *protos.MessageListRequest, stream protos.AIDecisionMessageService_ListServer) error {
	ctx := stream.Context()
	logger := log.FromContext(ctx).With(
		log.Int(LogKeyAppID, int(req.AppId)),
		log.String(LogKeyTemplateType, req.Type),
		log.Int(LogKeyTemplateMinVersion, int(req.MinVersion)),
		log.Int(LogKeyTemplateMaxVersion, int(req.MaxVersion)),
	)
	var generatedAtFrom, generatedAtTo *time.Time
	if req.GenerationTimeFrom != nil {
		v, _ := ptypes.Timestamp(req.GenerationTimeFrom)
		generatedAtFrom = &v
		logger = logger.With(log.Time(LogKeyGenerationTimeFrom, v))
	}
	if req.GenerationTimeTo != nil {
		v, _ := ptypes.Timestamp(req.GenerationTimeTo)
		logger = logger.With(log.Time(LogKeyGenerationTimeTo, v))
		generatedAtTo = &v
	}
	ctx = log.WithLogger(ctx, logger)

	if err := s.validateListRequest(ctx, req); err != nil {
		return err
	}

	messages, err := s.messageStorage.ListMessages(ctx, req.AppId, req.Type, req.MinVersion, req.MaxVersion, generatedAtFrom, generatedAtTo)
	if err == storage.ErrNotFound {
		return grpc.ErrNotFound(ctx, err)
	} else if err != nil {
		return grpc.ErrUnavailable(ctx, err)
	}

	for _, msg := range messages {
		// render message
		mt, err := message.NewTemplate(msg.Template)
		if err != nil {
			// should never happen, likely an invalid template in db WITH a message that refers to it
			// which would mean someone has gone and done something stupid manually
			return grpc.ErrFailedPrecondition(ctx, err)
		}
		tmplData, err := message.UnmarshalTemplateData(msg.Data)
		if err != nil {

		}
		rendered, err := mt.RenderString(tmplData)
		if err != nil {
			return grpc.ErrInvalidArgument(ctx, err)
		}

		// send to requester
		genTime, _ := ptypes.TimestampProto(msg.GeneratedAt)
		if err := stream.Send(&protos.Message{
			AppId:          msg.AppID,
			Type:           msg.Template.Type,
			Version:        msg.Template.Version,
			Data:           msg.Data,
			GenerationTime: genTime,
			Message:        rendered,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *AIDecisionMessageService) validateListRequest(ctx context.Context, req *protos.MessageListRequest) error {
	return validate(ctx, validatePositiveInt("app_id", req.AppId))
}
