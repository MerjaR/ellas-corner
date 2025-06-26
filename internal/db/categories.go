package db

// CuratedBabyBox provides hardcoded example items grouped by category.
// This is used to render themed suggestions on the baby box feature page.
// Eventually, this could be moved or an admin could update it as part of other categories.
// For now, this is a simple structure to demonstrate how the babybox page would work.

type BabyBoxItem struct {
	Title string
	Image string
}

var CuratedBabyBox = map[string][]BabyBoxItem{
	"Winter clothes": {
		{Title: "Winter Outfit", Image: "winter1.jpg"},
		{Title: "Warm Hat", Image: "winter2.png"},
		{Title: "Thermal Onesie", Image: "winter3.jpg"},
	},
	"Baby's health": {
		{Title: "Infrared Forehead Thermometer", Image: "baby1.jpg"},
		{Title: "Nasal Aspirator", Image: "baby2.jpg"},
		{Title: "Nail Scissors", Image: "baby3.jpg"},
	},
	"Mother's health and comfort": {
		{Title: "Nursing Pillow", Image: "mother1.jpg"},
		{Title: "Cream", Image: "mother2.jpg"},
		{Title: "Relaxation Tea", Image: "mother3.jpg"},
	},
	"Development books": {
		{Title: "Baby Development", Image: "book1.jpg"},
		{Title: "New Mother Mindset", Image: "book2.jpg"},
		{Title: "Parenting Psychology", Image: "book3.jpg"},
	},
	"Recommended for parents": {
		{Title: "Noise Cancelling Headphones", Image: "parents1.jpg"},
		{Title: "Decaf coffee", Image: "parents2.jpg"},
		{Title: "Baby Footprint Kit", Image: "parents3.jpg"},
	},
	"Travel essentials": {
		{Title: "Diaper Bag", Image: "travel1.jpg"},
		{Title: "Travel Bed", Image: "travel2.jpg"},
		{Title: "Play mat", Image: "travel3.jpg"},
	},
}
