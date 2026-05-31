# Dataset Schema: ecommerce_products_5k.csv

## Overview

| Property | Value |
|---|---|
| File | `ecommerce_products_5k.csv` |
| Encoding | UTF-8 |
| Delimiter | Comma (`,`) |
| Total rows | 5,000 |
| Total columns | 65 |
| Currency | Vietnamese Dong (VND) |
| Purpose | Content-Based Filtering evaluation for e-commerce recommendation systems |

---

## Column Reference

### Identity & Core Fields

| Column | Type | Description | Example |
|---|---|---|---|
| `product_id` | string | Unique product identifier, format `PRD{5-digit}` | `PRD00001` |
| `name` | string | Product display name: `{Brand} {ProductType}` | `Apple Smartphones` |
| `brand` | string | Manufacturer or brand name | `Samsung`, `Nike`, `IKEA` |

### Category Hierarchy

| Column | Type | Description | Example |
|---|---|---|---|
| `main_category` | string | Top-level category (10 distinct values) | `Electronics` |
| `sub_category` | string | Second-level category (20 distinct values) | `Smartphones` |
| `category_path` | string | Full path, format `{main} > {sub}` | `Electronics > Smartphones` |

**All 20 category paths:**
- `Electronics > Smartphones`
- `Electronics > Laptops`
- `Electronics > Headphones & Audio`
- `Electronics > Smart TVs`
- `Electronics > Tablets`
- `Fashion > Men's Clothing`
- `Fashion > Women's Clothing`
- `Fashion > Footwear`
- `Home & Living > Furniture`
- `Home & Living > Kitchen Appliances`
- `Home & Living > Bedding & Bath`
- `Sports & Outdoors > Fitness Equipment`
- `Sports & Outdoors > Outdoor & Camping`
- `Beauty & Personal Care > Skincare`
- `Beauty & Personal Care > Haircare`
- `Books & Education > Books`
- `Toys & Games > Children's Toys`
- `Food & Beverages > Health Foods`
- `Automotive > Car Accessories`
- `Health & Wellness > Vitamins & Supplements`

Each category contains exactly **250 products**.

### Text & Semantic Fields

| Column | Type | Description | Usage |
|---|---|---|---|
| `description` | string | 1–2 sentence natural language description combining product type, key attributes, and use case | Primary text feature for TF-IDF / embeddings |
| `tags` | string | Pipe-separated (`\|`) keywords: category terms, brand name, product type | Secondary text feature; tokenize by splitting on `\|` |

**Tags example:** `communication|phone|gadget|apple|mobile`  
**Tags cardinality:** 5–8 tags per product.

### Pricing Fields

| Column | Type | Description | Range |
|---|---|---|---|
| `actual_price_vnd` | integer | Original price before discount (VND) | 50,000 – 100,000,000 |
| `discounted_price_vnd` | integer | Final selling price after discount (VND) | 45,000 – 100,000,000 |
| `discount_percentage` | integer | Discount rate applied | 0, 5, 10, 15, 20, 25, 30, 40, 50 |

**Relationship:** `discounted_price_vnd = actual_price_vnd × (1 - discount_percentage / 100)`

**Price ranges by category:**
| Category | Min (VND) | Max (VND) |
|---|---|---|
| Electronics > Smartphones | 1,729,000 | 34,919,000 |
| Electronics > Laptops | 8,026,000 | 79,630,000 |
| Electronics > Smart TVs | 3,478,000 | 99,936,000 |
| Electronics > Tablets | 2,008,000 | 29,934,000 |
| Electronics > Headphones & Audio | 241,000 | 11,987,000 |
| Home & Living > Furniture | 632,000 | 49,661,000 |
| Home & Living > Kitchen Appliances | 274,000 | 14,845,000 |
| Sports & Outdoors > Fitness Equipment | 95,000 | 19,910,000 |
| Sports & Outdoors > Outdoor & Camping | 218,000 | 19,796,000 |
| Automotive > Car Accessories | 61,000 | 9,994,000 |
| Fashion > Footwear | 266,000 | 7,989,000 |
| Fashion > Men's Clothing | 152,000 | 4,994,000 |
| Fashion > Women's Clothing | 138,000 | 5,969,000 |
| Home & Living > Bedding & Bath | 204,000 | 4,945,000 |
| Toys & Games > Children's Toys | 104,000 | 4,981,000 |
| Beauty & Personal Care > Skincare | 114,000 | 3,493,000 |
| Food & Beverages > Health Foods | 206,000 | 2,974,000 |
| Health & Wellness > Vitamins & Supplements | 150,000 | 2,494,000 |
| Beauty & Personal Care > Haircare | 80,000 | 1,996,000 |
| Books & Education > Books | 81,000 | 798,000 |

