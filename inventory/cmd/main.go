package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	inventoryV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = ":50052"

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	mu          sync.RWMutex
	inventories map[string]*inventoryV1.Part
}

func (s *InventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	part := s.inventories[req.Uuid]
	if part == nil {
		return nil, status.Errorf(codes.NotFound, "part: not found")
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	parts := make([]*inventoryV1.Part, 0, len(s.inventories))
	for _, part := range s.inventories {
		parts = append(parts, part)
	}
	s.mu.RUnlock()

	filter := req.GetFilter()
	if filter == nil {
		return &inventoryV1.ListPartsResponse{Parts: parts}, nil
	}

	if len(filter.Uuids) > 0 {
		uuidSet := make(map[string]struct{}, len(filter.Uuids))
		for _, u := range filter.Uuids {
			uuidSet[u] = struct{}{}
		}
		tmp := parts[:0]
		for _, part := range parts {
			if _, ok := uuidSet[part.Uuid]; ok {
				tmp = append(tmp, part)
			}
		}
		parts = tmp
	}

	if len(filter.Names) > 0 {
		nameSet := make(map[string]struct{}, len(filter.Names))
		for _, n := range filter.Names {
			nameSet[n] = struct{}{}
		}
		tmp := parts[:0]
		for _, part := range parts {
			if _, ok := nameSet[part.Name]; ok {
				tmp = append(tmp, part)
			}
		}
		parts = tmp
	}

	if len(filter.Categories) > 0 {
		catSet := make(map[inventoryV1.Category]struct{}, len(filter.Categories))
		for _, c := range filter.Categories {
			catSet[c] = struct{}{}
		}
		tmp := parts[:0]
		for _, part := range parts {
			if _, ok := catSet[part.Category]; ok {
				tmp = append(tmp, part)
			}
		}
		parts = tmp
	}

	if len(filter.ManufacturerCountries) > 0 {
		countrySet := make(map[string]struct{}, len(filter.ManufacturerCountries))
		for _, c := range filter.ManufacturerCountries {
			countrySet[c] = struct{}{}
		}
		tmp := parts[:0]
		for _, part := range parts {
			if part.Manufacturer != nil {
				if _, ok := countrySet[part.Manufacturer.Cuntry]; ok {
					tmp = append(tmp, part)
				}
			}
		}
		parts = tmp
	}

	if len(filter.Tags) > 0 {
		tagSet := make(map[string]struct{}, len(filter.Tags))
		for _, t := range filter.Tags {
			tagSet[t] = struct{}{}
		}
		tmp := parts[:0]
		for _, part := range parts {
			found := false
			for _, tag := range part.Tags {
				if _, ok := tagSet[tag]; ok {
					found = true
					break
				}
			}
			if found {
				tmp = append(tmp, part)
			}
		}
		parts = tmp
	}

	if len(parts) == 0 {
		return nil, status.Errorf(codes.NotFound, "parts: not found")
	}

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}

func (s *InventoryService) seedTestData() {
	parts := []*inventoryV1.Part{
		{
			Uuid:          uuid.NewString(),
			Name:          "Двигатель",
			Description:   "Мощный ракетный двигатель",
			Price:         15000.0,
			StockQuantity: 5,
			Category:      inventoryV1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 2.0, Width: 1.0, Height: 1.2, Weight: 300.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name: "RocketMotors", Cuntry: "Russia", Website: "https://rocketmotors.example.com",
			},
			Tags:      []string{"основной", "мотор"},
			Metadata:  map[string]*structpb.Value{"серия": structpb.NewStringValue("X100")},
			CreatedAt: timestamppb.Now(), UpdatedAt: timestamppb.Now(),
		},
		{
			Uuid:          uuid.NewString(),
			Name:          "Топливный бак",
			Description:   "Бак для хранения топлива",
			Price:         8000.0,
			StockQuantity: 8,
			Category:      inventoryV1.Category_CATEGORY_FUEL,
			Dimensions: &inventoryV1.Dimensions{
				Length: 1.5, Width: 1.5, Height: 2.0, Weight: 200.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name: "FuelTech", Cuntry: "Germany", Website: "https://fueltech.example.com",
			},
			Tags:      []string{"топливо", "бак"},
			Metadata:  map[string]*structpb.Value{"материал": structpb.NewStringValue("титан")},
			CreatedAt: timestamppb.Now(), UpdatedAt: timestamppb.Now(),
		},
		{
			Uuid:          uuid.NewString(),
			Name:          "Иллюминатор",
			Description:   "Прочный иллюминатор для ракеты",
			Price:         3000.0,
			StockQuantity: 15,
			Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 0.5, Width: 0.5, Height: 0.1, Weight: 20.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name: "GlassSpace", Cuntry: "USA", Website: "https://glassspace.example.com",
			},
			Tags:      []string{"стекло", "иллюминатор"},
			Metadata:  map[string]*structpb.Value{"прозрачность": structpb.NewStringValue("99%")},
			CreatedAt: timestamppb.Now(), UpdatedAt: timestamppb.Now(),
		},
		{
			Uuid:          uuid.NewString(),
			Name:          "Крыло",
			Description:   "Аэродинамическое крыло",
			Price:         5000.0,
			StockQuantity: 12,
			Category:      inventoryV1.Category_CATEGORY_WING,
			Dimensions: &inventoryV1.Dimensions{
				Length: 3.0, Width: 0.5, Height: 0.2, Weight: 50.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name: "WingPro", Cuntry: "France", Website: "https://wingpro.example.com",
			},
			Tags:      []string{"крыло", "аэродинамика"},
			Metadata:  map[string]*structpb.Value{"тип": structpb.NewStringValue("стабилизатор")},
			CreatedAt: timestamppb.Now(), UpdatedAt: timestamppb.Now(),
		},
		{
			Uuid:          uuid.NewString(),
			Name:          "Панель управления",
			Description:   "Электронная панель управления",
			Price:         7000.0,
			StockQuantity: 7,
			Category:      inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED,
			Dimensions: &inventoryV1.Dimensions{
				Length: 0.8, Width: 0.4, Height: 0.1, Weight: 10.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name: "ControlSys", Cuntry: "Japan", Website: "https://controlsys.example.com",
			},
			Tags:      []string{"электроника", "панель"},
			Metadata:  map[string]*structpb.Value{"версия": structpb.NewStringValue("2.1")},
			CreatedAt: timestamppb.Now(), UpdatedAt: timestamppb.Now(),
		},
	}

	for _, part := range parts {
		s.inventories[part.Uuid] = part
	}
}

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	inventories := map[string]*inventoryV1.Part{}

	service := &InventoryService{
		inventories: inventories,
	}

	service.seedTestData()

	inventoryV1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("gRPC server listening on %s\n", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}
