package handler

import (
	"context"
	"math/rand"
	"starting_grpc/api/gen/api"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	// パンケーキの仕上がりに影響する seed を初期化します。
	rand.Seed(time.Now().UnixNano())
}

// BakerHandler はケーキを焼きます
type BakerHandler struct {
	report *report
}

type report struct {
	sync.Mutex // 複数人が同時に焼いても大丈夫にしておきます
	data       map[api.Pancake_Menu]int
}

// NewBakerHandler は Bakerhandlerを初期化します
func NewBakerHandler() *BakerHandler {
	return &BakerHandler{
		report: &report{
			data: make(map[api.Pancake_Menu]int),
		},
	}
}

// Bake は指定されたメニューのパンケーキを焼いて、焼けたパンを Response として返します
func (h *BakerHandler) Bake(ctx context.Context, req *api.BakeRequest) (*api.BakeResponse, error) {
	// fmt.Printf("Baked a pancake for %v !\n", ctx.Value("UserName"))

	// バリデーション
	if req.Menu == api.Pancake_UNKNOWN || req.Menu > api.Pancake_SPICY_CURRY {
		return nil, status.Errorf(codes.InvalidArgument, "パンケーキを選んでください！")
	}

	// パンケーキを焼いて、数を記録します
	now := time.Now()
	h.report.Lock()
	h.report.data[req.Menu] = h.report.data[req.Menu] + 1
	h.report.Unlock()

	// レスポンスを作って返す
	return &api.BakeResponse{
		Pancake: &api.Pancake{
			Menu:           req.Menu,
			ChefName:       "gami", // ワンオペ
			TechnicalScore: rand.Float32(),
			CreateTime: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		},
	}, nil
}

// Report は、焼けたパンの数を申告します
func (h *BakerHandler) Report(ctx context.Context, req *api.ReportRequest) (*api.ReportResponse, error) {
	counts := make([]*api.Report_BakeCount, len(h.report.data))

	// レポートを作ります
	h.report.Lock()
	for k, v := range h.report.data {
		counts = append(counts, &api.Report_BakeCount{
			Menu:  k,
			Count: int32(v),
		})
	}
	h.report.Unlock()

	// レスポンスを作って返す
	return &api.ReportResponse{
		Report: &api.Report{
			BakeCounts: counts,
		},
	}, nil
}