### Engagement & Inventory Fields

| Column | Type | Description | Range |
|---|---|---|---|
| `rating` | float | Average user rating | 3.0 – 5.0 |
| `review_count` | integer | Number of user reviews | 10 – 50,000 |
| `stock_quantity` | integer | Units available in inventory | 0 – 500 |
| `in_stock` | string | Availability flag | `Yes` / `No` |

---

## Attribute Columns (`attr_*`)

There are **50 attribute columns**, all prefixed with `attr_`. Attributes are **sparse**: each product only populates the attributes relevant to its category. Empty cells represent "not applicable" for that category.

**Important:** Always filter by `category_path` before using attribute columns to avoid comparing sparse/empty values across unrelated categories.

### Attribute Map by Category

#### Electronics > Smartphones
| Column | Possible Values |
|---|---|
| `attr_ram` | 4GB, 6GB, 8GB, 12GB, 16GB |
| `attr_storage` | 64GB, 128GB, 256GB, 512GB |
| `attr_display` | 6.1 inch, 6.4 inch, 6.6 inch, 6.7 inch, 6.8 inch |
| `attr_battery` | 3000mAh, 4000mAh, 4500mAh, 5000mAh, 5500mAh |
| `attr_camera` | 48MP, 50MP, 64MP, 108MP, 200MP |
| `attr_os` | Android 12, Android 13, Android 14, iOS 16, iOS 17 |

#### Electronics > Laptops
| Column | Possible Values |
|---|---|
| `attr_processor` | Intel Core i3/i5/i7/i9, AMD Ryzen 5/7/9, Apple M2, Apple M3 |
| `attr_ram` | 8GB, 16GB, 32GB, 64GB |
| `attr_storage` | 256GB SSD, 512GB SSD, 1TB SSD, 2TB SSD |
| `attr_display` | 13.3 inch FHD, 14 inch FHD, 15.6 inch FHD, 16 inch QHD, 17.3 inch FHD |
| `attr_battery` | 45Wh, 56Wh, 72Wh, 86Wh, 100Wh |

#### Electronics > Headphones & Audio
| Column | Possible Values |
|---|---|
| `attr_type` | Over-ear, On-ear, In-ear, True Wireless |
| `attr_connectivity` | Bluetooth 5.0, Bluetooth 5.3, Wired 3.5mm, USB-C |
| `attr_battery` | 8h, 20h, 30h, 40h, 60h |
| `attr_feature` | Active Noise Cancelling, Ambient Mode, Hi-Res Audio, Spatial Audio |
| `attr_driver` | 9mm, 10mm, 40mm, 50mm |

#### Electronics > Smart TVs
| Column | Possible Values |
|---|---|
| `attr_size` | 32 inch, 43 inch, 50 inch, 55 inch, 65 inch, 75 inch, 85 inch |
| `attr_resolution` | Full HD 1080p, 4K UHD, 8K UHD, QLED 4K, OLED 4K |
| `attr_os` | Android TV, Tizen OS, webOS, Google TV, Fire TV |
| `attr_refresh_rate` | 60Hz, 120Hz, 144Hz |
| `attr_hdr` | HDR10, Dolby Vision, QLED HDR, OLED HDR |

