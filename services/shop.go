package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/hackclub/hackatime/config"
	"github.com/hackclub/hackatime/models"
	"github.com/patrickmn/go-cache"
)

type ShopService struct {
	config *config.Config
	cache  *cache.Cache
}

func NewShopService() *ShopService {
	return &ShopService{
		config: config.Get(),
		cache:  cache.New(1*time.Minute, 1*time.Minute),
	}
}

type HackClubProduct struct {
	Name        string `json:"name"`
	SmallName   string `json:"smallName"`
	Description string `json:"description"`
	Hours       int    `json:"hours"`
	ImageURL    string `json:"imageURL"`
	Stock       *int   `json:"stock"` // Change to pointer to allow null
}

func (srv *ShopService) GetProducts() ([]*models.Product, error) {
	// Check if products are in cache
	if cachedProducts, found := srv.cache.Get("products"); found {
		return cachedProducts.([]*models.Product), nil
	}

	// Fetch products from Hack Club API
	resp, err := http.Get("https://hackclub.com/api/arcade/shop/")
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %v", err)
	}
	defer resp.Body.Close()

	var hackClubProducts []HackClubProduct
	if err := json.NewDecoder(resp.Body).Decode(&hackClubProducts); err != nil {
		return nil, fmt.Errorf("error decoding products: %v", err)
	}

	formattedProducts := []*models.Product{}

	for i, product := range hackClubProducts {

		stock := -1
		if product.Stock != nil {
			stock = *product.Stock
		}

		formattedProducts = append(formattedProducts, &models.Product{
			ID:          uint(i + 1), // Use index + 1 as ID
			Name:        product.Name,
			SmallName:   product.SmallName,
			Description: product.Description,
			Price:       product.Hours,
			Stock:       stock,
			Image:       product.ImageURL,
		})
	}

	// Sort products by price
	sort.Slice(formattedProducts, func(i, j int) bool {
		return formattedProducts[i].Price < formattedProducts[j].Price
	})

	// Cache the sorted formatted products
	srv.cache.Set("products", formattedProducts, cache.DefaultExpiration)

	return formattedProducts, nil
}

func (srv *ShopService) ClearProductsCache() {
	srv.cache.Delete("products")
}
