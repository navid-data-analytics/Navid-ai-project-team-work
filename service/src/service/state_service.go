package service

import (
	"context"
	"time"

	"github.com/callstats-io/ai-decision/service/gen/protos"
	"github.com/callstats-io/ai-decision/service/src/grpc"
	"github.com/callstats-io/ai-decision/service/src/storage"
	"github.com/callstats-io/go-common/log"
	"github.com/golang/protobuf/ptypes"
)

// StateStorage defines the interface state service expects from applicable storages
type StateStorage interface {
	SaveState(ctx context.Context, state *storage.AidAnalyticsState) error
	GetState(ctx context.Context, state *storage.AidAnalyticsState) error
	ListStates(ctx context.Context, appID int32, keyword string, from, to *time.Time) ([]*storage.AidAnalyticsState, error)
}

// AIDecisionStateService implements the protos AIDecisionStateServiceServer
type AIDecisionStateService struct {
	stateStorage StateStorage
}

var _ = protos.AIDecisionStateServiceServer(&AIDecisionStateService{})

//NewAIDecisionStateService returns a new AIDecisionStateService or an error if initialization fails
func NewAIDecisionStateService(storage StateStorage) (*AIDecisionStateService, error) {
	s := &AIDecisionStateService{
		stateStorage: storage,
	}
	return s, nil
}

// Save stores AI decision analytics state
func (s *AIDecisionStateService) Save(ctx context.Context, req *protos.StateSaveRequest) (*protos.State, error) {
	savedAt, _ := ptypes.Timestamp(req.GenerationTime)
	ctx = log.WithLogger(ctx, log.FromContext(ctx).With(
		log.Int(LogKeyAppID, int(req.AppId)),
		log.String(LogKeyKeyword, req.Keyword),
		log.Time(LogKeyGenerationTime, savedAt),
	))
	if err := s.validateSaveRequest(ctx, req); err != nil {
		return nil, err
	}

	state := &storage.AidAnalyticsState{
		AppID:   req.AppId,
		Keyword: req.Keyword,
		Data:    req.Data,
		SavedAt: savedAt,
	}
	if err := s.stateStorage.SaveState(ctx, state); err != nil {
		return nil, grpc.ErrUnavailable(ctx, err)
	}

	// echo state back to caller
	savedAtProto, _ := ptypes.TimestampProto(state.SavedAt)
	return &protos.State{
		AppId:          state.AppID,
		Keyword:        state.Keyword,
		Data:           state.Data,
		GenerationTime: savedAtProto,
	}, nil
}

// Get retrieves AI decision analytics state
func (s *AIDecisionStateService) Get(ctx context.Context, req *protos.StateGetRequest) (*protos.State, error) {
	savedAt, _ := ptypes.Timestamp(req.GenerationTime)
	ctx = log.WithLogger(ctx, log.FromContext(ctx).With(
		log.Int(LogKeyAppID, int(req.AppId)),
		log.String(LogKeyKeyword, req.Keyword),
		log.Time(LogKeyGenerationTime, savedAt),
	))
	if err := s.validateGetRequest(ctx, req); err != nil {
		return nil, err
	}

	state := &storage.AidAnalyticsState{
		AppID:   req.AppId,
		Keyword: req.Keyword,
		SavedAt: savedAt,
	}
	if err := s.stateStorage.GetState(ctx, state); err == storage.ErrNotFound {
		return nil, grpc.ErrNotFound(ctx, err)
	} else if err != nil {
		return nil, grpc.ErrUnavailable(ctx, err)
	}

	// echo state back to caller
	savedAtProto, _ := ptypes.TimestampProto(state.SavedAt)
	return &protos.State{
		AppId:          state.AppID,
		Keyword:        state.Keyword,
		Data:           state.Data,
		GenerationTime: savedAtProto,
	}, nil
}

// List retrieves AI decision analytics states within a time range
func (s *AIDecisionStateService) List(req *protos.StateListRequest, stream protos.AIDecisionStateService_ListServer) error {
	ctx := stream.Context()
	logger := log.FromContext(ctx).With(
		log.Int(LogKeyAppID, int(req.AppId)),
		log.String(LogKeyKeyword, req.Keyword),
	)
	var savedAtFrom, savedAtTo *time.Time
	if req.GenerationTimeFrom != nil {
		v, _ := ptypes.Timestamp(req.GenerationTimeFrom)
		savedAtFrom = &v
		logger = logger.With(log.Time(LogKeyGenerationTimeFrom, v))
	}
	if req.GenerationTimeTo != nil {
		v, _ := ptypes.Timestamp(req.GenerationTimeTo)
		logger = logger.With(log.Time(LogKeyGenerationTimeTo, v))
		savedAtTo = &v
	}
	ctx = log.WithLogger(ctx, logger)

	if err := s.validateListRequest(ctx, req); err != nil {
		return err
	}
	states, err := s.stateStorage.ListStates(ctx, req.AppId, req.Keyword, savedAtFrom, savedAtTo)
	if err == storage.ErrNotFound {
		return grpc.ErrNotFound(ctx, err)
	} else if err != nil {
		return grpc.ErrUnavailable(ctx, err)
	}

	for _, s := range states {
		genTime, _ := ptypes.TimestampProto(s.SavedAt)
		if err := stream.Send(&protos.State{
			AppId:          s.AppID,
			Keyword:        s.Keyword,
			Data:           s.Data,
			GenerationTime: genTime,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *AIDecisionStateService) validateSaveRequest(ctx context.Context, req *protos.StateSaveRequest) error {
	return validate(ctx,
		validatePositiveInt("app_id", req.AppId),
		validateNonEmptyString("keyword", req.Keyword),
		validateNonEmptyBytes("data", req.Data),
		validateTimestamp("generation_time", req.GenerationTime),
	)
}

func (s *AIDecisionStateService) validateGetRequest(ctx context.Context, req *protos.StateGetRequest) error {
	return validate(ctx,
		validatePositiveInt("app_id", req.AppId),
		validateNonEmptyString("keyword", req.Keyword),
		validateTimestamp("generation_time", req.GenerationTime),
	)
}

func (s *AIDecisionStateService) validateListRequest(ctx context.Context, req *protos.StateListRequest) error {
	return validate(ctx, validatePositiveInt("app_id", req.AppId))
}
