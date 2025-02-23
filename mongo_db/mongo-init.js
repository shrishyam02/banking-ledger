
db = db.getSiblingDB("banking_ledger_db");

db.createCollection("transactions");

db.transactions.insertMany([
  {
    account_id: "account1",
    amount: 100.00,
    transaction_type: "deposit",
    timestamp: new Date(), // Current timestamp
    details: { description: "Initial deposit" },
    status: "completed"
  },
  {
    account_id: "account2",
    amount: 50.00,
    transaction_type: "deposit",
    timestamp: new Date(),
    details: { description: "Initial deposit" },
    status: "completed"
  }
], { ordered: false });

db.transactions.createIndex( { account_id: 1 } );
db.transactions.createIndex( { timestamp: 1 } );
db.transactions.createIndex( { transaction_id: 1 }, { unique: true, sparse: true} );

// Optionally create a capped collection for transaction logs that automatically removes older entries after a certain size or time limit is reached.
// Useful if you only need to keep a limited history of transactions.
db.createCollection("transactions", { capped: true, size: 104857600, max: 100000 }); // 100MB size limit, 100000 documents max. Adjust as needed.