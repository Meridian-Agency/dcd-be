package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type ServiceSeed struct {
	Slug        string
	ParentSlug  string
	Name        string
	Category    string
	BasePrice   float64
	DurationMin int
	Description string
}

var servicesSeed = []ServiceSeed{
	// Parent Services
	{
		Slug:        "deep-clean-valet",
		Name:        "Deep Clean Valet",
		Category:    "VALETING",
		BasePrice:   140.00,
		DurationMin: 240,
		Description: "A full interior and exterior reset — every panel, vent and seam.",
	},
	{
		Slug:        "maintenance-valet",
		Name:        "Maintenance Valet",
		Category:    "VALETING",
		BasePrice:   70.00,
		DurationMin: 120,
		Description: "Keep a detailed car looking its best between bigger jobs.",
	},
	{
		Slug:        "single-stage-enhancement",
		Name:        "Single-Stage Enhancement",
		Category:    "DETAILING",
		BasePrice:   450.00,
		DurationMin: 480,
		Description: "Restore gloss and remove the majority of light swirls in one stage.",
	},
	{
		Slug:        "two-stage-paint-correction",
		Name:        "Two-Stage Paint Correction",
		Category:    "DETAILING",
		BasePrice:   650.00,
		DurationMin: 960,
		Description: "The deepest correction we offer — a near-flawless, better-than-showroom finish.",
	},
	{
		Slug:        "gtechniq-ceramic-coating",
		Name:        "New Car Protection & Ceramic Coating",
		Category:    "PROTECTION",
		BasePrice:   495.00,
		DurationMin: 480,
		Description: "Long-term gloss, easier cleaning and serious protection.",
	},
	{
		Slug:        "paint-protection-film",
		Name:        "Paint Protection Film (PPF)",
		Category:    "PROTECTION",
		BasePrice:   0.00,
		DurationMin: 0,
		Description: "SunTek film that shields paint from stone chips and scratches.",
	},
	{
		Slug:        "alloy-wheel-refurbishment",
		Name:        "Alloy Wheel Refurbishment",
		Category:    "UPGRADES",
		BasePrice:   75.00,
		DurationMin: 180,
		Description: "Kerbed or corroded alloys brought back to factory finish.",
	},
	{
		Slug:        "vehicle-upgrades",
		Name:        "Vehicle Security",
		Category:    "UPGRADES",
		BasePrice:   249.00,
		DurationMin: 120,
		Description: "Immobilisers and Thatcham-approved trackers, fitted to a factory standard.",
	},
	{
		Slug:        "apple-carplay",
		Name:        "Apple CarPlay & Android Auto",
		Category:    "UPGRADES",
		BasePrice:   299.00,
		DurationMin: 60,
		Description: "Genuine wireless CarPlay, coded or retrofitted to an OEM standard.",
	},
	// PPF Subtypes
	{
		Slug:        "ppf-front-end",
		ParentSlug:  "paint-protection-film",
		Name:        "Front End Coverage",
		Category:    "PROTECTION",
		BasePrice:   1395.00,
		DurationMin: 480,
		Description: "Protects high-impact areas: front bumper, hood, fenders, and mirrors.",
	},
	{
		Slug:        "ppf-extended",
		ParentSlug:  "paint-protection-film",
		Name:        "Extended Coverage",
		Category:    "PROTECTION",
		BasePrice:   1895.00,
		DurationMin: 720,
		Description: "Adds protection for side skirts, lower doors, and rear splash areas.",
	},
	{
		Slug:        "ppf-full-body",
		ParentSlug:  "paint-protection-film",
		Name:        "Full Body Coverage",
		Category:    "PROTECTION",
		BasePrice:   3500.00,
		DurationMin: 960,
		Description: "Complete vehicle coverage for absolute peace of mind.",
	},
}

func SeedServices(db *gorm.DB) error {
	log.Println("Seeding service packages...")

	// First pass: upsert parent services (those without a ParentSlug)
	for _, seed := range servicesSeed {
		if seed.ParentSlug != "" {
			continue
		}

		var existing ServicePackage
		err := db.Where("slug = ?", seed.Slug).First(&existing).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("error checking existing service %s: %w", seed.Slug, err)
		}

		if err == gorm.ErrRecordNotFound {
			// Create new
			newService := ServicePackage{
				Slug:        seed.Slug,
				Name:        seed.Name,
				Category:    seed.Category,
				BasePrice:   seed.BasePrice,
				DurationMin: seed.DurationMin,
				Description: seed.Description,
			}
			if err := db.Create(&newService).Error; err != nil {
				return fmt.Errorf("failed to create service %s: %w", seed.Slug, err)
			}
			log.Printf("Created service package: %s", seed.Slug)
		} else {
			// Update existing
			existing.Name = seed.Name
			existing.Category = seed.Category
			existing.BasePrice = seed.BasePrice
			existing.DurationMin = seed.DurationMin
			existing.Description = seed.Description
			if err := db.Save(&existing).Error; err != nil {
				return fmt.Errorf("failed to update service %s: %w", seed.Slug, err)
			}
			log.Printf("Updated service package: %s", seed.Slug)
		}
	}

	// Second pass: upsert child services (those with a ParentSlug)
	for _, seed := range servicesSeed {
		if seed.ParentSlug == "" {
			continue
		}

		// Find parent ID
		var parent ServicePackage
		if err := db.Where("slug = ?", seed.ParentSlug).First(&parent).Error; err != nil {
			return fmt.Errorf("failed to find parent service %s for %s: %w", seed.ParentSlug, seed.Slug, err)
		}

		var existing ServicePackage
		err := db.Where("slug = ?", seed.Slug).First(&existing).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("error checking existing service subtype %s: %w", seed.Slug, err)
		}

		if err == gorm.ErrRecordNotFound {
			// Create new child
			newService := ServicePackage{
				ParentID:    &parent.ID,
				Slug:        seed.Slug,
				Name:        seed.Name,
				Category:    seed.Category,
				BasePrice:   seed.BasePrice,
				DurationMin: seed.DurationMin,
				Description: seed.Description,
			}
			if err := db.Create(&newService).Error; err != nil {
				return fmt.Errorf("failed to create service subtype %s: %w", seed.Slug, err)
			}
			log.Printf("Created service subtype: %s (parent: %s)", seed.Slug, parent.Slug)
		} else {
			// Update existing child
			existing.ParentID = &parent.ID
			existing.Name = seed.Name
			existing.Category = seed.Category
			existing.BasePrice = seed.BasePrice
			existing.DurationMin = seed.DurationMin
			existing.Description = seed.Description
			if err := db.Save(&existing).Error; err != nil {
				return fmt.Errorf("failed to update service subtype %s: %w", seed.Slug, err)
			}
			log.Printf("Updated service subtype: %s", seed.Slug)
		}
	}

	log.Println("Service packages seeding complete.")
	return nil
}
