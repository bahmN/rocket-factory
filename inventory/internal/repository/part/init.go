package part

import (
	"time"

	repoModel "github.com/bahmN/rocket-factory/inventory/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) InitTestData() {
	parts := []*repoModel.Part{
		{
			UUID:          uuid.NewString(),
			Name:          "Двигатель",
			Description:   "Мощный ракетный двигатель",
			Price:         15000.0,
			StockQuantity: 5,
			Category:      1, // например, 1 = ENGINE
			Dimensions: &repoModel.Dimensions{
				Length: 200,
				Width:  100,
				Height: 120,
				Weight: 300,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "RocketMotors",
				Cuntry:  "Russia",
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
			Category:      2, // например, 2 = FUEL
			Dimensions: &repoModel.Dimensions{
				Length: 150,
				Width:  150,
				Height: 200,
				Weight: 200,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "FuelTech",
				Cuntry:  "Germany",
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
			Category:      3, // например, 3 = PORTHOLE
			Dimensions: &repoModel.Dimensions{
				Length: 50,
				Width:  50,
				Height: 10,
				Weight: 20,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "GlassSpace",
				Cuntry:  "USA",
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
			Category:      4, // например, 4 = WING
			Dimensions: &repoModel.Dimensions{
				Length: 300,
				Width:  50,
				Height: 20,
				Weight: 50,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "WingPro",
				Cuntry:  "France",
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
			Category:      0, // например, 0 = UNKNOWN
			Dimensions: &repoModel.Dimensions{
				Length: 80,
				Width:  40,
				Height: 10,
				Weight: 10,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "ControlSys",
				Cuntry:  "Japan",
				Website: "https://controlsys.example.com",
			},
			Tags:      []string{"электроника", "панель"},
			Metadata:  map[string]string{"версия": "2.1"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, part := range parts {
		r.data[part.UUID] = part
	}
}
