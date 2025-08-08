#!/usr/bin/env python3
"""
Database Structure Check for Blog Service
"""

import pymysql
import sys

def check_database():
    try:
        # Connect to MySQL database
        connection = pymysql.connect(
            host='65.1.94.25',
            user='phpmyadmin',
            password='mFVarH2LCrQK',
            database='mejona_unified',
            charset='utf8mb4'
        )
        
        cursor = connection.cursor()
        
        print("Connected to production database successfully!")
        print("\nEXISTING TABLES:")
        print("=" * 50)
        
        # Show all tables
        cursor.execute("SHOW TABLES")
        tables = cursor.fetchall()
        
        for table in tables:
            print(f"  - {table[0]}")
        
        print(f"\nTotal tables: {len(tables)}")
        
        # Check if blogs table exists
        cursor.execute("SHOW TABLES LIKE 'blogs'")
        blog_table = cursor.fetchone()
        
        if blog_table:
            print("\nEXISTING BLOGS TABLE STRUCTURE:")
            print("=" * 50)
            cursor.execute("DESCRIBE blogs")
            columns = cursor.fetchall()
            
            for column in columns:
                print(f"  {column[0]:25} {column[1]:15} {column[2]:8} {column[3]:8}")
            
            print(f"\nTotal columns in blogs table: {len(columns)}")
            
            # Check for data
            cursor.execute("SELECT COUNT(*) FROM blogs")
            count = cursor.fetchone()[0]
            print(f"Existing blog records: {count}")
            
        else:
            print("\nNo 'blogs' table found - will need to create")
        
        # Check contacts table for CRM integration
        cursor.execute("SHOW TABLES LIKE 'contacts'")
        contacts_table = cursor.fetchone()
        
        if contacts_table:
            print("\nCONTACTS TABLE (for CRM integration):")
            cursor.execute("SELECT COUNT(*) FROM contacts")
            contact_count = cursor.fetchone()[0]
            print(f"  Existing contacts: {contact_count}")
        
        cursor.close()
        connection.close()
        
        return True
        
    except Exception as e:
        print(f"Database connection failed: {str(e)}")
        return False

if __name__ == "__main__":
    if check_database():
        print("\nDatabase is ready for blog service implementation!")
    else:
        print("\nNeed to fix database connection before proceeding")
        sys.exit(1)