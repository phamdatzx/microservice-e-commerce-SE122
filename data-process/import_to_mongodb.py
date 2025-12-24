#!/usr/bin/env python3
"""
Import JSON data into MongoDB collection
"""

import json
from pymongo import MongoClient
from pymongo.errors import ConnectionFailure, BulkWriteError
import os
from typing import List, Dict, Any
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()


def load_json_file(filepath: str) -> List[Dict[str, Any]]:
    """
    Load data from JSON file
    
    Args:
        filepath: Path to the JSON file
        
    Returns:
        List of dictionaries containing the data
    """
    print(f"Loading data from {filepath}...")
    
    try:
        with open(filepath, 'r', encoding='utf-8') as file:
            data = json.load(file)
        print(f"Successfully loaded {len(data)} records")
        return data
    except FileNotFoundError:
        print(f"Error: File {filepath} not found")
        return []
    except json.JSONDecodeError as e:
        print(f"Error: Invalid JSON format - {e}")
        return []
    except Exception as e:
        print(f"Error loading file: {e}")
        return []


def connect_to_mongodb(connection_string: str = None, host: str = 'localhost', port: int = 27017):
    """
    Connect to MongoDB
    
    Args:
        connection_string: MongoDB connection string (optional)
        host: MongoDB host (default: localhost)
        port: MongoDB port (default: 27017)
        
    Returns:
        MongoClient instance or None if connection fails
    """
    try:
        if connection_string:
            client = MongoClient(connection_string)
        else:
            client = MongoClient(host, port)
        
        # Test the connection
        client.admin.command('ping')
        print("Successfully connected to MongoDB")
        return client
    except ConnectionFailure:
        print("Error: Failed to connect to MongoDB. Make sure MongoDB is running.")
        return None
    except Exception as e:
        print(f"Error connecting to MongoDB: {e}")
        return None


def import_data_to_mongodb(
    data: List[Dict[str, Any]], 
    client: MongoClient, 
    database_name: str, 
    collection_name: str,
    batch_size: int = 1000
) -> bool:
    """
    Import data into MongoDB collection
    
    Args:
        data: List of documents to import
        client: MongoClient instance
        database_name: Name of the database
        collection_name: Name of the collection
        batch_size: Number of documents to insert per batch
        
    Returns:
        True if successful, False otherwise
    """
    try:
        db = client[database_name]
        collection = db[collection_name]
        
        print(f"Importing data to database '{database_name}', collection '{collection_name}'...")
        
        # Insert in batches for better performance
        total_inserted = 0
        for i in range(0, len(data), batch_size):
            batch = data[i:i + batch_size]
            try:
                result = collection.insert_many(batch, ordered=False)
                total_inserted += len(result.inserted_ids)
                print(f"Inserted batch {i//batch_size + 1}: {len(result.inserted_ids)} documents")
            except BulkWriteError as e:
                # Some documents might have been inserted even if there was an error
                total_inserted += e.details.get('nInserted', 0)
                print(f"Warning: Batch {i//batch_size + 1} had errors. Inserted: {e.details.get('nInserted', 0)}")
        
        print(f"\nTotal documents inserted: {total_inserted} out of {len(data)}")
        
        # Show collection statistics
        doc_count = collection.count_documents({})
        print(f"Collection '{collection_name}' now contains {doc_count} documents")
        
        return True
    except Exception as e:
        print(f"Error importing data: {e}")
        return False


def main():
    """
    Main function to run the import process
    """
    # Configuration
    JSON_FILE = 'output_products.json'
    
    # MongoDB connection settings - Modify these as needed
    MONGO_CONNECTION_STRING = os.getenv('MONGO_CONNECTION_STRING', None)
    MONGO_HOST = os.getenv('MONGO_HOST', 'localhost')
    MONGO_PORT = int(os.getenv('MONGO_PORT', 27017))
    
    # Database and collection settings
    DATABASE_NAME = os.getenv('MONGO_DATABASE', 'products_db')
    COLLECTION_NAME = os.getenv('MONGO_COLLECTION', 'products')
    
    print("=" * 60)
    print("MongoDB Data Import Tool")
    print("=" * 60)
    print(f"JSON File: {JSON_FILE}")
    print(f"Database: {DATABASE_NAME}")
    print(f"Collection: {COLLECTION_NAME}")
    print("=" * 60)
    print()
    
    # Load JSON data
    data = load_json_file(JSON_FILE)
    if not data:
        print("No data to import. Exiting.")
        return
    
    # Connect to MongoDB
    client = connect_to_mongodb(
        connection_string=MONGO_CONNECTION_STRING,
        host=MONGO_HOST,
        port=MONGO_PORT
    )
    
    if not client:
        print("Failed to connect to MongoDB. Exiting.")
        return
    
    try:
        # Import data
        success = import_data_to_mongodb(
            data=data,
            client=client,
            database_name=DATABASE_NAME,
            collection_name=COLLECTION_NAME,
            batch_size=1000
        )
        
        if success:
            print("\n✓ Data import completed successfully!")
        else:
            print("\n✗ Data import failed.")
    
    finally:
        # Close the connection
        client.close()
        print("\nMongoDB connection closed.")


if __name__ == "__main__":
    main()
