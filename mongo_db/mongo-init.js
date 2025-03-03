
db = db.getSiblingDB("banking_ledger_db");

db.createUser({
  user: "ledger",
  pwd: "ledger",
  roles: [
    {
      role: "readWrite",
      db: "banking_ledger_db"
    }
  ]
});

db.createCollection("transactions");

db.transactions.insertMany([
  {
    id: "d6263bc8-0eeb-4195-9e64-81abd6d5685c",
    accountId: "8db6626d-5e84-4c4e-8cec-7dc54cb20ff5",
    amount: 100.00,
    transactionType: "credit",
    acceptedAt: new Date(),
    processedAt: new Date(),
    details: "Initial deposit",
    status: "success"
  },
  {
    id: "152a42be-63b6-46f9-919e-ba3996eaa890",
    accountId: "8db6626d-5e84-4c4e-8cec-7dc54cb20ff5",
    amount: 50.00,
    transactionType: "debit",
    acceptedAt: new Date(),
    processedAt: new Date(),    
    details: "debit",
    status: "success"
  }
], { ordered: false });

db.transactions.createIndex( { id: 1 }, { unique: true, sparse: true} );
db.transactions.createIndex( { accountId: 1 } );
db.transactions.createIndex( { acceptedAt: 1 } );

// Optionally create a capped collection for transaction logs that automatically removes older entries after a certain size or time limit is reached.
// Useful if you only need to keep a limited history of transactions.
db.createCollection("transactions", { capped: true, size: 104857600, max: 100000 }); // 100MB size limit, 100000 documents max. Adjust as needed.