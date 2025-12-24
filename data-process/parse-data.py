import json
import uuid
import random
import csv
from datetime import datetime
from itertools import product

def now_iso():
    return {"$date": datetime.utcnow().isoformat() + "Z"}

def gen_id():
    return str(uuid.uuid4())

def parse_product(raw_json: str):
    data = json.loads(raw_json)

    title = data.get("title", "")
    price = data.get("final_price", 0)
    description = data.get("Product Desciption", "")
    image_urls = data.get("image", [])
    variations = data.get("variations", [])

    # -------------------------
    # Images
    # -------------------------
    images = []
    for idx, url in enumerate(image_urls, start=1):
        images.append({
            "_id": gen_id(),
            "url": url,
            "order": idx
        })

    # -------------------------
    # Option Groups
    # -------------------------
    option_groups = []
    option_map = {}

    for opt in variations:
        key = opt["name"].lower()
        values = opt["variations"]

        option_groups.append({
            "key": key,
            "values": values
        })

        option_map[key] = values

    # -------------------------
    # Variants (Cartesian Product)
    # -------------------------
    variant_list = []
    keys = list(option_map.keys())
    combinations = list(product(*option_map.values()))

    for combo in combinations:
        options = dict(zip(keys, combo))
        sku = f"S1-" + "-".join(combo)

        variant_list.append({
            "_id": gen_id(),
            "sku": sku,
            "options": options,
            "price": price,
            "stock": 100,
            "image": random.choice(images)["url"] if images else None
        })

    # -------------------------
    # Final Product Object
    # -------------------------
    product = {
        "_id": gen_id(),
        "name": title,
        "description": description,
        "price": price,
        "images": images,
        "option_groups": option_groups,
        "variants": variant_list,

        # Defaults
        "category_ids": [],
        "seller_category_ids": [],
        "seller_id": "",
        "status": "",
        "is_active": True,
        "rating": 0,
        "rate_count": 0,
        "sold_count": 0,
        "created_at": now_iso(),
        "updated_at": now_iso()
    }

    return product


def parse_product_from_csv_row(row: dict):
    """
    Parse a product from a CSV row with individual columns.
    """
    title = row.get("title", "").strip()
    price = float(row.get("final_price", 0) or 0)
    description = row.get("Product Description", "").strip()
    
    # Parse image URLs (can be JSON array string or comma-separated)
    image_data = row.get("image", "[]").strip()
    try:
        if image_data.startswith('['):
            image_urls = json.loads(image_data)
        else:
            image_urls = [img.strip() for img in image_data.split(',') if img.strip()]
    except:
        image_urls = []
    
    # Parse variations (JSON array string)
    variations_data = row.get("variations", "[]").strip()
    try:
        if variations_data and variations_data != "[]":
            variations = json.loads(variations_data)
        else:
            variations = []
    except:
        variations = []
    
    # -------------------------
    # Images
    # -------------------------
    images = []
    for idx, url in enumerate(image_urls, start=1):
        images.append({
            "_id": gen_id(),
            "url": url,
            "order": idx
        })

    # -------------------------
    # Option Groups
    # -------------------------
    option_groups = []
    option_map = {}

    for opt in variations:
        key = opt.get("name", "").lower()
        values = opt.get("variations", [])
        
        if key and values:
            option_groups.append({
                "key": key,
                "values": values
            })
            option_map[key] = values

    # -------------------------
    # Variants (Cartesian Product)
    # -------------------------
    variant_list = []
    if option_map:
        keys = list(option_map.keys())
        combinations = list(product(*option_map.values()))

        for combo in combinations:
            options = dict(zip(keys, combo))
            sku = f"S1-" + "-".join(combo)

            variant_list.append({
                "_id": gen_id(),
                "sku": sku,
                "options": options,
                "price": price,
                "stock": 100,
                "image": random.choice(images)["url"] if images else None
            })
    else:
        # No variations, create a single default variant
        variant_list.append({
            "_id": gen_id(),
            "sku": f"S1-DEFAULT",
            "options": {},
            "price": price,
            "stock": 100,
            "image": images[0]["url"] if images else None
        })

    # -------------------------
    # Final Product Object
    # -------------------------
    product_obj = {
        "_id": gen_id(),
        "name": title,
        "description": description,
        "price": price,
        "images": images,
        "option_groups": option_groups,
        "variants": variant_list,

        # Defaults
        "category_ids": [],
        "seller_category_ids": [],
        "seller_id": row.get("seller_id", ""),
        "status": "active" if row.get("is_available", "").lower() == "true" else "inactive",
        "is_active": True,
        "rating": float(row.get("rating", 0) or 0),
        "rate_count": 0,
        "sold_count": int(row.get("sold", 0) or 0),
        "created_at": now_iso(),
        "updated_at": now_iso()
    }

    return product_obj


def read_csv_file(csv_path: str):
    """
    Read products from CSV file with individual columns.
    """
    products = []
    
    with open(csv_path, 'r', encoding='utf-8') as file:
        csv_reader = csv.DictReader(file)
        
        # Print available columns for debugging
        fieldnames = csv_reader.fieldnames
        print(f"CSV columns found: {len(fieldnames)} columns")
        
        row_count = 0
        for row in csv_reader:
            row_count += 1
            
            try:
                product_obj = parse_product_from_csv_row(row)
                products.append(product_obj)
                
                if row_count <= 5 or row_count % 100 == 0:
                    print(f"Processed {row_count} products... Latest: {product_obj.get('name', 'Unknown')[:50]}")
            except Exception as e:
                print(f"Error parsing row {row_count}: {e}")
                continue
        
        print(f"Total rows processed: {row_count}")
    
    return products


def export_to_json(products: list, output_path: str):
    """
    Export parsed products to JSON file.
    """
    with open(output_path, 'w', encoding='utf-8') as file:
        json.dump(products, file, indent=2, ensure_ascii=False)
    
    print(f"Exported {len(products)} products to {output_path}")


def export_to_csv(products: list, output_path: str):
    """
    Export parsed products to CSV file (flattened).
    """
    if not products:
        print("No products to export")
        return
    
    with open(output_path, 'w', encoding='utf-8', newline='') as file:
        fieldnames = [
            '_id', 'name', 'description', 'price', 'status', 'is_active',
            'rating', 'rate_count', 'sold_count', 'variant_count', 'image_count'
        ]
        
        writer = csv.DictWriter(file, fieldnames=fieldnames)
        writer.writeheader()
        
        for product in products:
            writer.writerow({
                '_id': product['_id'],
                'name': product['name'],
                'description': product['description'],
                'price': product['price'],
                'status': product['status'],
                'is_active': product['is_active'],
                'rating': product['rating'],
                'rate_count': product['rate_count'],
                'sold_count': product['sold_count'],
                'variant_count': len(product.get('variants', [])),
                'image_count': len(product.get('images', []))
            })
    
    print(f"Exported {len(products)} products to {output_path}")


def main():
    """
    Main function to process CSV and export results.
    """
    # Configuration
    input_csv = "input_products.csv"  # Change to your input CSV file path
    output_json = "output_products.json"
    output_csv = "output_products.csv"
    
    print(f"Reading products from {input_csv}...")
    products = read_csv_file(input_csv)
    
    print(f"Parsed {len(products)} products")
    
    # Export to JSON
    export_to_json(products, output_json)
    
    # Export to CSV
    export_to_csv(products, output_csv)
    
    print("Processing complete!")


if __name__ == "__main__":
    main()
