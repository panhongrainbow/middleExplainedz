# MongoDB index

## Introduction

(empty)

## Operation

### In mongosh

```bash
$ mongosh

$ use test

$ db.products.drop
$ db.createCollection("products")

$ for (let i = 0; i < 100000; i++) {
   db.products.insertOne({ 
       item: "product-" + i, 
       price: Math.round(Math.random() * 100),  
       quantity: Math.round(Math.random() * 50)
   })
}

$ db.products.createIndex({item: 1, price: -1})

$ db.products.find({"price": {"$gt": 50}}).explain("executionStats")

$ db.products.dropIndex({item: 1, price: 1})

db.products.getIndexes()
[
  { v: 2, key: { _id: 1 }, name: '_id_' },
  { v: 2, key: { item: 1, price: -1 }, name: 'item_1_price_-1' }
]
```