#### Electronics > Tablets
| Column | Possible Values |
|---|---|
| `attr_display` | 8 inch, 10.1 inch, 10.5 inch, 11 inch, 12.4 inch |
| `attr_storage` | 64GB, 128GB, 256GB, 512GB |
| `attr_ram` | 4GB, 6GB, 8GB, 12GB |
| `attr_battery` | 5000mAh, 7040mAh, 8000mAh, 10090mAh |
| `attr_connectivity` | WiFi, WiFi + 4G LTE, WiFi + 5G |

#### Fashion > Men's Clothing
| Column | Possible Values |
|---|---|
| `attr_type` | T-Shirt, Polo Shirt, Dress Shirt, Jacket, Hoodie, Jeans, Chinos, Shorts, Blazer, Sweater |
| `attr_material` | 100% Cotton, Polyester Blend, Linen, Denim, Fleece, Merino Wool |
| `attr_fit` | Slim Fit, Regular Fit, Relaxed Fit, Oversized |
| `attr_color` | Black, White, Navy Blue, Grey, Olive Green, Burgundy, Beige |
| `attr_size` | XS, S, M, L, XL, XXL |

#### Fashion > Women's Clothing
| Column | Possible Values |
|---|---|
| `attr_type` | Dress, Blouse, T-Shirt, Jeans, Skirt, Jacket, Cardigan, Jumpsuit, Leggings, Coat |
| `attr_material` | 100% Cotton, Chiffon, Silk Blend, Linen, Polyester, Knit, Denim |
| `attr_fit` | Regular, Fitted, Loose, Oversized |
| `attr_color` | Black, White, Blush Pink, Sky Blue, Emerald Green, Cream, Red |
| `attr_style` | Casual, Formal, Party, Boho, Streetwear, Minimalist |

#### Fashion > Footwear
| Column | Possible Values |
|---|---|
| `attr_type` | Running Shoes, Casual Sneakers, Boots, Sandals, Loafers, High Heels, Slip-Ons, Trail Shoes, Basketball Shoes, Flip Flops |
| `attr_material` | Leather, Mesh, Canvas, Synthetic, Suede, Knit |
| `attr_sole` | Rubber Sole, EVA Sole, Foam Cushion, Gum Sole |
| `attr_feature` | Waterproof, Breathable, Anti-Slip, Memory Foam, Air Cushion |
| `attr_gender` | Men, Women, Unisex |

#### Home & Living > Furniture
| Column | Possible Values |
|---|---|
| `attr_type` | Sofa, Dining Table, Bed Frame, Wardrobe, Bookshelf, Coffee Table, Office Desk, Dresser, Armchair, TV Console |
| `attr_material` | Solid Wood, MDF Board, Metal Frame, Glass Top, Upholstered Fabric, Faux Leather |
| `attr_color` | Natural Wood, White, Black, Walnut Brown, Grey, Oak |
| `attr_style` | Modern, Scandinavian, Industrial, Traditional, Minimalist, Rustic |
| `attr_dimension` | Small, Medium, Large, Extra Large |

#### Home & Living > Kitchen Appliances
| Column | Possible Values |
|---|---|
| `attr_type` | Rice Cooker, Air Fryer, Blender, Microwave Oven, Coffee Maker, Electric Kettle, Food Processor, Stand Mixer, Toaster Oven, Slow Cooker |
| `attr_capacity` | 1L, 1.5L, 2L, 3L, 4L, 5L, 6L |
| `attr_power` | 600W, 800W, 1000W, 1200W, 1500W, 1800W, 2000W |
| `attr_feature` | Auto Shut-Off, Keep Warm Function, Non-Stick Coating, Digital Display, Timer Function, Multiple Speed Settings |
| `attr_color` | White, Black, Silver, Red, Stainless Steel |

