package part

import (
	"context"
	"time"

	repoModel "github.com/bahmN/rocket-factory/inventory/internal/repository/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) InitTestData(ctx context.Context) error {
	parts := []*repoModel.Part{
		{
			UUID:          uuid.NewString(),
			Name:          "Двигатель",
			Description:   "Мощный ракетный двигатель",
			Price:         15000.0,
			StockQuantity: 5,
			Category:      1, // ENGINE
			Dimensions: &repoModel.Dimensions{
				Length: 200,
				Width:  100,
				Height: 120,
				Weight: 300,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "RocketMotors",
				Country: "Russia",
				Website: "https://rocketmotors.example.com",
			},
			Tags:      []string{"основной", "мотор"},
			Metadata:  map[string]string{"серия": "X100"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Топливный бак",
			Description:   "Бак для хранения топлива",
			Price:         8000.0,
			StockQuantity: 8,
			Category:      2, // FUEL
			Dimensions: &repoModel.Dimensions{
				Length: 150,
				Width:  150,
				Height: 200,
				Weight: 200,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "FuelTech",
				Country: "Germany",
				Website: "https://fueltech.example.com",
			},
			Tags:      []string{"топливо", "бак"},
			Metadata:  map[string]string{"материал": "титан"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Иллюминатор",
			Description:   "Прочный иллюминатор для ракеты",
			Price:         3000.0,
			StockQuantity: 15,
			Category:      3, // PORTHOLE
			Dimensions: &repoModel.Dimensions{
				Length: 50,
				Width:  50,
				Height: 10,
				Weight: 20,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "GlassSpace",
				Country: "USA",
				Website: "https://glassspace.example.com",
			},
			Tags:      []string{"стекло", "иллюминатор"},
			Metadata:  map[string]string{"прозрачность": "99%"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Крыло",
			Description:   "Аэродинамическое крыло",
			Price:         5000.0,
			StockQuantity: 12,
			Category:      4, // WING
			Dimensions: &repoModel.Dimensions{
				Length: 300,
				Width:  50,
				Height: 20,
				Weight: 50,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "WingPro",
				Country: "France",
				Website: "https://wingpro.example.com",
			},
			Tags:      []string{"крыло", "аэродинамика"},
			Metadata:  map[string]string{"тип": "стабилизатор"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Панель управления",
			Description:   "Электронная панель управления",
			Price:         7000.0,
			StockQuantity: 7,
			Category:      0, // UNKNOWN
			Dimensions: &repoModel.Dimensions{
				Length: 80,
				Width:  40,
				Height: 10,
				Weight: 10,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "ControlSys",
				Country: "Japan",
				Website: "https://controlsys.example.com",
			},
			Tags:      []string{"электроника", "панель"},
			Metadata:  map[string]string{"версия": "2.1"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Очистим коллекцию перед вставкой, чтобы не дублировать
	if _, err := r.coll.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}

	// Преобразуем []interface{} для bulk вставки
	docs := make([]interface{}, len(parts))
	for i, part := range parts {
		docs[i] = part
	}

	// Вставляем все документы
	_, err := r.coll.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}