#### Home & Living > Bedding & Bath
| Column | Possible Values |
|---|---|
| `attr_type` | Bed Sheet Set, Duvet Cover, Pillow, Bath Towel Set, Comforter, Mattress Topper, Blanket, Quilt, Bath Mat, Shower Curtain |
| `attr_material` | 100% Cotton, Egyptian Cotton, Microfiber, Bamboo, Flannel, Linen, Percale Cotton |
| `attr_thread_count` | 200TC, 300TC, 400TC, 600TC, 800TC |
| `attr_size` | Single, Queen, King, Double |
| `attr_color` | White, Ivory, Grey, Navy, Blush, Sage Green |

#### Sports & Outdoors > Fitness Equipment
| Column | Possible Values |
|---|---|
| `attr_type` | Yoga Mat, Resistance Bands, Dumbbells, Kettlebell, Jump Rope, Pull-Up Bar, Foam Roller, Ab Wheel, Weight Bench, Barbell Set |
| `attr_material` | Natural Rubber, NBR Foam, Cast Iron, Steel, TPE, PVC |
| `attr_weight` | 2kg, 5kg, 10kg, 15kg, 20kg, 30kg, 50kg |
| `attr_feature` | Anti-Slip Surface, Extra Thick, Adjustable, Foldable, Heavy Duty |
| `attr_level` | Beginner, Intermediate, Advanced, Professional |

#### Sports & Outdoors > Outdoor & Camping
| Column | Possible Values |
|---|---|
| `attr_type` | Tent, Sleeping Bag, Backpack, Camping Stove, Trekking Poles, Headlamp, Water Filter, Hammock, Portable Charger, Compass |
| `attr_capacity` | 1 person, 2 person, 3 person, 4 person, 20L, 30L, 40L, 60L, 80L |
| `attr_season` | 3-Season, 4-Season, Summer, Winter |
| `attr_material` | Ripstop Nylon, Gore-Tex, Polyester, Aluminum, Carbon Fiber |
| `attr_feature` | Waterproof, Lightweight, Ultra-Compact, Wind-Resistant, UV Protection |

#### Beauty & Personal Care > Skincare
| Column | Possible Values |
|---|---|
| `attr_type` | Moisturizer, Serum, Sunscreen, Toner, Cleanser, Eye Cream, Face Mask, Exfoliant, Retinol, Vitamin C Serum |
| `attr_skin_type` | Oily Skin, Dry Skin, Combination Skin, Sensitive Skin, Normal Skin, Acne-Prone Skin |
| `attr_key_ingredient` | Hyaluronic Acid, Niacinamide, Retinol, Vitamin C, Salicylic Acid, Ceramides, Peptides, AHA/BHA |
| `attr_spf` | SPF 15, SPF 30, SPF 50, SPF 50+, No SPF |
| `attr_volume` | 30ml, 50ml, 75ml, 100ml, 150ml, 200ml |

#### Beauty & Personal Care > Haircare
| Column | Possible Values |
|---|---|
| `attr_type` | Shampoo, Conditioner, Hair Mask, Hair Serum, Dry Shampoo, Hair Oil, Leave-In Conditioner, Hair Spray, Heat Protectant, Scalp Treatment |
| `attr_hair_type` | Normal Hair, Dry Hair, Oily Hair, Curly Hair, Color-Treated Hair, Fine Hair, Thick Hair |
| `attr_key_ingredient` | Argan Oil, Keratin, Biotin, Coconut Oil, Collagen, Vitamin E, Hyaluronic Acid |
| `attr_volume` | 200ml, 250ml, 300ml, 400ml, 500ml, 1L |
| `attr_benefit` | Moisturizing, Volumizing, Strengthening, Smoothing, Repairing, Anti-Frizz |

#### Books & Education > Books
| Column | Possible Values |
|---|---|
| `attr_genre` | Self-Help, Business, Fiction, Science Fiction, Biography, History, Technology, Health & Wellness, Psychology, Travel |
| `attr_format` | Hardcover, Paperback, E-book, Audiobook |
| `attr_language` | English, Vietnamese, French, Spanish |
| `attr_pages` | 150-250 pages, 250-350 pages, 350-500 pages, 500+ pages |
| `attr_level` | Beginner, Intermediate, Advanced, All Levels |

#### Toys & Games > Children's Toys
| Column | Possible Values |
|---|---|
| `attr_type` | Building Blocks, Action Figures, Dolls, Board Games, Educational Toys, Remote Control Toys, Stuffed Animals, Puzzles, Arts & Crafts, Outdoor Play |
| `attr_age_group` | 0-2 years, 3-5 years, 6-8 years, 9-12 years, 12+ years |
| `attr_material` | BPA-Free Plastic, Wood, Fabric, Metal, Foam |
| `attr_skill` | Motor Skills, Creativity, Problem Solving, STEM, Language Skills |
| `attr_players` | Solo, 2 Players, 2-4 Players, 4+ Players |

#### Food & Beverages > Health Foods
| Column | Possible Values |
|---|---|
| `attr_type` | Protein Powder, Protein Bar, Granola, Green Superfood, Probiotic, Multivitamin, Omega-3, Collagen Powder, Meal Replacement, Energy Bar |
| `attr_flavor` | Chocolate, Vanilla, Strawberry, Unflavored, Mixed Berry, Peanut Butter, Matcha |
| `attr_dietary` | Vegan, Gluten-Free, Keto-Friendly, Organic, Non-GMO, Dairy-Free |
| `attr_serving` | 15 servings, 20 servings, 30 servings, 60 servings |
| `attr_protein` | 15g protein, 20g protein, 25g protein, 30g protein |

#### Automotive > Car Accessories
| Column | Possible Values |
|---|---|
| `attr_type` | Car Seat Cover, Dashboard Camera, Car Air Freshener, Phone Mount, Car Vacuum, Tire Inflator, Jump Starter, Car Mat, Steering Wheel Cover, GPS Navigator |
| `attr_compatibility` | Universal Fit, Sedan, SUV, Truck, Hatchback |
| `attr_material` | Neoprene, Leather, Mesh, Silicone, Rubber |
| `attr_feature` | Waterproof, Easy Install, 360-Degree Rotation, Wide Angle, Wireless |
| `attr_power` | 12V DC, USB Powered, Battery Powered, Wireless |

#### Health & Wellness > Vitamins & Supplements
| Column | Possible Values |
|---|---|
| `attr_type` | Multivitamin, Vitamin D3, Vitamin C, Omega-3 Fish Oil, Magnesium, Zinc, B-Complex, Probiotics, Turmeric, Collagen |
| `attr_form` | Capsule, Tablet, Softgel, Gummy, Powder, Liquid |
| `attr_dosage` | 500mg, 1000mg, 1500mg, 2000mg, 5000IU, 10000IU |
| `attr_count` | 30 count, 60 count, 90 count, 120 count, 180 count, 250 count |
| `attr_certification` | GMP Certified, USP Verified, NSF Certified, Non-GMO Verified, Organic |

---

## Recommended Feature Engineering for Content-Based Filtering

### Text Features
```
text_combined = description + " " + tags.replace("|", " ")
```
Apply TF-IDF or sentence embeddings (e.g., `paraphrase-multilingual-MiniLM`) on `text_combined`.

### Categorical Features (One-Hot Encode)
- `main_category`, `sub_category`, `brand`
- All populated `attr_*` columns within the same category group

### Numerical Features (Normalize with MinMaxScaler or StandardScaler)
- `actual_price_vnd` or `discounted_price_vnd`
- `rating`
- `review_count` (consider log-transform: `log1p(review_count)`)
- `discount_percentage`

### Similarity Strategy
- Within-category comparison: use full attribute vector + text features
- Cross-category comparison: use only `main_category`, `brand`, text features, and price
- Recommended similarity metric: **Cosine Similarity** on the combined feature vector

### Known Sparsity Pattern
`attr_*` columns are sparse by design. A product in `Electronics > Smartphones` will have `attr_ram`, `attr_storage`, etc. populated, but columns like `attr_genre`, `attr_skin_type`, `attr_flavor` will be empty. **Do not impute cross-category attributes.** Filter by `category_path` first when building category-specific feature vectors.